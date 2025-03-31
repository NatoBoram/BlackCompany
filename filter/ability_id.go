package filter

import (
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/api"
)

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

func IsOrderedToAny(abilities ...api.AbilityID) scl.Filter {
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

func IsNotOrderedToAny(abilities ...api.AbilityID) scl.Filter {
	return func(u *scl.Unit) bool {
		if len(u.Orders) <= 0 {
			return true
		}

		for _, ability := range abilities {
			if IsOrderedTo(ability)(u) {
				return false
			}
		}

		return true
	}
}

// HasTargetAbility filters units that are using a specific ability
func HasTargetAbility(ability api.AbilityID) scl.Filter {
	return func(u *scl.Unit) bool {
		return u.TargetAbility() == ability
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
