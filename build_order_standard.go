package main

import (
	"math"
	"math/rand"

	"github.com/aiseeq/s2l/lib/point"
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/api"
	"github.com/aiseeq/s2l/protocol/enums/ability"
	"github.com/aiseeq/s2l/protocol/enums/terran"
)

var Standard = Strategy{
	Name: "Standard",
	Steps: BuildOrder{
		&SupplyDepotStep,
		buildingStep("Barracks", terran.Barracks, ability.Build_Barracks, 1, terran.SupplyDepot),
		refineryStep(1),
		orbitalCommandStep(1),
		addonStep("Barracks Reactor", terran.Barracks, terran.BarracksReactor, ability.Build_Reactor_Barracks, 1),
		// TODO: Turn Expand() into a step and insert it here. It is extremely funny
		// to see it continuously create command centers at home and send them
		// flying.
		// buildingStep("Command Center", terran.CommandCenter, ability.Build_CommandCenter, 2),
		&TrainMarine,
		buildingStep("Barracks", terran.Barracks, ability.Build_Barracks, 3, terran.SupplyDepot),
		orbitalCommandStep(2),
		addonStep("Barracks Tech Lab", terran.Barracks, terran.BarracksTechLab, ability.Build_TechLab_Barracks, 2),
		upgradeStep("Combat Shield", ability.Research_CombatShield, terran.BarracksTechLab),
		upgradeStep("Stimpack", ability.Research_Stimpack, terran.BarracksTechLab),
		buildingStep("Factory", terran.Factory, ability.Build_Factory, 1, terran.BarracksTechLab),
		buildingStep("Engineering Bay", terran.EngineeringBay, ability.Build_EngineeringBay, 1, terran.SupplyDepot),
		upgradeStep("Infantry Weapons Level 1", ability.Research_TerranInfantryWeaponsLevel1, terran.EngineeringBay),
		attackWaveStep(firstWaveConfig()),
		refineryStep(4),
		buildingStep("Barracks", terran.Barracks, ability.Build_Barracks, 5, terran.SupplyDepot),
		buildingStep("Starport", terran.Starport, ability.Build_Starport, 1, terran.Factory),
		addonStep("Factory Tech Lab", terran.Factory, terran.FactoryTechLab, ability.Build_TechLab_Factory, 1), // Factory Reactor
		addonStep("Barracks Reactor", terran.Barracks, terran.BarracksReactor, ability.Build_Reactor_Barracks, 3),
		// Switch Starport and Factory
		addonStep("Starport Reactor", terran.Starport, terran.StarportReactor, ability.Build_Reactor_Starport, 1), // Factory Tech Lab
		// Medivac (x4)
		// Siege Tank (x2)
		upgradeStep("Infantry Armor Level 1", ability.Research_TerranInfantryArmorLevel1, terran.EngineeringBay),

		// At this point, we should have enough units to launch a bigger attack.
		// TODO: Update to a second wave
		attackWaveStep(firstWaveConfig()),

		// These are just in the meantime
		buildingStep("Armory", terran.Armory, ability.Build_Armory, 1, terran.Factory),

		upgradeStep("Infantry Weapons Level 2", ability.Research_TerranInfantryWeaponsLevel2, terran.EngineeringBay),
		upgradeStep("Infantry Armor Level 2", ability.Research_TerranInfantryArmorLevel2, terran.EngineeringBay),
		attackWaveStep(firstWaveConfig()),

		upgradeStep("Infantry Weapons Level 3", ability.Research_TerranInfantryWeaponsLevel3, terran.EngineeringBay),
		upgradeStep("Infantry Armor Level 3", ability.Research_TerranInfantryArmorLevel3, terran.EngineeringBay),
		attackWaveStep(firstWaveConfig()),

		attackWaveStep(fullSupplyWaveConfig()),
	},
}

