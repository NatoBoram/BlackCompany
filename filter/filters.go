package filter

import (
	"github.com/aiseeq/s2l/lib/point"
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/api"
)

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
		return units.CloserThan(distance, unit).Exists()
	}
}

// NotCloserThan returns a filter that selects units that are not closer than
// the specified distance to any unit in the given set.
func NotCloserThan(distance float64, units scl.Units) scl.Filter {
	return func(unit *scl.Unit) bool {
		return units.CloserThan(distance, unit).Empty()
	}
}

// IsOrderedToTarget filters units that are currently ordered to use a specific
// ability to a specific coordinate.
func IsOrderedToTarget(ability api.AbilityID, target point.Point) scl.Filter {
	return func(u *scl.Unit) bool {
		if len(u.Orders) <= 0 {
			return false
		}

		for _, order := range u.Orders {
			orderPos := order.GetTargetWorldSpacePos()
			orderPoint := point.Pt3(orderPos)

			if order.AbilityId == ability && orderPoint.Dist(target) < 1 {
				return true
			}
		}

		return false
	}
}

func IsOrderedToTag(ability api.AbilityID, tag api.UnitTag) scl.Filter {
	return func(u *scl.Unit) bool {
		if len(u.Orders) <= 0 {
			return false
		}

		for _, order := range u.Orders {
			if order.AbilityId == ability && order.GetTargetUnitTag() == tag {
				return true
			}
		}

		return false
	}
}

// IsNotOrderedOnTag filters units that aren't being ordered to something by
// other units.
//
// For example, to check that no workers in that list are ordered to build a
// refinery on that vespene geyser, you would do:
//
//	geysers.Filter(filter.IsNotOrderedOnTag(ability.Build_Refinery, workers))
func IsNotOrderedOnTag(ability api.AbilityID, units scl.Units) scl.Filter {
	return func(target *scl.Unit) bool {
		for _, unit := range units {
			if len(unit.Orders) <= 0 {
				continue
			}

			for _, order := range unit.Orders {
				if order.AbilityId == ability && order.GetTargetUnitTag() == target.Tag {
					return false
				}
			}
		}

		return true
	}
}

// IsNotOrderedToTarget filters units that are not currently ordered to use a
// specific ability to a specific coordinate.
func IsNotOrderedToTarget(ability api.AbilityID, target point.Point) scl.Filter {
	return func(u *scl.Unit) bool {
		return !IsOrderedToTarget(ability, target)(u)
	}
}

// IsCcAtExpansion checks if a command center is at its designated expansion
// location. Useful to exclude bases that need to be flown to their expansion.
func IsCcAtExpansion(ccForExp map[api.UnitTag]point.Point) scl.Filter {
	return func(u *scl.Unit) bool {
		expansion, ok := ccForExp[u.Tag]
		return ok && expansion.Dist(u) < 1
	}
}

// IsNotCcAtExpansion filters command centers that are not at their designated
// expansion location.
func IsNotCcAtExpansion(ccForExp map[api.UnitTag]point.Point) scl.Filter {
	return func(u *scl.Unit) bool {
		expansion, ok := ccForExp[u.Tag]
		return !ok || expansion.Dist(u) >= 1
	}
}
