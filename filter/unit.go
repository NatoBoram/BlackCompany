package filter

import (
	"math"
	"slices"

	"github.com/NatoBoram/BlackCompany/adapter"
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/api"
	"github.com/aiseeq/s2l/protocol/enums/ability"
)

// IsEmptyMineral checks if a mineral field is depleted
func IsEmptyMineral(u *scl.Unit) bool {
	if u.IsMineral() {
		return u.MineralContents == 0
	}
	return false
}

// HasMinerals checks if a mineral field has minerals.
//
// Debug notes:
//
// Minerals that have been seen before:
// display_type:Visible alliance:Neutral tag:4310171649 unit_type:666 owner:16
// pos:<x:138 y:20.5 z:11.989014> facing:4.712389 radius:1.125
// build_progress:1 cloak:NotCloaked health:500 health_max:500
// mineral_contents:860
//
// Minerals that have never been seen before:
// display_type:Snapshot alliance:Neutral tag:8754364417 unit_type:665
// owner:16 pos:<x:23 y:140.5 z:11.989014> radius:1.125 build_progress:1
//
// So check if it has been explored before relying on this function.
func HasMinerals(u *scl.Unit) bool {
	if u.IsMineral() {
		return u.MineralContents > 0
	}
	return false
}

// IsStructure checks if a unit is a structure
func IsStructure(u *scl.Unit) bool {
	return u.IsStructure()
}

// HasGas checks if a vespene geyser or refinery has vespene gas remaining
func HasGas(u *scl.Unit) bool {
	return u.VespeneContents > 0
}

// IsNotNil checks if a unit is not nil
func IsNotNil(u *scl.Unit) bool {
	return u != nil
}

// IsInProgress filters units that are currently being built
func IsInProgress(unit *scl.Unit) bool {
	return unit.BuildProgress < 1
}

// IsNotGathering filters units that are not gathering resources
func IsNotGathering(u *scl.Unit) bool {
	return !IsGathering(u)
}

// IsNotReturning filters units that are not returning resources
func IsNotReturning(u *scl.Unit) bool {
	return !IsReturning(u)
}

// IsGathering filters units that are gathering resources
func IsGathering(u *scl.Unit) bool {
	if u.IsGathering() {
		return true
	}

	gathering := adapter.ToKeys(scl.GatheringAbilities)
	return IsOrderedToAny(gathering...)(u) || HasAnyTargetAbility(gathering)(u)
}

// IsReturning filters units that are returning resources
func IsReturning(u *scl.Unit) bool {
	if u.IsReturning() {
		return true
	}

	returning := adapter.ToKeys(scl.ReturningAbilities)
	return IsOrderedToAny(returning...)(u) || HasAnyTargetAbility(returning)(u)
}

func IsGatheringOrReturning(u *scl.Unit) bool {
	if u.IsGathering() || u.IsReturning() {
		return true
	}

	gathering := adapter.ToKeys(scl.GatheringAbilities)
	returning := adapter.ToKeys(scl.ReturningAbilities)
	abilities := slices.Concat(gathering, returning)

	return IsOrderedToAny(abilities...)(u) || HasAnyTargetAbility(abilities)(u)
}

// IsNotBuilding filters units that are not currently ordered to build a
// structure
func IsNotBuilding(u *scl.Unit) bool {
	if len(u.Orders) <= 0 {
		return true
	}

	buildingAbilities := []api.AbilityID{
		ability.Build_Armory,
		ability.Build_Assimilator,
		ability.Build_BanelingNest,
		ability.Build_Barracks,
		ability.Build_Bunker,
		ability.Build_CommandCenter,
		ability.Build_CreepTumor,
		ability.Build_CreepTumor_Queen,
		ability.Build_CreepTumor_Tumor,
		ability.Build_CyberneticsCore,
		ability.Build_DarkShrine,
		ability.Build_EngineeringBay,
		ability.Build_EvolutionChamber,
		ability.Build_Extractor,
		ability.Build_Factory,
		ability.Build_FleetBeacon,
		ability.Build_Forge,
		ability.Build_FusionCore,
		ability.Build_Gateway,
		ability.Build_GhostAcademy,
		ability.Build_Hatchery,
		ability.Build_HydraliskDen,
		ability.Build_InfestationPit,
		ability.Build_Interceptors,
		ability.Build_LurkerDen,
		ability.Build_MissileTurret,
		ability.Build_Nexus,
		ability.Build_Nuke,
		ability.Build_NydusNetwork,
		ability.Build_NydusWorm,
		ability.Build_PhotonCannon,
		ability.Build_Pylon,
		ability.Build_Reactor,
		ability.Build_Reactor_Barracks,
		ability.Build_Reactor_Factory,
		ability.Build_Reactor_Starport,
		ability.Build_Refinery,
		ability.Build_RoachWarren,
		ability.Build_RoboticsBay,
		ability.Build_RoboticsFacility,
		ability.Build_SensorTower,
		ability.Build_ShieldBattery,
		ability.Build_SpawningPool,
		ability.Build_SpineCrawler,
		ability.Build_Spire,
		ability.Build_SporeCrawler,
		ability.Build_Stargate,
		ability.Build_Starport,
		ability.Build_StasisTrap,
		ability.Build_SupplyDepot,
		ability.Build_TechLab,
		ability.Build_TechLab_Barracks,
		ability.Build_TechLab_Factory,
		ability.Build_TechLab_Starport,
		ability.Build_TemplarArchive,
		ability.Build_TwilightCouncil,
		ability.Build_UltraliskCavern,
	}

	return !IsOrderedToAny(buildingAbilities...)(u)
}

func SameHeightAs(u *scl.Unit) scl.Filter {
	return func(u2 *scl.Unit) bool {
		return math.Abs(float64(u.Pos.Z-u2.Pos.Z)) < 1
	}
}

func InSightOf(target *scl.Unit) scl.Filter {
	return func(unit *scl.Unit) bool {
		return unit.IsCloserThan(target.SightRange(), target)
	}
}
