package bot

import (
	"fmt"
	"math"
	"math/rand"
	"runtime/debug"
	"slices"

	"github.com/NatoBoram/BlackCompany/filter"
	"github.com/NatoBoram/BlackCompany/log"
	"github.com/aiseeq/s2l/lib/point"
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/api"
	"github.com/aiseeq/s2l/protocol/enums/ability"
	"github.com/aiseeq/s2l/protocol/enums/terran"
)

var Standard = Strategy{
	Name: "Standard",
	Steps: BuildOrder{
		&defenseWaveStep,
		&supplyDepotStep,
		chatVersionStep(),
		buildingStep("Barracks", terran.Barracks, ability.Build_Barracks, 1, terran.SupplyDepot),
		refineryStep(1),
		orbitalCommandStep(1),
		addonStep("Barracks Reactor", terran.Barracks, terran.BarracksReactor, ability.Build_Reactor_Barracks, 1),
		// TODO: Turn Expand() into a step and insert it here. It is extremely funny
		// to see it continuously create command centers at home and send them
		// flying.
		townHallStep(terran.CommandCenter, ability.Build_CommandCenter, 2),
		&trainMarine,
		attackWaveStep(fullSupplyWaveConfig()),
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

		// Things to do over and over again
		refineryStep(0),

		// These are just in the meantime
		buildingStep("Armory", terran.Armory, ability.Build_Armory, 1, terran.Factory),

		upgradeStep("Infantry Weapons Level 2", ability.Research_TerranInfantryWeaponsLevel2, terran.EngineeringBay),
		upgradeStep("Infantry Armor Level 2", ability.Research_TerranInfantryArmorLevel2, terran.EngineeringBay),

		upgradeStep("Infantry Weapons Level 3", ability.Research_TerranInfantryWeaponsLevel3, terran.EngineeringBay),
		upgradeStep("Infantry Armor Level 3", ability.Research_TerranInfantryArmorLevel3, terran.EngineeringBay),
	},
}

// supplyDepotStep builds supply depots.
//
// TODO: Continuously make one supply depot after the other. Don't wait until
// we're supply blocked.
var supplyDepotStep = BuildStep{
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

		townHalls := b.FindTownHalls()
		production := b.findProductionStructures()

		structures := slices.Concat(townHalls, production)
		if structures.Empty() {
			return false
		}

		// Calculate how much supply we'll use during depot construction
		timeForScv := float64(BuildTimeSCV) / float64(structures.Len())
		scvDuringDepots := uint32(math.Ceil(float64(BuildTimeSupplyDepot) / timeForScv))

		// Don't build if we have enough supply or if already building
		depotsOrdered := b.FindWorkers().Filter(filter.IsOrderedTo(ability.Build_SupplyDepot)).Len() >= 1
		depotsInProgress := b.Units.My.OfType(terran.SupplyDepot).Filter(filter.IsInProgress).Len() >= 1

		return supplyLeft <= scvDuringDepots && !depotsOrdered && !depotsInProgress
	},

	Execute: func(b *Bot) {
		townHalls := b.FindTownHalls().Filter(filter.IsCcAtExpansion(b.State.CcForExp))
		if townHalls.Empty() {
			return
		}

		randomTownHall := townHalls[rand.Intn(len(townHalls))]

		// Find a good position for the supply depot
		pos := b.WhereToBuild(randomTownHall.Point(), scl.S2x2, terran.SupplyDepot, ability.Build_SupplyDepot)
		if pos == nil {
			return
		}

		builder := b.FindIdleOrGatheringWorkers().ClosestTo(pos)
		if builder == nil {
			return
		}

		// Go back to the closest resource after building
		if resource := b.findResourcesNearTownHalls(townHalls).ClosestTo(pos); resource != nil {
			log.Info("Building supply depot at %v and queuing to gather at %v", *pos, resource.Point())

			builder.CommandPos(ability.Build_SupplyDepot, pos)
			builder.CommandTagQueue(ability.Smart, resource.Tag)

			if resource.IsMineral() {
				b.Miners.MineralForMiner[builder.Tag] = resource.Tag
			}

			if resource.IsGeyser() {
				b.Miners.GasForMiner[builder.Tag] = resource.Tag
			}
		} else {
			log.Info("Building supply depot at %v", *pos)
			builder.CommandPos(ability.Build_SupplyDepot, pos)
		}

		b.DeductResources(ability.Build_SupplyDepot)
	},

	Next: func(b *Bot) bool {
		return b.Units.My.OfType(terran.SupplyDepot).Len() >= 1 ||
			b.FindWorkers().Filter(filter.IsOrderedTo(ability.Build_SupplyDepot)).Exists()
	},
}

