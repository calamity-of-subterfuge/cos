package utilsys

import "time"

// Thinker describes the standard Thinker interface which can be wrapped
// with NewThinkerBuilder to produce an ActionBuilder.
//
// Note that when we use the word "Thinker" in this package we are almost never
// referring to this interface. We are simply referring to any ActionBuilder
// which selects which ActionBuilder to delegate to based on its score. The most
// common way to implement that concept of Thinker is by implementing this
// interface and using NewThinkerBuilder as the constructor for the
// ActionBuilder.
type Thinker interface {
	// Select the index of the child which should be exected from the given
	// slice of scored actions.
	Select(children []ScoredAction) int
}

type thinker struct {
	thinker       Thinker
	scoredActions []ScoredAction
	currentIndex  int
}

func (t *thinker) State() ActionState {
	if t.currentIndex == -1 {
		return ActionStateInit
	}

	return t.scoredActions[t.currentIndex].Action.State()
}

func (t *thinker) Attached(world, actor interface{}) {
	for _, scoredAction := range t.scoredActions {
		scoredAction.Action.Attached(world, actor)
		scoredAction.Scorer.Attached(world, actor)
	}

	t.currentIndex = t.thinker.Select(t.scoredActions)
}

func (t *thinker) Execute(delta time.Duration) {
	t.scoredActions[t.currentIndex].Action.Execute(delta)
}

func (t *thinker) Cancel() {
	t.scoredActions[t.currentIndex].Action.Cancel()
}

func (t *thinker) FinishCanceling(delta time.Duration) {
	t.scoredActions[t.currentIndex].Action.FinishCanceling(delta)
}

func (t *thinker) Reset() {
	t.scoredActions[t.currentIndex].Action.Reset()
	t.currentIndex = t.thinker.Select(t.scoredActions)
}

type thinkerBuilder struct {
	thinker Thinker
	actions []ScoredActionBuilder
}

func (b thinkerBuilder) Build() Action {
	initializedActions := make([]ScoredAction, len(b.actions))
	for idx, actB := range b.actions {
		initializedActions[idx] = ScoredAction{
			Action: actB.Action.Build(),
			Scorer: actB.Scorer.Build(),
		}
	}
	return &thinker{
		thinker:       b.thinker,
		scoredActions: initializedActions,
		currentIndex:  -1,
	}
}

// NewThinkerBuilder produces an ActionBuilder out of something implementing
// the Thinker interface and the non-empty slice of children.
func NewThinkerBuilder(thinker Thinker, actions []ScoredActionBuilder) ActionBuilder {
	if len(actions) == 0 {
		panic("error: actions cannot be empty")
	}
	return thinkerBuilder{
		thinker: thinker,
		actions: actions,
	}
}
