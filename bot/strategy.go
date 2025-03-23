package bot

// Strategy represents a build order strategy
type Strategy struct {
	Name  string
	Steps BuildOrder
}

// BuildOrder represents a sequence of build steps
type BuildOrder []*BuildStep

// BuildStep represents a single step in a build order
type BuildStep struct {
	Name string

	// Called to check if we should execute this step
	Predicate func(b *Bot) bool

	// Called to execute this step
	Execute func(b *Bot)

	// Called to check if we should move to the next step
	Next func(b *Bot) bool
}

// ExecuteStrategy executes the current strategy
func (b *Bot) ExecuteStrategy(strategy *Strategy) {
	for _, step := range strategy.Steps {
		if !step.Predicate(b) {
			continue
		}

		step.Execute(b)

		if !step.Next(b) {
			return
		}
	}
}