// chatVersionStep announces the version number of the bot.
func chatVersionStep() *BuildStep {
	announced := false

	return &BuildStep{
		Name: "Announce version number",
		Predicate: func(b *Bot) bool {
			return true
		},
		Execute: func(b *Bot) {
			if announced {
				return
			}

			// Print version information for everyone to enjoy
			info, ok := debug.ReadBuildInfo()
			if ok {
				message := fmt.Sprintf("BlackCompany %s", info.Main.Version)
				b.Actions.ChatSend(message, api.ActionChat_Team)
			}

			announced = true
		},
		Next: func(b *Bot) bool {
			return true
		},
	}
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
			ordered := b.FindWorkers().Filter(filter.IsOrderedTo(abilityId))
			inProgress := buildings.Filter(filter.IsInProgress)
			notStarted := ordered.Len() - inProgress.Len()
			return buildings.Len()+notStarted < quantity
		},

		Execute: func(b *Bot) {
			b.build(name, buildingId, abilityId, scl.S5x3)
		},

		Next: func(b *Bot) bool {
			buildings := b.Units.My.OfType(buildingId)
			ordered := b.FindWorkers().Filter(filter.IsOrderedTo(abilityId))
			inProgress := buildings.Filter(filter.IsInProgress)
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

			if quantity == 0 {
				return true
			}

			refineries := b.Units.My.OfType(terran.Refinery, terran.RefineryRich)
			ordered := b.FindWorkers().Filter(filter.IsOrderedTo(ability.Build_Refinery))
			inProgress := refineries.Filter(filter.IsInProgress)
			notStarted := ordered.Len() - inProgress.Len()

			return refineries.Len()+notStarted < quantity
		},

		Execute: func(b *Bot) {
			townHalls := b.FindTownHalls()
			if townHalls.Empty() {
				return
			}

			vespeneGeysers := b.findVespeneGeysersNearTownHalls(townHalls)
			claimed := b.FindClaimedVespeneGeysersNearTownHalls(townHalls)
			buildable := vespeneGeysers.Filter(filter.NotCloserThan(1, claimed))
			if buildable.Empty() {
				return
			}

			ordered := b.FindWorkers().Filter(filter.IsOrderedTo(ability.Build_Refinery))
			unplanned := buildable.Filter(filter.IsNotOrderedOnTag(ability.Build_Refinery, ordered))
			if unplanned.Empty() {
				return
			}

			random := unplanned[rand.Intn(len(unplanned))]
			worker := b.FindIdleOrGatheringWorkers().ClosestTo(random)
			if worker == nil {
				return
			}

			log.Info("Building refinery at %v", random.Point())
			worker.CommandTag(ability.Build_Refinery, random.Tag)
			worker.CommandTagQueue(ability.Smart, random.Tag)
			b.DeductResources(ability.Build_Refinery)
			b.Miners.GasForMiner[worker.Tag] = random.Tag
		},

		Next: func(b *Bot) bool {
			if quantity == 0 {
				return true
			}

			refineries := b.Units.My.OfType(terran.Refinery, terran.RefineryRich)
			ordered := b.FindWorkers().Filter(filter.IsOrderedTo(ability.Build_Refinery))
			inProgress := refineries.Filter(filter.IsInProgress)
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
			inProgress := b.Units.My.OfType(terran.CommandCenter).Filter(filter.IsOrderedTo(ability.Morph_OrbitalCommand))
			if orbitalCommands.Len()+inProgress.Len() >= quantity {
				b.State.CcForOrbitalCommand = 0
				return false
			}

			return true
		},

		Execute: func(b *Bot) {
			if b.State.CcForOrbitalCommand == 0 {
				// There's no command center marked for morphing into an orbital
				// command, so let's mark one
				commandCenters := b.Units.My.OfType(terran.CommandCenter).Filter(scl.Ready)
				if commandCenters.Empty() {
					return
				}

				randomCommandCenter := commandCenters[rand.Intn(len(commandCenters))]
				b.State.CcForOrbitalCommand = randomCommandCenter.Tag
			}

			// Check if the marked command center is still valid
			reserved := b.Units.ByTag[b.State.CcForOrbitalCommand]
			if reserved == nil || !reserved.Is(terran.CommandCenter, terran.CommandCenterFlying) {
				b.State.CcForOrbitalCommand = 0
				return
			}

			// If it's ordered to do anything else, cancel it
			ordered := filter.IsOrderedTo(ability.Morph_OrbitalCommand)(reserved)
			if len(reserved.Orders) > 0 && !ordered {
				reserved.Command(ability.Cancel_Last)
				return
			}

			// If it's not morphing yet, morph it
			if !ordered {
				log.Info("Morphing orbital command at %v", reserved.Point())
				reserved.Command(ability.Morph_OrbitalCommand)
				b.DeductResources(ability.Morph_OrbitalCommand)
				b.State.CcForOrbitalCommand = 0
			}
		},

		Next: func(b *Bot) bool {
			orbitalCommands := b.Units.My.OfType(terran.OrbitalCommand, terran.OrbitalCommandFlying)
			inProgress := b.Units.My.OfType(terran.CommandCenter).Filter(filter.IsOrderedTo(ability.Morph_OrbitalCommand))
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
			if b.State.BuildingForAddOn == 0 {
				randomBuilding := buildings[rand.Intn(len(buildings))]
				b.State.BuildingForAddOn = randomBuilding.Tag
			}

			// Check if the marked building is still valid
			reserved := b.Units.ByTag[b.State.BuildingForAddOn]
			if reserved == nil || !reserved.Is(buildingId) || reserved.AddOnTag != 0 {
				b.State.BuildingForAddOn = 0
				return
			}

			// If it's ordered to do anything else, cancel it
			ordered := filter.IsOrderedTo(abilityId)(reserved)
			if len(reserved.Orders) > 0 && !ordered {
				reserved.Command(ability.Cancel_Last)
				return
			}

			log.Info("Building %s at %v", name, reserved.Point())
			reserved.Command(abilityId)

			// In case it fails, queue the add-on to a new location
			elsewhere := b.WhereToBuild(reserved.Point(), scl.S5x3, addonId, abilityId)
			reserved.CommandPosQueue(abilityId, elsewhere)

			b.DeductResources(abilityId)
			b.State.BuildingForAddOn = 0
		},

		Next: func(b *Bot) bool {
			return b.Units.My.OfType(addonId).Len() >= quantity
		},
	}
}

