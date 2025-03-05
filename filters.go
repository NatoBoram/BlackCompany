package main

import "github.com/aiseeq/s2l/lib/scl"

// IsEmptyMineral checks if a mineral field is depleted
func IsEmptyMineral(u *scl.Unit) bool {
	if u.IsMineral() {
		return u.MineralContents == 0
	}
	return false
}

// HasMinerals checks if a mineral field has minerals
func HasMinerals(u *scl.Unit) bool {
	if u.IsMineral() {
		return u.MineralContents > 0
	}
	return false
}

func IsStructure(u *scl.Unit) bool {
	return u.IsStructure()
}

// HasGas checks if a vespene geyser or refinery has vespene gas remaining
func HasGas(u *scl.Unit) bool {
	return u.VespeneContents > 0
}

// IsNotNil checks if a unit is not nil
func IsNotNil(u *scl.Unit) bool {
	return u != nil
}