// SupplyDepotStep builds supply depots.
//
// TODO: Continuously make one supply depot after the other. Don't wait until
// we're supply blocked.
var SupplyDepotStep = BuildStep{
	Name: "Supply Depot",
	Predicate: func(b *Bot) bool {
		if !b.CanBuy(ability.Build_SupplyDepot) {
			return false
		}

		currSupply := b.Obs.PlayerCommon.FoodUsed
		maxSupply := b.Obs.PlayerCommon.FoodCap
		supplyLeft := maxSupply - currSupply

		if maxSupply >= 200 {
			return false
		}

		townHalls := b.findTownHalls()
		production := b.findProductionStructures()

		structures := make(scl.Units, 0, len(townHalls)+len(production))
		structures = append(structures, townHalls...)
		structures = append(structures, production...)
		if structures.Empty() {
			return false
		}

		// Calculate how much supply we'll use during depot construction
		timeForScv := float64(BuildTimeSCV) / float64(structures.Len())
		scvDuringDepots := uint32(math.Ceil(float64(BuildTimeSupplyDepot) / timeForScv))

		// Don't build if we have enough supply or if already building
		depotsOrdered := b.findWorkers().Filter(IsOrderedTo(ability.Build_SupplyDepot)).Len() >= 1
		depotsInProgress := b.Units.My.OfType(terran.SupplyDepot).Filter(IsInProgress).Len() >= 1

		return supplyLeft <= scvDuringDepots && !depotsOrdered && !depotsInProgress
	},

	Execute: func(b *Bot) {
		townHalls := b.findTownHalls().Filter(IsCcAtExpansion(b.state.CcForExp))
		if townHalls.Empty() {
			return
		}

		randomTownHall := townHalls[rand.Intn(len(townHalls))]

		// Find a good position for the supply depot
		pos := b.whereToBuild(randomTownHall.Point(), scl.S2x2, terran.SupplyDepot, ability.Build_SupplyDepot)
		if pos == nil {
			return
		}

		builder := b.findIdleOrGatheringWorkers().ClosestTo(pos)
		if builder == nil {
			return
		}

		// Go back to the closest resource after building
		if resource := b.findResourcesNearTownHalls(townHalls).ClosestTo(pos); resource != nil {
			logger.Info("Building supply depot at %v and queuing to gather at %v", *pos, resource.Point())

			builder.CommandPos(ability.Build_SupplyDepot, pos)
			builder.CommandTagQueue(ability.Smart, resource.Tag)

			if resource.IsMineral() {
				b.Miners.MineralForMiner[builder.Tag] = resource.Tag
			}

			if resource.IsGeyser() {
				b.Miners.GasForMiner[builder.Tag] = resource.Tag
			}
		} else {
			logger.Info("Building supply depot at %v", *pos)
			builder.CommandPos(ability.Build_SupplyDepot, pos)
		}

		b.DeductResources(ability.Build_SupplyDepot)
	},

	Next: func(b *Bot) bool {
		return b.Units.My.OfType(terran.SupplyDepot).Len() >= 1 || b.findWorkers().Filter(IsOrderedTo(ability.Build_SupplyDepot)).Exists()
	},
}

func buildingStep(name string, buildingId api.UnitTypeID, abilityId api.AbilityID, quantity int, requirements ...api.UnitTypeID) *BuildStep {
	return &BuildStep{
		Name: name,
		Predicate: func(b *Bot) bool {
			for _, requirement := range requirements {
				if b.Units.My.OfType(requirement).Empty() {
					return false
				}
			}

			if !b.CanBuy(abilityId) {
				return false
			}

			buildings := b.Units.My.OfType(buildingId)
			ordered := b.findWorkers().Filter(IsOrderedTo(abilityId))
			inProgress := buildings.Filter(IsInProgress)
			notStarted := ordered.Len() - inProgress.Len()
			return buildings.Len()+notStarted < quantity
		},

		Execute: func(b *Bot) {
			b.build(name, buildingId, abilityId, scl.S5x3)
		},

		Next: func(b *Bot) bool {
			buildings := b.Units.My.OfType(buildingId)
			ordered := b.findWorkers().Filter(IsOrderedTo(abilityId))
			inProgress := buildings.Filter(IsInProgress)
			notStarted := ordered.Len() - inProgress.Len()
			return buildings.Len()+notStarted >= quantity
		},
	}
}

