package bot

// buildDuring calculates how many times a fast action can be done during a
// slower action.
func buildDuring(fast, slow float64) int {
	return int(slow / fast)
}