var trainMarine = BuildStep{
	Name: "Train Marine",
	Predicate: func(b *Bot) bool {
		return b.CanBuy(ability.Train_Marine)
	},

	Execute: func(b *Bot) {
		barracks := b.Units.My.OfType(terran.Barracks).Filter(scl.Ready, scl.Ground, scl.Idle, filter.IsNotTag(b.State.BuildingForAddOn))
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
				log.Info("Training one marine at %v", barrack.Point())
			}

			if amount == 2 {
				barrack.CommandQueue(ability.Train_Marine)
				barrack.CommandQueue(ability.Train_Marine)
				log.Info("Training two marines at %v", barrack.Point())
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
	townHalls := b.FindTownHalls().Filter(filter.IsCcAtExpansion(b.State.CcForExp))
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

	townHalls := b.FindTownHalls().Filter(filter.IsCcAtExpansion(b.State.CcForExp))
	if townHalls.Empty() {
		return
	}

	randomTownHall := townHalls[rand.Intn(len(townHalls))]

	pos := b.WhereToBuild(randomTownHall.Point(), size, buildingId, abilityId)
	if pos == nil {
		return
	}

	builder := b.FindIdleOrGatheringWorkers().ClosestTo(pos)
	if builder == nil {
		return
	}

	if resource := b.findResourcesNearTownHalls(townHalls).ClosestTo(pos); resource != nil {
		log.Info("Building %s at %v and queuing to gather at %v", name, *pos, resource.Point())

		builder.CommandPos(abilityId, pos)
		builder.CommandTagQueue(ability.Smart, resource.Tag)

		if resource.IsMineral() {
			b.Miners.MineralForMiner[builder.Tag] = resource.Tag
		}

		if resource.IsGeyser() {
			b.Miners.GasForMiner[builder.Tag] = resource.Tag
		}
	} else {
		log.Info("Building %s at %v", name, *pos)
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

			if b.Units.My.OfType(buildingId).Filter(filter.IsOrderedTo(abilityId)).Exists() {
				return false
			}

			return true
		},

		Execute: func(b *Bot) {
			buildings := b.Units.My.OfType(buildingId).Filter(scl.Ready, scl.Idle)
			if buildings.Empty() {
				return
			}

			log.Info("Researching %s", name)
			buildings.First().Command(abilityId)
			b.DeductResources(abilityId)
		},

		Next: func(b *Bot) bool {
			if b.Upgrades[abilityId] {
				return true
			}

			if b.Units.My.OfType(buildingId).Filter(filter.IsOrderedTo(abilityId)).Exists() {
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

// defenseWaveStep assigns an attack wave to the defense of the bases.
var defenseWaveStep = BuildStep{
	Name: "Defense Wave",
	Predicate: func(b *Bot) bool {
		return true
	},
	Execute: func(b *Bot) {
		enemyInBases := b.FindEnemiesInBases()
		if len(enemyInBases) == 0 {
			return
		}

		base, enemies := b.mostThreatenedBase(enemyInBases)
		if base == nil || enemies.Empty() {
			return
		}

		// Where is the enemy cluster at that base?
		cluster := findClusterAtBase(base, enemies)
		if cluster.Empty() {
			log.Error("No cluster found at base %v", base.Point())
			return
		}

		inWaves := b.State.AttackWaves.Units(b)
		marines := b.Units.My.OfType(terran.Marine).Filter(scl.Ready, filter.NotIn(inWaves))
		if marines.Empty() {
			return
		}

		wave := AttackWave{
			Tags:   marines.Tags(),
			Target: cluster.Center(),
		}
		b.State.AttackWaves = append(b.State.AttackWaves, wave)
		log.Info("Sending %d marines defend base at %v", marines.Len(), base.Point())
	},
	Next: func(b *Bot) bool {
		return true
	},
}

func (b *Bot) mostThreatenedBase(enemyInBases map[api.UnitTag]scl.Units) (*scl.Unit, scl.Units) {
	var mostEnemyBase *scl.Unit
	var mostEnemyUnits scl.Units

	for base, units := range enemyInBases {
		if len(units) > mostEnemyUnits.Len() {
			mostEnemyBase = b.Units.MyAll.ByTag(base)
			mostEnemyUnits = units
		}
	}

	return mostEnemyBase, mostEnemyUnits
}

func findClusterAtBase(base *scl.Unit, enemies scl.Units) scl.Units {
	closest := enemies.ClosestTo(base)
	cluster := clusterBySight(closest, enemies)
	return cluster
}

func clusterBySight(target *scl.Unit, units scl.Units) scl.Units {
	cluster := make(scl.Units, 0, len(units))
	cluster.Add(target)

	// For each unit in the cluster, add all units in sight to the cluster.
	for _, unit := range cluster {
		inSight := units.Filter(filter.NotIn(cluster), filter.InSightOf(unit))
		cluster = append(cluster, inSight...)
	}

	return cluster
}

func (b *Bot) FindEnemyClusterAtHome() scl.Units {
	enemyInBases := b.FindEnemiesInBases()
	if len(enemyInBases) == 0 {
		return scl.Units{}
	}

	base, enemies := b.mostThreatenedBase(enemyInBases)
	if base == nil || enemies.Empty() {
		return scl.Units{}
	}

	// Where is the enemy cluster at that base?
	cluster := findClusterAtBase(base, enemies)
	if cluster.Empty() {
		log.Error("No cluster found at base %v", base.Point())
		return scl.Units{}
	}

	return cluster
}

func townHallStep(buildingId api.UnitTypeID, abilityId api.AbilityID, quantity int) *BuildStep {
	return &BuildStep{
		Name: "Town Hall",
		Predicate: func(b *Bot) bool {
			if !b.CanBuy(abilityId) {
				return false
			}

			townHalls := b.FindTownHalls()
			inProgress := townHalls.Filter(filter.IsInProgress)
			notStarted := b.FindWorkers().Filter(filter.IsOrderedTo(abilityId)).Len() - inProgress.Len()

			if quantity == 0 && inProgress.Len() == 0 && notStarted == 0 {
				return true
			}

			return townHalls.Len()+notStarted < quantity
		},
		Execute: func(b *Bot) {
			b.build("Town Hall", buildingId, abilityId, scl.S5x3)
		},
		Next: func(b *Bot) bool {
			if quantity == 0 {
				return true
			}

			townHalls := b.FindTownHalls()
			inProgress := townHalls.Filter(filter.IsInProgress)
			notStarted := b.FindWorkers().Filter(filter.IsOrderedTo(abilityId)).Len() - inProgress.Len()
			return townHalls.Len()+notStarted >= quantity
		},
	}
}