func refineryStep(quantity int) *BuildStep {
	return &BuildStep{
		Name: "Refinery",
		Predicate: func(b *Bot) bool {
			if !b.CanBuy(ability.Build_Refinery) {
				return false
			}

			refineries := b.Units.My.OfType(terran.Refinery, terran.RefineryRich)
			ordered := b.findWorkers().Filter(IsOrderedTo(ability.Build_Refinery))
			inProgress := refineries.Filter(IsInProgress)
			notStarted := ordered.Len() - inProgress.Len()

			return refineries.Len()+notStarted < quantity
		},

		Execute: func(b *Bot) {
			townHalls := b.findTownHalls()
			vespeneGeysers := b.findVespeneGeysersNearTownHalls(townHalls)
			if townHalls.Empty() || vespeneGeysers.Empty() {
				return
			}

			randomVespeneGeyser := vespeneGeysers[rand.Intn(len(vespeneGeysers))]
			worker := b.findIdleOrGatheringWorkers().ClosestTo(randomVespeneGeyser)
			if worker == nil {
				return
			}

			logger.Info("Building refinery at %v", randomVespeneGeyser.Point())
			worker.CommandTag(ability.Build_Refinery, randomVespeneGeyser.Tag)
			worker.CommandTagQueue(ability.Smart, randomVespeneGeyser.Tag)
			b.DeductResources(ability.Build_Refinery)
			b.Miners.GasForMiner[worker.Tag] = randomVespeneGeyser.Tag
		},

		Next: func(b *Bot) bool {
			refineries := b.Units.My.OfType(terran.Refinery, terran.RefineryRich)
			ordered := b.findWorkers().Filter(IsOrderedTo(ability.Build_Refinery))
			inProgress := refineries.Filter(IsInProgress)
			notStarted := ordered.Len() - inProgress.Len()
			return refineries.Len()+notStarted >= quantity
		},
	}
}

func orbitalCommandStep(quantity int) *BuildStep {
	return &BuildStep{
		Name: "Orbital Command",
		Predicate: func(b *Bot) bool {
			barracks := b.Units.My.OfType(terran.Barracks).Filter(scl.Ready, scl.Ground)
			if barracks.Empty() {
				return false
			}

			if !b.CanBuy(ability.Morph_OrbitalCommand) {
				return false
			}

			orbitalCommands := b.Units.My.OfType(terran.OrbitalCommand, terran.OrbitalCommandFlying)
			inProgress := b.Units.My.OfType(terran.CommandCenter).Filter(IsOrderedTo(ability.Morph_OrbitalCommand))
			if orbitalCommands.Len()+inProgress.Len() >= quantity {
				b.state.CcForOrbitalCommand = 0
				return false
			}

			return true
		},

		Execute: func(b *Bot) {
			if b.state.CcForOrbitalCommand == 0 {
				// There's no command center marked for morphing into an orbital
				// command, so let's mark one
				commandCenters := b.Units.My.OfType(terran.CommandCenter).Filter(scl.Ready)
				if commandCenters.Empty() {
					return
				}

				randomCommandCenter := commandCenters[rand.Intn(len(commandCenters))]
				b.state.CcForOrbitalCommand = randomCommandCenter.Tag
			}

			// Check if the marked command center is still valid
			reserved := b.Units.ByTag[b.state.CcForOrbitalCommand]
			if reserved == nil || !reserved.Is(terran.CommandCenter, terran.CommandCenterFlying) {
				b.state.CcForOrbitalCommand = 0
				return
			}

			// If it's ordered to do anything else, cancel it
			ordered := IsOrderedTo(ability.Morph_OrbitalCommand)(reserved)
			if len(reserved.Orders) > 0 && !ordered {
				reserved.Command(ability.Cancel_Last)
				return
			}

			// If it's not morphing yet, morph it
			if !ordered {
				logger.Info("Morphing orbital command at %v", reserved.Point())
				reserved.Command(ability.Morph_OrbitalCommand)
				b.DeductResources(ability.Morph_OrbitalCommand)
				b.state.CcForOrbitalCommand = 0
			}
		},

		Next: func(b *Bot) bool {
			orbitalCommands := b.Units.My.OfType(terran.OrbitalCommand, terran.OrbitalCommandFlying)
			inProgress := b.Units.My.OfType(terran.CommandCenter).Filter(IsOrderedTo(ability.Morph_OrbitalCommand))
			return orbitalCommands.Len()+inProgress.Len() >= quantity
		},
	}
}

