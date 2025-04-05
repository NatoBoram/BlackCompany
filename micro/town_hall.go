package micro

import (
	"github.com/NatoBoram/BlackCompany/bot"
	"github.com/NatoBoram/BlackCompany/filter"
	"github.com/NatoBoram/BlackCompany/log"
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/enums/ability"
	"github.com/aiseeq/s2l/protocol/enums/terran"
)

func handleTownHalls(b *bot.Bot) {
	flyToExpansion(b)
	trainWorkers(b)
}

func flyToExpansion(b *bot.Bot) {
	// Make sure every command center is destined to an expansion
	b.AssignCCToExpansion()

	// Find the misplaced ones and move them
	misplaced := b.State.CcForExp.Misplaced(b)
	for tag, expansion := range misplaced {
		unit := b.Units.ByTag[tag]
		if unit == nil {
			delete(b.State.CcForExp, tag)
			continue
		}

		if !unit.IsReady() {
			continue
		}

		if filter.IsOrderedToAny(
			ability.Lift, ability.Lift_CommandCenter, ability.Lift_OrbitalCommand,
			ability.Land, ability.Land_CommandCenter, ability.Land_OrbitalCommand,
		)(unit) {
			continue
		}

		if len(unit.Orders) > 0 {
			unit.Command(ability.Cancel_Last)
			continue
		}

		if unit.Is(terran.CommandCenter) {
			log.Info("Lifting Command Center from %s to expansion %s", unit.Point(), expansion)
			unit.Command(ability.Lift_CommandCenter)
		}

		if unit.Is(terran.OrbitalCommand) {
			log.Info("Lifting Orbital Command from %s to expansion %s", unit.Point(), expansion)
			unit.Command(ability.Lift_OrbitalCommand)
		}

		if unit.Is(terran.CommandCenterFlying) {
			log.Info("Flying Command Center from %s to expansion %s", unit.Point(), expansion)
			unit.CommandPos(ability.Land_CommandCenter, expansion)
		}

		if unit.Is(terran.OrbitalCommandFlying) {
			log.Info("Flying Orbital Command from %s to expansion %s", unit.Point(), expansion)
			unit.CommandPos(ability.Land_OrbitalCommand, expansion)
		}
	}

	flying := b.Units.My.OfType(
		terran.CommandCenterFlying, terran.OrbitalCommandFlying,
	).Filter(
		filter.IsNotOrderedToAny(
			ability.Lift, ability.Lift_CommandCenter, ability.Lift_OrbitalCommand,
			ability.Land, ability.Land_CommandCenter, ability.Land_OrbitalCommand,
		),
	)
	if flying.Empty() {
		return
	}

	// Land them where they want to be landed
	for _, unit := range flying {
		expansion, ok := b.State.CcForExp[unit.Tag]
		if !ok || expansion == 0 {
			log.Warn("Command center %s is flying but has no expansion assigned", unit.Point())
			continue
		}

		location := b.WhereToBuild(expansion, scl.S5x5, unit.UnitType, ability.Build_CommandCenter)

		if unit.Is(terran.CommandCenterFlying) {
			log.Info("Landing Command Center at %s", location)
			unit.CommandPos(ability.Land_CommandCenter, location)
		}

		if unit.Is(terran.OrbitalCommandFlying) {
			log.Info("Landing Orbital Command at %s", location)
			unit.CommandPos(ability.Land_OrbitalCommand, location)
		}
	}
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
	if !b.CanBuy(ability.Train_SCV) {
		return
	}

	townHalls := b.Units.My.OfType(
		terran.CommandCenter, terran.OrbitalCommand, terran.PlanetaryFortress,
	)
	if townHalls.Empty() {
		return
	}

	workers := b.FindWorkers()
	inProgress := townHalls.Filter(filter.IsOrderedToAny(ability.Train_SCV, ability.Train_Probe, ability.Train_Drone))
	if workers.Len()+inProgress.Len() >= bot.MaxWorkers {
		return
	}

	resources := b.FindUnsaturatedResourcesNearTownHalls(townHalls)
	if resources.Empty() {
		return
	}

	idleTownHalls := townHalls.Filter(
		scl.Ready, scl.Idle,
		filter.IsCcAtExpansion(b.State.CcForExp),
	)
	if idleTownHalls.Empty() {
		return
	}

	for _, cc := range idleTownHalls {
		if !b.CanBuy(ability.Train_SCV) || resources.Empty() {
			break
		}

		// Ignore command centers that are reserved for morphing into an orbital
		// command or a planetary fortress.
		if b.State.CcForOrbitalCommand == cc.Tag || b.State.CcForPlanetaryFortress == cc.Tag {
			if cc.Is(terran.OrbitalCommand, terran.OrbitalCommandFlying) {
				b.State.CcForOrbitalCommand = 0
			} else if cc.Is(terran.PlanetaryFortress) {
				b.State.CcForPlanetaryFortress = 0
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
		cc.CommandTag(ability.Rally_Building, resource.Tag)
		cc.CommandQueue(ability.Train_SCV)
		b.DeductResources(ability.Train_SCV)
		resources.Remove(resource)
	}
}
