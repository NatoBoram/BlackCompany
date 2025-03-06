package main

import (
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/api"
)

// IsEmptyMineral checks if a mineral field is depleted
func IsEmptyMineral(u *scl.Unit) bool {
	if u.IsMineral() {
		return u.MineralContents == 0
	}
	return false
}

// HasMinerals checks if a mineral field has minerals
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

// IsUnsaturatedMineralField returns a filter function that filters mineral
// fields with saturation below a target level
func IsUnsaturatedMineralField(saturation map[api.UnitTag]int, target int) scl.Filter {
	return func(mineralField *scl.Unit) bool {
		return saturation[mineralField.Tag] < target
	}
}

// IsUnsaturatedVespeneGeyser returns a filter function that filters vespene
// geysers with saturation below a target level
func IsUnsaturatedVespeneGeyser(saturation map[api.UnitTag]int, target int) scl.Filter {
	return func(refinery *scl.Unit) bool {
		return saturation[refinery.Tag] < target
	}
}

// CloserThan returns a filter function that selects units that are close to at
// least one unit in the list
func CloserThan(distance float64, units scl.Units) scl.Filter {
	return func(unit *scl.Unit) bool {
		for _, resource := range units {
			if unit.IsCloserThan(distance, resource) {
				return true
			}
		}
		return false
	}
}

// IsInProgress filters units that are currently being built
func IsInProgress(unit *scl.Unit) bool {
	return unit.BuildProgress < 1
}

// IsGatheringOrIdle filters units that are gathering resources or idle
func IsGatheringOrIdle(unit *scl.Unit) bool {
	if len(unit.Orders) <= 0 {
		return true
	}
	return unit.IsGathering()
}

// IsOrderedTo filters units that are currently ordered to use a specific
// ability.
func IsOrderedTo(ability api.AbilityID) scl.Filter {
	return func(u *scl.Unit) bool {
		if len(u.Orders) <= 0 {
			return false
		}

		for _, order := range u.Orders {
			if order.AbilityId == ability {
				return true
			}
		}

		return false
	}
}

func IsOrderedToAny(abilities []api.AbilityID) scl.Filter {
	return func(u *scl.Unit) bool {
		if len(u.Orders) <= 0 {
			return false
		}

		for _, ability := range abilities {
			if IsOrderedTo(ability)(u) {
				return true
			}
		}

		return false
	}
}

// HasTargetTag filters units that are targeting a specific unit
func HasTargetTag(tag api.UnitTag) scl.Filter {
	return func(u *scl.Unit) bool {
		if u.TargetTag() == tag {
			return true
		}

		for _, order := range u.Orders {
			if order.GetTargetUnitTag() == tag {
				return true
			}
		}

		return false
	}
}

// HasAnyTargetTag filters units that are targeting any of the specified units
func HasAnyTargetTag(tags scl.Tags) scl.Filter {
	return func(u *scl.Unit) bool {
		for _, tag := range tags {
			if u.TargetTag() == tag {
				return true
			}

			for _, order := range u.Orders {
				if order.GetTargetUnitTag() == tag {
					return true
				}
			}
		}

		return false
	}
}

// IsNotGathering filters units that are not gathering resources
func IsNotGathering(u *scl.Unit) bool {
	return !IsGathering(u)
}

// IsNotReturning filters units that are not returning resources
func IsNotReturning(u *scl.Unit) bool {
	return !IsReturning(u)
}

// HasTargetAbility filters units that are using a specific ability
func HasTargetAbility(ability api.AbilityID) scl.Filter {
	return func(u *scl.Unit) bool {
		if u.TargetAbility() == ability {
			return true
		}

		return false
	}
}

// HasAnyTargetAbility filters units that are using any of the specified
// abilities
func HasAnyTargetAbility(abilities []api.AbilityID) scl.Filter {
	return func(u *scl.Unit) bool {
		for _, ability := range abilities {
			if HasTargetAbility(ability)(u) {
				return true
			}
		}

		return false
	}
}

// IsGathering filters units that are gathering resources
func IsGathering(u *scl.Unit) bool {
	if u.IsGathering() {
		return true
	}

	gathering := ToKeys[api.AbilityID](scl.GatheringAbilities)
	return IsOrderedToAny(gathering)(u) || HasAnyTargetAbility(gathering)(u)
}

// IsReturning filters units that are returning resources
func IsReturning(u *scl.Unit) bool {
	if u.IsReturning() {
		return true
	}

	returning := ToKeys[api.AbilityID](scl.ReturningAbilities)
	return IsOrderedToAny(returning)(u) || HasAnyTargetAbility(returning)(u)
}

func IsGatheringOrReturning(u *scl.Unit) bool {
	if u.IsGathering() || u.IsReturning() {
		return true
	}

	gathering := ToKeys[api.AbilityID](scl.GatheringAbilities)
	returning := ToKeys[api.AbilityID](scl.ReturningAbilities)

	abilities := make([]api.AbilityID, 0, len(gathering)+len(returning))
	abilities = append(abilities, gathering...)
	abilities = append(abilities, returning...)

	return IsOrderedToAny(abilities)(u) || HasAnyTargetAbility(abilities)(u)
}
