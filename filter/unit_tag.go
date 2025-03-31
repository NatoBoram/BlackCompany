package filter

import (
	"slices"

	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/api"
)

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

func IsNotTag(tag api.UnitTag) scl.Filter {
	return func(u *scl.Unit) bool {
		return u.Tag != tag
	}
}

func IsNotTags(tags ...api.UnitTag) scl.Filter {
	return func(u *scl.Unit) bool {
		return !slices.Contains(tags, u.Tag)
	}
}
