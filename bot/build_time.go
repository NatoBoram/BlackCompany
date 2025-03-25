package bot

import "math"

type BuildTime uint

const (
	BuildTimeCommandCenter BuildTime = 71
	BuildTimeMarine        BuildTime = 25
	BuildTimeSCV           BuildTime = 12
	BuildTimeSupplyDepot   BuildTime = 40
)

// BuildDuring calculates the amount of X you can build during the production of
// one Y.
func BuildDuring(x, y BuildTime) int {
	return int(math.Ceil(float64(y) / float64(x)))
}
