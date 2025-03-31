package filter

import "github.com/aiseeq/s2l/lib/scl"

// NotIn returns units that are not in the list.
func NotIn(units scl.Units) scl.Filter {
	return func(u *scl.Unit) bool {
		return units.ByTag(u.Tag) == nil
	}
}
