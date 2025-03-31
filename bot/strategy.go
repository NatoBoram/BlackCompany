package bot

// Strategy is a build order to be executed.
type Strategy struct {
	Name  string
	Steps BuildOrder
}

// BuildOrder is a list of build steps in a strategy.
type BuildOrder []*BuildStep

// BuildStep is a step in a build order.
type BuildStep struct {
	// Name of the step.
	Name string

	// Predicate determines if this step should be executed.
	Predicate func(*Bot) bool

	// Execute is the action to be taken.
	Execute func(*Bot)

	// Next determines if we're ready to advance to the next step.
	Next func(*Bot) bool
}
