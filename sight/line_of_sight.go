package sight

// LineOfSight represents the area around a unit or building that can be seen by
// it.
//
// https://liquipedia.net/starcraft2/Sight
type LineOfSight float64

const (
	// Unless otherwise noted, buildings have a sight range of 9.
	LineOfSightBuildings LineOfSight = 9

	// The Xel'Naga Tower provides a Sight range of 22.
	LineOfSightXelNagaTower LineOfSight = 22
)

func (los LineOfSight) Float64() float64 {
	return float64(los)
}
