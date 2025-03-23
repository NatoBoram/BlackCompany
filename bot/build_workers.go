package bot

import (
	"github.com/NatoBoram/BlackCompany/filter"
	"github.com/NatoBoram/BlackCompany/log"
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/enums/ability"
	"github.com/aiseeq/s2l/protocol/enums/terran"
)

const (
	// MaxWorkers is the maximum number of workers that can be trained.
	MaxWorkers = 80
)

// BuildWorker trains SCVs from command centers.
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
func (b *Bot) BuildWorker() {
	if !b.CanBuy(ability.Train_SCV) || b.findMiners().Len() >= MaxWorkers {
		return
	}

	townHalls := b.Units.My.OfType(
		terran.CommandCenter, terran.OrbitalCommand, terran.PlanetaryFortress,
	)
	if townHalls.Empty() {
		return
	}

	resources := b.findUnsaturatedResourcesNearTownHalls(townHalls)
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