// addonStep manages building add-ons for buildings.
//
// TODO: Check if there's enough space for the reactor, and if not, fly the
// building somewhere safe.
func addonStep(name string, buildingId api.UnitTypeID, addonId api.UnitTypeID, abilityId api.AbilityID, quantity int) *BuildStep {
	return &BuildStep{
		Name: name,
		Predicate: func(b *Bot) bool {
			buildings := b.Units.My.OfType(buildingId).Filter(scl.Ready, scl.Ground, scl.NoAddon)
			if buildings.Empty() {
				return false
			}

			if !b.CanBuy(abilityId) {
				return false
			}

			if b.Units.My.OfType(addonId).Len() >= quantity {
				return false
			}

			return true
		},

		Execute: func(b *Bot) {
			buildings := b.Units.My.OfType(buildingId).Filter(scl.Ready, scl.Ground, scl.NoAddon)
			if buildings.Empty() {
				return
			}

			// If there's no building marked for add-on, mark one
			if b.state.BuildingForAddOn == 0 {
				randomBuilding := buildings[rand.Intn(len(buildings))]
				b.state.BuildingForAddOn = randomBuilding.Tag
			}

			// Check if the marked building is still valid
			reserved := b.Units.ByTag[b.state.BuildingForAddOn]
			if reserved == nil || !reserved.Is(buildingId) || reserved.AddOnTag != 0 {
				b.state.BuildingForAddOn = 0
				return
			}

			// If it's ordered to do anything else, cancel it
			ordered := IsOrderedTo(abilityId)(reserved)
			if len(reserved.Orders) > 0 && !ordered {
				reserved.Command(ability.Cancel_Last)
				return
			}

			logger.Info("Building %s at %v", name, reserved.Point())
			reserved.Command(abilityId)

			// In case it fails, queue the add-on to a new location
			elsewhere := b.whereToBuild(reserved.Point(), scl.S5x3, addonId, abilityId)
			reserved.CommandPosQueue(abilityId, elsewhere)

			b.DeductResources(abilityId)
			b.state.BuildingForAddOn = 0
		},

		Next: func(b *Bot) bool {
			return b.Units.My.OfType(addonId).Len() >= quantity
		},
	}
}

var TrainMarine = BuildStep{
	Name: "Train Marine",
	Predicate: func(b *Bot) bool {
		return b.CanBuy(ability.Train_Marine)
	},

	Execute: func(b *Bot) {
		barracks := b.Units.My.OfType(terran.Barracks).Filter(scl.Ready, scl.Ground, scl.Idle, IsNotTag(b.state.BuildingForAddOn))
		if barracks.Empty() {
			return
		}

		for _, barrack := range barracks {
			amount := b.amountTrainMarines(barrack)
			if amount == 0 {
				break
			}

			if rally := b.rallyPoint(); rally != nil {
				barrack.CommandPos(ability.Rally_Building, rally)
			}

			if amount == 1 {
				barrack.CommandQueue(ability.Train_Marine)
				logger.Info("Training one marine at %v", barrack.Point())
			}

			if amount == 2 {
				barrack.CommandQueue(ability.Train_Marine)
				barrack.CommandQueue(ability.Train_Marine)
				logger.Info("Training two marines at %v", barrack.Point())
			}
		}
	},

	Next: func(b *Bot) bool {
		return true
	},
}

