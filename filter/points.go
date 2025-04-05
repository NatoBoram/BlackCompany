package filter

import (
	"github.com/aiseeq/s2l/lib/point"
	"github.com/aiseeq/s2l/lib/scl"
)

func IsNotAtAny(points point.Points) scl.Filter {
	return func(u *scl.Unit) bool {
		return points.CloserThan(1, u).Empty()
	}
}

func NotInPoints(points point.Points) point.Filter {
	return func(p point.Point) bool {
		return points.CloserThan(1, p).Empty()
	}
}
