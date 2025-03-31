package filter

import (
	"github.com/aiseeq/s2l/lib/point"
	"github.com/aiseeq/s2l/lib/scl"
)

func IsNotAt(point point.Point) scl.Filter {
	return func(u *scl.Unit) bool {
		return u.IsFurtherThan(1, point)
	}
}
