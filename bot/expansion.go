package bot

import (
	"fmt"

	"github.com/NatoBoram/BlackCompany/filter"
	"github.com/NatoBoram/BlackCompany/log"
	"github.com/NatoBoram/BlackCompany/sight"
	"github.com/aiseeq/s2l/lib/point"
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/api"
	"github.com/aiseeq/s2l/protocol/enums/ability"
	"github.com/aiseeq/s2l/protocol/enums/terran"
	"github.com/jwalton/gchalk"
)

const (
	// MaxWorkers is the maximum number of workers that can be trained.
	MaxWorkers = 80
)

// FindExpansionLocations finds the next best available expansion locations.
func (b *Bot) FindExpansionLocations() point.Points {
	locations := make(point.Points, 0, b.Locs.MyExps.Len()+1)
	expansions := append(b.Locs.MyExps, b.Locs.MyStart)
	townHalls := b.FindTownHalls()

	for _, expansion := range expansions {
		// Skip existing expansions
		if townHalls.CloserThan(1, expansion).Exists() {
			continue
		}

		// Skip reserved locations
		if b.State.CcForExp.IsReserved(b, expansion) {
			continue
		}

		// Skip locations that would be unsafe
		if b.Enemies.Visible.Filter(scl.DpsGt5).CloserThan(sight.LineOfSightScannerSweep.Float64(), expansion).Exists() {
			continue
		}

		// If the expansion is not explored, then its mineral content shows up as
		// empty. Let's just assume it's full.
		fieldsWithMinerals := b.Units.Minerals.All().
			CloserThan(scl.ResourceSpreadDistance, expansion).
			Filter(filter.HasMinerals)
		isExplored := b.Grid.IsExplored(expansion)
		if isExplored && fieldsWithMinerals.Empty() {
			continue
		}

		// It's already reserved for a town hall.
		townHalls := b.State.CcForExp.ByExpansion(b, expansion)
		if townHalls.Exists() {
			continue
		}

		locations = append(locations, expansion)
	}

	return locations
}

// FindAvailableCommandCenters finds command centers that are not at an
// expansion location that we can use to expand by lifting them.
func (b *Bot) FindAvailableCommandCenters() scl.Units {
	expansions := append(b.Locs.MyExps, b.Locs.MyStart)

	return b.Units.My.
		OfType(
			terran.CommandCenter, terran.OrbitalCommand,
			terran.CommandCenterFlying, terran.OrbitalCommandFlying,
		).
		Filter(
			filter.IsNotAtAny(expansions),
			filter.IsNotCcAtExpansion(b.State.CcForExp),
			filter.IsNotOrderedToAny(
				ability.Lift, ability.Lift_CommandCenter, ability.Lift_OrbitalCommand,
				ability.Land, ability.Land_CommandCenter, ability.Land_OrbitalCommand,
			),
		)
}

// ShouldExpand returns whether we should expand or not.
//
// The current strategy is as follows:
//
//   - Don't build if a Command Center is in progress
//   - Don't build if there's more resource slots than [MaxWorkers]
func (b *Bot) ShouldExpand() bool {
	if !b.CanBuy(ability.Build_CommandCenter) {
		return false
	}

	ccOrdered := b.FindWorkers().Filter(filter.IsOrderedTo(ability.Build_CommandCenter))
	if ccOrdered.Exists() {
		return false
	}

	ccInProgress := b.Units.My.OfType(terran.CommandCenter).Filter(filter.IsInProgress)
	if ccInProgress.Exists() {
		return false
	}

	townHalls := b.FindTownHalls()
	mineralFields := b.FindMineralFieldsNearTownHalls(townHalls)
	claimedVespeneGeysers := b.FindClaimedVespeneGeysersNearTownHalls(townHalls)

	mineralSlots := mineralFields.Len() * 2
	vespeneGeyserSlots := claimedVespeneGeysers.Len() * 3
	resourceSlots := mineralSlots + vespeneGeyserSlots

	// Literally enough resources to saturate the maximum amount of workers we'll
	// train. This limits us to 4 fresh bases.
	if resourceSlots > MaxWorkers {
		return false
	}

	return true
}

// IsFlyingFaster calculates whether flying a CC to the target location is
// faster than having a worker walk there to build it.
//
// Does not take into account lifting time nor landing time.
func (b *Bot) IsFlyingFaster(worker *scl.Unit, base point.Pointer, expansion point.Point) bool {
	walkTime := b.walkTime(worker, expansion)
	flyTime := b.flyTime(base, terran.CommandCenterFlying, expansion)

	if walkTime == flyTime {
		log.Debug("Walking time: %f, fly time: %f", walkTime, flyTime)
	} else if walkTime > flyTime {
		flyString := fmt.Sprintf("%f", flyTime)
		log.Debug("Walking time: %f, fly time: %s", walkTime, gchalk.Bold(flyString))
	} else if walkTime < flyTime {
		walkString := fmt.Sprintf("%f", walkTime)
		log.Debug("Walking time: %s, fly time: %f", gchalk.Bold(walkString), flyTime)
	}

	return walkTime >= flyTime
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
