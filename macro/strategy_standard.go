package macro

import (
	"github.com/NatoBoram/BlackCompany/bot"
	"github.com/aiseeq/s2l/protocol/enums/ability"
	"github.com/aiseeq/s2l/protocol/enums/terran"
)

var Standard = bot.Strategy{
	Name: "Standard",
	Steps: bot.BuildOrder{
		&defenseWaveStep,
		&supplyDepotStep,
		chatVersionStep(),
		buildingStep("Barracks", terran.Barracks, ability.Build_Barracks, 1, terran.SupplyDepot),
		refineryStep(1),
		orbitalCommandStep(1),
		addonStep("Barracks Reactor", terran.Barracks, terran.BarracksReactor, ability.Build_Reactor_Barracks, 1),
		expandStep(2),
		&marineStep,
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
		&planetaryFortressStep,

		// These are just in the meantime
		buildingStep("Armory", terran.Armory, ability.Build_Armory, 1, terran.Factory),

		upgradeStep("Infantry Weapons Level 2", ability.Research_TerranInfantryWeaponsLevel2, terran.EngineeringBay),
		upgradeStep("Infantry Armor Level 2", ability.Research_TerranInfantryArmorLevel2, terran.EngineeringBay),

		upgradeStep("Infantry Weapons Level 3", ability.Research_TerranInfantryWeaponsLevel3, terran.EngineeringBay),
		upgradeStep("Infantry Armor Level 3", ability.Research_TerranInfantryArmorLevel3, terran.EngineeringBay),
	},
}
