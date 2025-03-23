package bot

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

// BuildOrder is a list of build steps in a strategy.
type BuildOrder []*BuildStep

// Strategy is a set of build steps to be executed.
type Strategy struct {
	Name  string
	Steps BuildOrder
}

// ExecuteStrategy executes a strategy.
func (b *Bot) ExecuteStrategy(s *Strategy) {
	for _, step := range s.Steps {
		if step.Predicate(b) {
			step.Execute(b)
		}

		if !step.Next(b) {
			break
		}
	}
}
