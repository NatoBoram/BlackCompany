package micro

import (
	"github.com/NatoBoram/BlackCompany/bot"
	"github.com/NatoBoram/BlackCompany/filter"
	"github.com/NatoBoram/BlackCompany/log"
	"github.com/NatoBoram/BlackCompany/sight"
	"github.com/aiseeq/s2l/lib/point"
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/api"
	"github.com/aiseeq/s2l/protocol/enums/ability"
	"github.com/aiseeq/s2l/protocol/enums/terran"
)

const (
	// MaxWorkers is the maximum number of workers that can be trained.
	MaxWorkers = 80
)

func handleTownHalls(b *bot.Bot) {
	trainWorkers(b)
	expand(b)
}

// expand expands the bot's base whenever enough resources are available.
func expand(b *bot.Bot) {
	expansions := findExpansionLocations(b)
	if expansions.Empty() {
		return
	}

	ccs := hasFreeCommandCenter(b)
	for i, cc := range ccs {
		if i >= len(expansions) {
			break
		}
		expansion := expansions[i]

		if !cc.IsFlying && cc.IsReady() {
			if cc.Is(terran.CommandCenter) {
				log.Info("Lifting Command Center from %s to expansion %s", cc.Point(), expansion)
				cc.Command(ability.Lift_CommandCenter)
			}

			if cc.Is(terran.OrbitalCommand) {
				log.Info("Lifting Orbital Command from %s to expansion %s", cc.Point(), expansion)
				cc.Command(ability.Lift_OrbitalCommand)
			}

			b.State.CcForExp[cc.Tag] = expansion
		}

		if cc.IsFlying && cc.IsIdle() {
			if cc.Is(terran.CommandCenterFlying) {
				cc.CommandPosQueue(ability.Land_CommandCenter, expansion)
			}
			if cc.Is(terran.OrbitalCommandFlying) {
				cc.CommandPosQueue(ability.Land_OrbitalCommand, expansion)
			}

			b.State.CcForExp[cc.Tag] = expansion
		}

	}
	if ccs.Exists() {
		return
	}

	expansion := expansions[0]

	shouldExpand := shouldExpand(b)
	if !shouldExpand {
		return
	}

	workers := b.FindIdleOrGatheringWorkers()
	if workers.Empty() {
		return
	}

	worker := workers.ClosestTo(expansion)
	if worker == nil {
		return
	}

	townHalls := b.FindTownHalls()
	if townHalls.Empty() {
		log.Info("No town halls found, building Command Center at expansion %s", expansion)
		location := b.WhereToBuild(expansion, scl.S5x5, terran.CommandCenter, ability.Build_CommandCenter)
		worker.CommandPos(ability.Build_CommandCenter, location)
		b.DeductResources(ability.Build_CommandCenter)
		return
	}

	nearestTownHall := townHalls.ClosestTo(worker)
	towards := nearestTownHall.Point().Towards(expansion, nearestTownHall.SightRange())
	location := b.WhereToBuild(towards, scl.S5x5, terran.CommandCenter, ability.Build_CommandCenter)

	// So do I build it there or near worker then fly it over?
	if isFlyingFaster(b, worker, location, expansion) {
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
func isFlyingFaster(b *bot.Bot, worker *scl.Unit, base point.Pointer, expansion point.Point) bool {
	flyTime := flyTime(b, base, terran.CommandCenterFlying, expansion)
	walkTime := walkTime(b, worker, expansion)

	log.Debug("Worker travel time: %f, fly time: %f", walkTime, flyTime)

	return flyTime < walkTime
}

func walkTime(b *bot.Bot, unit *scl.Unit, destination point.Point) float64 {
	walkDistance := b.RequestPathing(unit, destination)
	if walkDistance == 0 {
		return 0
	}

	return walkDistance / unit.Speed()
}

func flyTime(b *bot.Bot, origin point.Pointer, unit api.UnitTypeID, destination point.Point) float64 {
	flySpeed := float64(b.U.Types[unit].MovementSpeed)

	flyDistance := origin.Point().Dist(destination)
	if flyDistance == 0 {
		return 0
	}

	return flyDistance / flySpeed
}

// hasFreeCommandCenter checks if there's a command center that's not at an
// expansion location that we can use to expand by lifting it.
func hasFreeCommandCenter(b *bot.Bot) scl.Units {
	return b.Units.My.
		OfType(
			terran.CommandCenter, terran.OrbitalCommand,
			terran.CommandCenterFlying, terran.OrbitalCommandFlying,
		).
		Filter(
			filter.IsNotCcAtExpansion(b.State.CcForExp),
			filter.IsNotOrderedToAny(
				ability.Lift, ability.Lift_CommandCenter, ability.Lift_OrbitalCommand,
				ability.Land, ability.Land_CommandCenter, ability.Land_OrbitalCommand,
			),
		)
}

// findExpansionLocation finds the next best available expansion location.
func findExpansionLocations(b *bot.Bot) point.Points {
	locations := make(point.Points, 0, b.Locs.MyExps.Len()+1)
	expansions := append(b.Locs.MyExps, b.Locs.MyStart)
	townHalls := b.FindTownHalls()

	for _, expansion := range expansions {
		// Skip existing expansions
		if townHalls.CloserThan(scl.ResourceSpreadDistance, expansion).Exists() {
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
func shouldExpand(b *bot.Bot) bool {
	if !b.CanBuy(ability.Build_CommandCenter) {
		return false
	}

	ccOrdered := b.FindWorkers().Filter(filter.IsOrderedTo(ability.Build_CommandCenter)).Len() >= 1
	if ccOrdered {
		return false
	}

	ccInProgress := b.Units.My.OfType(terran.CommandCenter).Filter(filter.IsInProgress).Len() >= 1
	if ccInProgress {
		return false
	}

	miners := b.FindMiners()
	if miners.Empty() {
		return false
	}

	townHalls := b.FindTownHalls()
	mineralFields := b.FindMineralFieldsNearTownHalls(townHalls)
	claimedVespeneGeysers := b.FindClaimedVespeneGeysersNearTownHalls(townHalls)

	scvDuringCc := bot.BuildDuring(bot.BuildTimeSCV, bot.BuildTimeCommandCenter)

	mineralSlots := mineralFields.Len()*2 + claimedVespeneGeysers.Len()*3
	vespeneGeyserSlots := claimedVespeneGeysers.Len() * 3
	resourceSlots := mineralSlots + vespeneGeyserSlots

	// Literally enough resources to saturate our workers and then some
	if resourceSlots > miners.Len()+scvDuringCc {
		return false
	}

	return true
}

// trainWorkers trains SCVs from command centers.
//
//   - When SCVs can be afforded and there's less than 80 of them
//   - Find mineral fields that aren't depleted and count the missing SCVs to saturate them
//   - Find vespene geysers that aren't exhausted and count the missing SCVs to saturate them
//
// For each town halls:
//
//   - Get the closest resource that's not saturated
//   - Set the rally point to that resource
//   - Train a SCV
func trainWorkers(b *bot.Bot) {
	if !b.CanBuy(ability.Train_SCV) || b.FindMiners().Len() >= MaxWorkers {
		return
	}

	townHalls := b.Units.My.OfType(
		terran.CommandCenter, terran.OrbitalCommand, terran.PlanetaryFortress,
	)
	if townHalls.Empty() {
		return
	}

	resources := b.FindUnsaturatedResourcesNearTownHalls(townHalls)
	if resources.Empty() {
		return
	}

	idleTownHalls := townHalls.Filter(scl.Ready, scl.Idle, filter.IsCcAtExpansion(b.State.CcForExp))
	if idleTownHalls.Empty() {
		return
	}

	for _, cc := range idleTownHalls {
		if !b.CanBuy(ability.Train_SCV) || resources.Empty() {
			break
		}

		// Ignore command centers that are reserved for morphing into an orbital
		// command.
		if b.State.CcForOrbitalCommand == cc.Tag {
			if cc.Is(terran.OrbitalCommand, terran.OrbitalCommandFlying) {
				b.State.CcForOrbitalCommand = 0
			} else {
				continue
			}
		}

		var resource *scl.Unit

		resourcesNearby := resources.CloserThan(scl.ResourceSpreadDistance, cc)
		nearbyRefineries := resourcesNearby.Filter(filter.HasGas)
		nearbyMineralFields := resourcesNearby.Filter(filter.HasMinerals)
		if nearbyMineralFields.Exists() && nearbyRefineries.Exists() {
			if (nearbyMineralFields.Len()*2)/(nearbyRefineries.Len()*3) >= 16/6 {
				resource = nearbyRefineries.ClosestTo(cc)
			} else {
				resource = nearbyMineralFields.ClosestTo(cc)
			}
		} else {
			resource = resources.ClosestTo(cc)
		}

		log.Info("Training SCV for resource %v", resource.Point())
		cc.CommandTag(ability.Rally_CommandCenter, resource.Tag)
		cc.CommandQueue(ability.Train_SCV)
		b.DeductResources(ability.Train_SCV)
		resources.Remove(resource)
	}
}
