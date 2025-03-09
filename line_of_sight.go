package main

// LineOfSight represents the area around a unit or building that can be seen by
// it.
//
// https://liquipedia.net/starcraft2/Sight
type LineOfSight float64

const (
	// Unless otherwise noted, buildings have a sight range of 9.
	LineOfSightBuildings LineOfSight = 9

	// The Xel'Naga Tower provides a Sight range of 22.
	LineOfSightXelNagaTower LineOfSight = 22
)

const (
	// Protos Units
	LineOfSightProbe                    LineOfSight = 8
	LineOfSightZealot                   LineOfSight = 9
	LineOfSightStalker                  LineOfSight = 10
	LineOfSightSentry                   LineOfSight = 10
	LineOfSightAdept                    LineOfSight = 9
	LineOfSightShade                    LineOfSight = 4
	LineOfSightHighTemplar              LineOfSight = 10
	LineOfSightDarkTemplar              LineOfSight = 8
	LineOfSightArchon                   LineOfSight = 9
	LineOfSightObserver                 LineOfSight = 11
	LineOfSightObserverSurveillanceMode LineOfSight = 13.75
	LineOfSightWarpPrism                LineOfSight = 10
	LineOfSightImmortal                 LineOfSight = 9
	LineOfSightColossus                 LineOfSight = 10
	LineOfSightDisruptor                LineOfSight = 9
	LineOfSightPhoenix                  LineOfSight = 10
	LineOfSightVoidRay                  LineOfSight = 10
	LineOfSightOracle                   LineOfSight = 10
	LineOfSightStasisWard               LineOfSight = 4
	LineOfSightTempest                  LineOfSight = 12
	LineOfSightCarrier                  LineOfSight = 12
	LineOfSightInterceptor              LineOfSight = 7
	LineOfSightMothership               LineOfSight = 14
)

// Protos Buildings
const (
	LineOfSightNexus        LineOfSight = 11
	LineOfSightPhotonCannon LineOfSight = 11
)

// Terran Units
const (
	LineOfSightSCV           LineOfSight = 8
	LineOfSightMULE          LineOfSight = 8
	LineOfSightMarine        LineOfSight = 9
	LineOfSightReaper        LineOfSight = 9
	LineOfSightMarauder      LineOfSight = 10
	LineOfSightGhost         LineOfSight = 11
	LineOfSightHellion       LineOfSight = 10
	LineOfSightHellbat       LineOfSight = 10
	LineOfSightWidowMine     LineOfSight = 7
	LineOfSightCyclone       LineOfSight = 11
	LineOfSightSiegeTank     LineOfSight = 11
	LineOfSightThor          LineOfSight = 11
	LineOfSightViking        LineOfSight = 10
	LineOfSightMedivac       LineOfSight = 11
	LineOfSightLiberator     LineOfSight = 10
	LineOfSightRaven         LineOfSight = 11
	LineOfSightAutoTurret    LineOfSight = 7
	LineOfSightBanshee       LineOfSight = 10
	LineOfSightBattlecruiser LineOfSight = 12
)

// Terran Buildings
const (
	LineOfSightCommandCenter     LineOfSight = 11
	LineOfSightPlanetaryFortress LineOfSight = 11
	LineOfSightOrbitalCommand    LineOfSight = 11
	LineOfSightScannerSweep      LineOfSight = 13
	LineOfSightBunker            LineOfSight = 10
	LineOfSightSensorTower       LineOfSight = 12
	LineOfSightRadar             LineOfSight = 30
	LineOfSightMissileTurret     LineOfSight = 11
)

// Zerg Units
const (
	LineOfSightLarva                 LineOfSight = 5
	LineOfSightDrone                 LineOfSight = 8
	LineOfSightOverlord              LineOfSight = 11
	LineOfSightQueen                 LineOfSight = 9
	LineOfSightZergling              LineOfSight = 8
	LineOfSightBaneling              LineOfSight = 8
	LineOfSightRoach                 LineOfSight = 9
	LineOfSightOverseer              LineOfSight = 11
	LineOfSightOverseerOversightMode LineOfSight = 13.75
	LineOfSightOverseerChangeling    LineOfSight = 8
	LineOfSightHydralisk             LineOfSight = 9
	LineOfSightLurker                LineOfSight = 10
	LineOfSightMutalisk              LineOfSight = 11
	LineOfSightCorruptor             LineOfSight = 10
	LineOfSightSwarmHost             LineOfSight = 10
	LineOfSightLocust                LineOfSight = 6
	LineOfSightInfestor              LineOfSight = 10
	LineOfSightInfestedTerran        LineOfSight = 9
	LineOfSightInfestedTerranEgg     LineOfSight = 0
	LineOfSightViper                 LineOfSight = 11
	LineOfSightUltralisk             LineOfSight = 9
	LineOfSightBroodLord             LineOfSight = 12
	LineOfSightBroodling             LineOfSight = 7
)

// Zerg Buildings
const (
	LineOfSightHatchery     LineOfSight = 12
	LineOfSightSpineCrawler LineOfSight = 11
	LineOfSightSporeCrawler LineOfSight = 11
	LineOfSightCreepTumor   LineOfSight = 11
	LineOfSightLair         LineOfSight = 12
	LineOfSightNydusWorm    LineOfSight = 10
	LineOfSightHive         LineOfSight = 12
)

func (los LineOfSight) Float64() float64 {
	return float64(los)
}
