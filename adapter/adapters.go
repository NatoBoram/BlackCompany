package adapter

import (
	"maps"

	"github.com/aiseeq/s2l/lib/point"
	"github.com/aiseeq/s2l/lib/scl"
)

// ToPoints converts a slice of units to a slice of points.
func ToPoints(units scl.Units) point.Points {
	points := make(point.Points, 0, len(units))
	for _, unit := range units {
		points = append(points, unit.Point())
	}
	return points
}

// ToTags converts a slice of units to a slice of tags.
func ToTags(unit scl.Units) scl.Tags {
	tags := make(scl.Tags, 0, len(unit))
	for _, u := range unit {
		tags = append(tags, u.Tag)
	}
	return tags
}

// ToKeys returns a slice containing all keys from the input map
func ToKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range maps.Keys(m) {
		keys = append(keys, k)
	}
	return keys
}
