package main

import "math"

type BuildTime uint

const (
	BuildTimeCommandCenter BuildTime = 71
	BuildTimeSCV           BuildTime = 12
	BuildTimeSupplyDepot   BuildTime = 40
)

// buildDuring calculates the amount of X you can build during the production of
// one Y.
func buildDuring(x, y BuildTime) int {
	return int(math.Ceil(float64(y) / float64(x)))
}
