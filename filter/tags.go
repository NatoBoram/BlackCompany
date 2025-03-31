package filter

import "github.com/aiseeq/s2l/lib/scl"

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