func (b *Bot) amountTrainMarines(barracks *scl.Unit) int {
	if !b.CanBuy(ability.Train_Marine) {
		return 0
	}

	// Confirmed that we're about to train one marine.
	b.DeductResources(ability.Train_Marine)

	if barracks.AddOnTag == 0 {
		return 1
	}

	addon := b.Units.ByTag[barracks.AddOnTag]
	if addon == nil || !addon.Is(terran.BarracksReactor) || !b.CanBuy(ability.Train_Marine) {
		return 1
	}

	// Confirmed that we're about to train two marines.
	b.DeductResources(ability.Train_Marine)
	return 2
}

func (b *Bot) rallyPoint() *point.Point {
	townHalls := b.findTownHalls().Filter(IsCcAtExpansion(b.state.CcForExp))
	if townHalls.Empty() {
		return nil
	}

	closest := townHalls.ClosestTo(b.Locs.EnemyStart)
	rally := closest.Towards(b.Locs.EnemyStart, closest.SightRange())
	return &rally
}

func (b *Bot) build(name string, buildingId api.UnitTypeID, abilityId api.AbilityID, size scl.BuildingSize) {
	if !b.CanBuy(abilityId) {
		return
	}

	townHalls := b.findTownHalls().Filter(IsCcAtExpansion(b.state.CcForExp))
	if townHalls.Empty() {
		return
	}

	randomTownHall := townHalls[rand.Intn(len(townHalls))]

	pos := b.whereToBuild(randomTownHall.Point(), size, buildingId, abilityId)
	if pos == nil {
		return
	}

	builder := b.findIdleOrGatheringWorkers().ClosestTo(pos)
	if builder == nil {
		return
	}

	if resource := b.findResourcesNearTownHalls(townHalls).ClosestTo(pos); resource != nil {
		logger.Info("Building %s at %v and queuing to gather at %v", name, *pos, resource.Point())

		builder.CommandPos(abilityId, pos)
		builder.CommandTagQueue(ability.Smart, resource.Tag)

		if resource.IsMineral() {
			b.Miners.MineralForMiner[builder.Tag] = resource.Tag
		}

		if resource.IsGeyser() {
			b.Miners.GasForMiner[builder.Tag] = resource.Tag
		}
	} else {
		logger.Info("Building %s at %v", name, *pos)
		builder.CommandPos(abilityId, pos)
	}

	b.DeductResources(abilityId)
}

func upgradeStep(name string, abilityId api.AbilityID, buildingId api.UnitTypeID) *BuildStep {
	return &BuildStep{
		Name: name,
		Predicate: func(b *Bot) bool {
			if !b.CanBuy(abilityId) {
				return false
			}

			if b.Upgrades[abilityId] {
				return false
			}

			if b.Units.My.OfType(buildingId).Filter(IsOrderedTo(abilityId)).Exists() {
				return false
			}

			return true
		},

		Execute: func(b *Bot) {
			buildings := b.Units.My.OfType(buildingId).Filter(scl.Ready, scl.Idle)
			if buildings.Empty() {
				return
			}

			logger.Info("Researching %s", name)
			buildings.First().Command(abilityId)
			b.DeductResources(abilityId)
		},

		Next: func(b *Bot) bool {
			if b.Upgrades[abilityId] {
				return true
			}

			if b.Units.My.OfType(buildingId).Filter(IsOrderedTo(abilityId)).Exists() {
				return true
			}

			return false
		},
	}
}

func attackWaveStep(config *AttackWaveConfig) *BuildStep {
	return &BuildStep{
		Name:      config.Name,
		Predicate: config.Predicate,
		Execute:   config.Execute,
		Next: func(b *Bot) bool {
			return true
		},
	}
}
