package main

import (
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/enums/protoss"
	"github.com/aiseeq/s2l/protocol/enums/terran"
	"github.com/aiseeq/s2l/protocol/enums/zerg"
)

// findMineralFieldsNearTownHalls finds all mineral fields near town halls.
func (b *Bot) findMineralFieldsNearTownHalls(townHalls scl.Units) scl.Units {
	mineralFields := make(scl.Units, 0, len(townHalls)*8)
	for _, th := range townHalls {
		ccMineralFields := b.Units.Minerals.All().CloserThan(th.SightRange(), th).Filter(HasMinerals)
		mineralFields = append(mineralFields, ccMineralFields...)
	}

	return mineralFields
}

// findClaimedVespeneGeysersNearTownHalls finds all vespene geysers near town
// halls.
func (b *Bot) findClaimedVespeneGeysersNearTownHalls(townHalls scl.Units) scl.Units {
	refineries := make(scl.Units, 0, len(townHalls)*2)
	for _, th := range townHalls {
		ccRefineries := b.Units.My.OfType(
			protoss.Assimilator, protoss.AssimilatorRich,
			terran.Refinery, terran.RefineryRich,
			zerg.Extractor, zerg.ExtractorRich,
		).CloserThan(th.SightRange(), th).Filter(HasGas)

		refineries = append(refineries, ccRefineries...)
	}

	return refineries
}

// findVespeneGeysersNearTownHalls finds all vespene geysers near town halls.
func (b *Bot) findVespeneGeysersNearTownHalls(townHalls scl.Units) scl.Units {
	vespeneGeysers := make(scl.Units, 0, len(townHalls)*2)
	for _, th := range townHalls {
		ccGeysers := b.Units.Geysers.All().
			CloserThan(th.SightRange(), th).Filter(HasGas)

		vespeneGeysers = append(vespeneGeysers, ccGeysers...)
	}

	return vespeneGeysers
}

// findTownHalls finds all structures capable of training workers and miners.
func (b *Bot) findTownHalls() scl.Units {
	return b.Units.My.OfType(
		protoss.Nexus,
		terran.CommandCenter, terran.OrbitalCommand, terran.PlanetaryFortress,
		zerg.Hatchery, zerg.Lair, zerg.Hive,
	)
}

// findProductionStructures finds all structures capable of training military
// units.
func (b *Bot) findProductionStructures() scl.Units {
	return b.Units.My.OfType(
		protoss.Gateway, protoss.RoboticsFacility, protoss.Stargate,
		terran.Barracks, terran.Factory, terran.Starport,
		zerg.Hatchery, zerg.Lair, zerg.Hive,
	)
}

// findMiners finds all units capable of mining resources.
func (b *Bot) findMiners() scl.Units {
	return b.Units.My.OfType(protoss.Probe, terran.MULE, terran.SCV, zerg.Drone)
}

// findWorkers finds all units capable of building structures.
func (b *Bot) findWorkers() scl.Units {
	return b.Units.My.OfType(protoss.Probe, terran.SCV, zerg.Drone)
}

// findIdleOrGatheringWorkers finds idle or gathering workers that are not
// currently building a structure.
func (b *Bot) findIdleOrGatheringWorkers() scl.Units {
	workers := b.findWorkers().Filter(IsNotBuilding)

	if idle := workers.Filter(scl.Idle); !idle.Empty() {
		return idle
	}

	// If they're gathering, then it means they're not carrying resources
	if gathering := workers.Filter(IsGathering); !gathering.Empty() {
		return gathering
	}

	return workers.Filter(IsReturning)
}

// findTurretsNearResourcesNearTownHalls finds all missile turrets near mineral
// fields and vespene geysers near town halls.
func (b *Bot) findTurretsNearResourcesNearTownHalls(resources scl.Units) scl.Units {
	return b.Units.My.OfType(
		protoss.PhotonCannon,
		terran.AutoTurret, terran.MissileTurret,
		zerg.SpineCrawler, zerg.SporeCrawler,
	).Filter(CloserThan(scl.ResourceSpreadDistance, resources))
}

// findResourcesNearTownHalls finds all resources near town halls.
func (b *Bot) findResourcesNearTownHalls(townHalls scl.Units) scl.Units {
	mineralFields := b.findMineralFieldsNearTownHalls(townHalls)
	refineries := b.findClaimedVespeneGeysersNearTownHalls(townHalls)

	resources := make(scl.Units, 0, len(mineralFields)+len(refineries))
	resources = append(resources, mineralFields...)
	resources = append(resources, refineries...)

	return resources
}

// findUnsaturatedMineralFieldsNearTownHalls finds all unsaturated mineral
// fields near town halls
func (b *Bot) findUnsaturatedMineralFieldsNearTownHalls(townHalls scl.Units) scl.Units {
	mineralFields := b.findMineralFieldsNearTownHalls(townHalls)
	saturation := b.GetMineralsSaturation(mineralFields)
	return mineralFields.Filter(IsUnsaturatedMineralField(saturation, 2))
}

// findUnsaturatedVespeneGeysersNearTownHalls finds all unsaturated vespene
// geysers near town halls
func (b *Bot) findUnsaturatedVespeneGeysersNearTownHalls(townHalls scl.Units) scl.Units {
	refineries := b.findClaimedVespeneGeysersNearTownHalls(townHalls)
	saturation := b.GetGasSaturation(refineries)
	return refineries.Filter(IsUnsaturatedVespeneGeyser(saturation, 3))
}

// findUnsaturatedResourcesNearTownHalls finds all unsaturated resources near
// town halls.
func (b *Bot) findUnsaturatedResourcesNearTownHalls(townHalls scl.Units) scl.Units {
	mineralFields := b.findUnsaturatedMineralFieldsNearTownHalls(townHalls)
	vespeneGeysers := b.findUnsaturatedVespeneGeysersNearTownHalls(townHalls)

	resources := make(scl.Units, 0, len(mineralFields)+len(vespeneGeysers))
	resources = append(resources, mineralFields...)
	resources = append(resources, vespeneGeysers...)

	return resources
}
