package utilsys

// ScoredActionBuilder is just a convenience struct describing an ActionBuilder
// with an associated ScoreBuilder
type ScoredActionBuilder struct {
	// Action builds actions
	Action ActionBuilder

	// Scorer builds scorers
	Scorer ScorerBuilder
}

// ScoredAction is a convenience struct for describing an Action and a Score.
// This is typically only used within actual Actions, - for nesting
// ActionBuilder's use ScoredActionBuilder
type ScoredAction struct {
	Action Action
	Scorer Scorer
}
