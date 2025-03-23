package bot

import (
	"github.com/NatoBoram/BlackCompany/filter"
	"github.com/NatoBoram/BlackCompany/log"
	"github.com/NatoBoram/BlackCompany/sight"
	"github.com/aiseeq/s2l/lib/point"
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/api"
	"github.com/aiseeq/s2l/protocol/enums/ability"
	"github.com/aiseeq/s2l/protocol/enums/terran"
)

// Expand expands the bot's base whenever enough resources are available.
func (b *Bot) Expand() {
	expansions := b.findExpansionLocations()
	if expansions.Empty() {
		return
	}

	ccs := b.hasFreeCommandCenter()
	for i, cc := range ccs {
		if i >= len(expansions) {
			break
		}
		expansion := expansions[i]

		if b.State.CcForExp == nil {
			b.State.CcForExp = make(map[api.UnitTag]point.Point)
		}

		if !cc.IsFlying && cc.IsReady() {
			log.Info("Lifting Command Center from %s to expansion %s", cc.Point(), expansion)
			cc.Command(ability.Lift_CommandCenter)

			b.State.CcForExp[cc.Tag] = expansion
		}

		if cc.IsFlying && cc.IsIdle() {
			cc.CommandPosQueue(ability.Land_CommandCenter, expansion)

			b.State.CcForExp[cc.Tag] = expansion
		}

	}
	if ccs.Exists() {
		return
	}

	expansion := expansions[0]

	shouldExpand := b.shouldExpand()
	if !shouldExpand {
		return
	}

	workers := b.findIdleOrGatheringWorkers()
	if workers.Empty() {
		return
	}

	worker := workers.ClosestTo(expansion)
	if worker == nil {
		return
	}

	townHalls := b.findTownHalls()
	if townHalls.Empty() {
		log.Info("No town halls found, building Command Center at expansion %s", expansion)
		location := b.whereToBuild(expansion, scl.S5x5, terran.CommandCenter, ability.Build_CommandCenter)
		worker.CommandPos(ability.Build_CommandCenter, location)
		b.DeductResources(ability.Build_CommandCenter)
		return
	}

	nearestTownHall := townHalls.ClosestTo(worker)
	towards := nearestTownHall.Point().Towards(expansion, nearestTownHall.SightRange())
	location := b.whereToBuild(towards, scl.S5x5, terran.CommandCenter, ability.Build_CommandCenter)

	// So do I build it there or near worker then fly it over?
	if b.isFlyingFaster(worker, location, expansion) {
		log.Info("Building Command Center at base %s to fly to expansion %s", location, expansion)

		worker.CommandPos(ability.Build_CommandCenter, location)
		b.DeductResources(ability.Build_CommandCenter)

		closestMineralField := b.Units.Minerals.All().ClosestTo(expansion)
		worker.CommandTagQueue(ability.Smart, closestMineralField.Tag)

		return
	}

	log.Info("Expanding to %s", expansion)

	worker.CommandPos(ability.Build_CommandCenter, expansion)
	b.DeductResources(ability.Build_CommandCenter)

	closestMineralField := b.Units.Minerals.All().ClosestTo(expansion)
	worker.CommandTagQueue(ability.Smart, closestMineralField.Tag)
}

// isFlyingFaster calculates whether flying a CC to the target location is
// faster than having a worker walk there to build it.
//
// Does not take into account lifting time nor landing time.
func (b *Bot) isFlyingFaster(worker *scl.Unit, base point.Pointer, expansion point.Point) bool {
	flyTime := b.flyTime(base, terran.CommandCenterFlying, expansion)
	walkTime := b.walkTime(worker, expansion)

	log.Debug("Worker travel time: %f, fly time: %f", walkTime, flyTime)

	return flyTime < walkTime
}

func (b *Bot) walkTime(unit *scl.Unit, destination point.Point) float64 {
	walkDistance := b.RequestPathing(unit, destination)
	if walkDistance == 0 {
		return 0
	}

	return walkDistance / unit.Speed()
}

func (b *Bot) flyTime(origin point.Pointer, unit api.UnitTypeID, destination point.Point) float64 {
	flySpeed := float64(b.U.Types[unit].MovementSpeed)

	flyDistance := origin.Point().Dist(destination)
	if flyDistance == 0 {
		return 0
	}

	return flyDistance / flySpeed
}

// hasFreeCommandCenter checks if there's a command center that's not at an
// expansion location that we can use to expand by lifting it.
func (b *Bot) hasFreeCommandCenter() scl.Units {
	commandCenters := b.Units.My.OfType(terran.CommandCenter)
	flying := b.Units.My.OfType(terran.CommandCenterFlying)
	free := make(scl.Units, 0, commandCenters.Len()+flying.Len())

	if flying.Exists() {
		free = append(free, flying...)
	}

	mineralFields := b.Units.Minerals.All()
	for _, cc := range commandCenters {
		hasMineralFields := mineralFields.CloserThan(scl.ResourceSpreadDistance, cc).Filter(filter.SameHeightAs(cc)).Exists()
		if !hasMineralFields {
			free = append(free, cc)
		}
	}

	return free
}

// findExpansionLocation finds the next best available expansion location.
func (b *Bot) findExpansionLocations() point.Points {
	locations := make(point.Points, 0, b.Locs.MyExps.Len())

	for _, expansion := range b.Locs.MyExps {
		// Skip existing expansions
		if b.findTownHalls().
			CloserThan(scl.ResourceSpreadDistance, expansion).Exists() {
			continue
		}

		// Skip locations that would be unsafe
		if b.Enemies.Visible.Filter(scl.DpsGt5).CloserThan(sight.LineOfSightScannerSweep.Float64(), expansion).Exists() {
			continue
		}

		// If the expansion is not explored, then its mineral content shows up as
		// empty. Let's just assume it's full.
		hasMinerals := b.Units.Minerals.All().CloserThan(scl.ResourceSpreadDistance, expansion).Exists()
		isExplored := b.Grid.IsExplored(expansion)
		if isExplored && !hasMinerals {
			continue
		}

		locations = append(locations, expansion)
	}

	return locations
}

// shouldExpand returns whether we should expand or not.
//
// The current strategy is as follows:
//
//   - Don't build if a Command Center is in progress
//   - Don't build if there's more resource slots than miners + the amount of
//     SCVs it takes to build a Command Center
func (b *Bot) shouldExpand() bool {
	if !b.CanBuy(ability.Build_CommandCenter) {
		return false
	}

	ccOrdered := b.findWorkers().Filter(filter.IsOrderedTo(ability.Build_CommandCenter)).Len() >= 1
	if ccOrdered {
		return false
	}

	ccInProgress := b.Units.My.OfType(terran.CommandCenter).Filter(filter.IsInProgress).Len() >= 1
	if ccInProgress {
		return false
	}

	miners := b.findMiners()
	if miners.Empty() {
		return false
	}

	townHalls := b.findTownHalls()
	mineralFields := b.findMineralFieldsNearTownHalls(townHalls)
	claimedVespeneGeysers := b.findClaimedVespeneGeysersNearTownHalls(townHalls)

	scvDuringCc := buildDuring(BuildTimeSCV, BuildTimeCommandCenter)

	mineralSlots := mineralFields.Len()*2 + claimedVespeneGeysers.Len()*3
	vespeneGeyserSlots := claimedVespeneGeysers.Len() * 3
	resourceSlots := mineralSlots + vespeneGeyserSlots

	// Literally enough resources to saturate our workers and then some
	if resourceSlots > miners.Len()+scvDuringCc {
		return false
	}

	return true
}
