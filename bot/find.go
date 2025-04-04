package bot

import (
	"slices"

	"github.com/NatoBoram/BlackCompany/filter"
	"github.com/NatoBoram/BlackCompany/sight"
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/api"
	"github.com/aiseeq/s2l/protocol/enums/protoss"
	"github.com/aiseeq/s2l/protocol/enums/terran"
	"github.com/aiseeq/s2l/protocol/enums/zerg"
)

// findMineralFieldsNearTownHalls finds all mineral fields near town halls.
func (b *Bot) FindMineralFieldsNearTownHalls(townHalls scl.Units) scl.Units {
	mineralFields := make(scl.Units, 0, len(townHalls)*8)
	for _, th := range townHalls {
		ccMineralFields := b.Units.Minerals.All().CloserThan(th.SightRange(), th).Filter(filter.HasMinerals)
		mineralFields = append(mineralFields, ccMineralFields...)
	}

	return mineralFields
}

// findClaimedVespeneGeysersNearTownHalls finds all vespene geysers near town
// halls.
func (b *Bot) FindClaimedVespeneGeysersNearTownHalls(townHalls scl.Units) scl.Units {
	refineries := make(scl.Units, 0, len(townHalls)*2)
	for _, th := range townHalls {
		ccRefineries := b.Units.My.OfType(
			protoss.Assimilator, protoss.AssimilatorRich,
			terran.Refinery, terran.RefineryRich,
			zerg.Extractor, zerg.ExtractorRich,
		).CloserThan(th.SightRange(), th).Filter(filter.HasGas)

		refineries = append(refineries, ccRefineries...)
	}

	return refineries
}

// FindVespeneGeysersNearTownHalls finds all vespene geysers near town halls.
func (b *Bot) FindVespeneGeysersNearTownHalls(townHalls scl.Units) scl.Units {
	vespeneGeysers := make(scl.Units, 0, len(townHalls)*2)
	for _, th := range townHalls {
		ccGeysers := b.Units.Geysers.All().
			CloserThan(th.SightRange(), th).Filter(filter.HasGas)

		vespeneGeysers = append(vespeneGeysers, ccGeysers...)
	}

	return vespeneGeysers
}

// findTownHalls finds all structures capable of training workers and miners.
func (b *Bot) FindTownHalls() scl.Units {
	return b.Units.My.OfType(
		protoss.Nexus,
		terran.CommandCenter, terran.OrbitalCommand, terran.PlanetaryFortress,
		zerg.Hatchery, zerg.Lair, zerg.Hive,
	)
}

// FindProductionStructures finds all structures capable of training military
// units.
func (b *Bot) FindProductionStructures() scl.Units {
	return b.Units.My.OfType(
		protoss.Gateway, protoss.RoboticsFacility, protoss.Stargate,
		terran.Barracks, terran.Factory, terran.Starport,
		zerg.Hatchery, zerg.Lair, zerg.Hive,
	)
}

// findMiners finds all units capable of mining resources.
func (b *Bot) FindMiners() scl.Units {
	return b.Units.My.OfType(protoss.Probe, terran.MULE, terran.SCV, zerg.Drone)
}

// findWorkers finds all units capable of building structures.
func (b *Bot) FindWorkers() scl.Units {
	return b.Units.My.OfType(protoss.Probe, terran.SCV, zerg.Drone)
}

// findIdleOrGatheringWorkers finds idle or gathering workers that are not
// currently building a structure.
func (b *Bot) FindIdleOrGatheringWorkers() scl.Units {
	workers := b.FindWorkers().Filter(filter.IsNotBuilding)

	if idle := workers.Filter(scl.Idle); !idle.Empty() {
		return idle
	}

	// If they're gathering, then it means they're not carrying resources
	if gathering := workers.Filter(filter.IsGathering); !gathering.Empty() {
		return gathering
	}

	return workers.Filter(filter.IsReturning)
}

// findTurretsNearResourcesNearTownHalls finds all missile turrets near mineral
// fields and vespene geysers near town halls.
func (b *Bot) findTurretsNearResourcesNearTownHalls(resources scl.Units) scl.Units {
	return b.Units.My.OfType(
		protoss.PhotonCannon,
		terran.AutoTurret, terran.MissileTurret,
		zerg.SpineCrawler, zerg.SporeCrawler,
	).Filter(filter.CloserThan(scl.ResourceSpreadDistance, resources))
}

// FindResourcesNearTownHalls finds all resources near town halls.
func (b *Bot) FindResourcesNearTownHalls(townHalls scl.Units) scl.Units {
	mineralFields := b.FindMineralFieldsNearTownHalls(townHalls)
	refineries := b.FindClaimedVespeneGeysersNearTownHalls(townHalls)
	return slices.Concat(mineralFields, refineries)
}

// findUnsaturatedMineralFieldsNearTownHalls finds all unsaturated mineral
// fields near town halls
func (b *Bot) FindUnsaturatedMineralFieldsNearTownHalls(townHalls scl.Units) scl.Units {
	mineralFields := b.FindMineralFieldsNearTownHalls(townHalls)
	saturation := b.GetMineralsSaturation(mineralFields)
	return mineralFields.Filter(filter.IsUnsaturatedMineralField(saturation, 2))
}

// findUnsaturatedVespeneGeysersNearTownHalls finds all unsaturated vespene
// geysers near town halls
func (b *Bot) FindUnsaturatedVespeneGeysersNearTownHalls(townHalls scl.Units) scl.Units {
	refineries := b.FindClaimedVespeneGeysersNearTownHalls(townHalls)
	saturation := b.GetGasSaturation(refineries)
	return refineries.Filter(filter.IsUnsaturatedVespeneGeyser(saturation, 3))
}

// findUnsaturatedResourcesNearTownHalls finds all unsaturated resources near
// town halls.
func (b *Bot) FindUnsaturatedResourcesNearTownHalls(townHalls scl.Units) scl.Units {
	mineralFields := b.FindUnsaturatedMineralFieldsNearTownHalls(townHalls)
	vespeneGeysers := b.FindUnsaturatedVespeneGeysersNearTownHalls(townHalls)

	return slices.Concat(mineralFields, vespeneGeysers)
}

// FindEnemiesInBases finds enemies in range of bases. Bases are all buildings
// in sight of eachother, starting with the town hall.
func (b *Bot) FindEnemiesInBases() map[api.UnitTag]scl.Units {
	bases := b.FindTownHalls().Filter(filter.IsCcAtExpansion(b.State.CcForExp))
	enemiesInBases := make(map[api.UnitTag]scl.Units, len(bases))

	for _, base := range bases {
		buildings := scl.Units{base}
		enemies := b.Units.Enemy.All().CloserThan(base.SightRange(), base)

		for _, building := range buildings {
			buildingsInSight := b.Units.MyAll.
				Filter(scl.Structure, filter.NotIn(buildings)).
				CloserThan(building.SightRange(), building)

			buildings = append(buildings, buildingsInSight...)

			enemiesInRange := b.Units.Enemy.All().CloserThan(sight.LineOfSightScannerSweep.Float64(), building)
			enemies = append(enemies, enemiesInRange...)
		}

		enemiesInBases[base.Tag] = enemies
	}

	return enemiesInBases
}

// FindEnemyAirArmy finds all enemy air units that are not workers.
func (b *Bot) FindEnemyAirArmy() scl.Units {
	return b.Units.Enemy.All().Filter(scl.DpsGt5, scl.Flying)
}
