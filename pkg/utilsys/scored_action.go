package utilsys

// ScorerBuilderAndActionBuilder merges a scorer and an action. The
// scorer and action might not be independent.
type ScorerBuilderAndActionBuilder struct {
	// Action builds actions
	Action ActionBuilder

	// Scorer builds scorers
	Scorer ScorerBuilder
}

// Build implements ScoredActionBuilder but only works if the action and scorer
// are independent. Hence when receiving a ScorerBuilderAndActionBuilder as that
// type this function should not be called, but this allows receiving a
// ScorerBuilderAndActionBuilder type asserted as a ScoredActionBuilder, which
// is convenenient if you want to make a ScoredActionBuilder with an independent
// Scorer and Action.
func (b ScorerBuilderAndActionBuilder) Build() ScoredAction {
	return ScoredAction{
		Action: b.Action.Build(),
		Scorer: b.Scorer.Build(),
	}
}

// ScoredActionBuilder builds pairs of actions and scores. If the action and
// scorer are independent then the obvious implementation is done via
// ScorerBuilderAndActionBuilder. For example, when the Scorer is a fixed
// scorer, then the action and scorer are independent, i.e., you can build the
// action without building the scorer. If the Scorer is a cooldown scorer, where
// the cooldown starts when the action finishes, then the Scorer is not
// independent of the Action and hence they must be built in tandem.
type ScoredActionBuilder interface {
	Build() ScoredAction
}

// ScoredAction is a convenience struct for describing an Action and a Score.
// This is typically only used within actual Actions, - for nesting
// ActionBuilder's use ScoredActionBuilder
type ScoredAction struct {
	Action Action
	Scorer Scorer
}
