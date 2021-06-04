package utilsys

import (
	"math/rand"
	"time"
)

type idleAction struct {
	minDuration time.Duration
	maxDuration time.Duration

	state             ActionState
	durationRemaining time.Duration
}

func (a *idleAction) State() ActionState {
	return a.state
}

func (a *idleAction) Attached(world, actor interface{}) {
	a.state = ActionStateRequested
}

func (a *idleAction) Execute(delta time.Duration) {
	if a.state == ActionStateRequested {
		a.durationRemaining = a.minDuration + time.Duration(rand.Int63n(int64(a.maxDuration-a.minDuration)))
		a.state = ActionStateExecuting
	}

	a.durationRemaining -= delta
	if a.durationRemaining < 0 {
		a.state = ActionStateSuccess
	}
}

func (a *idleAction) Cancel() {
	a.state = ActionStateFailure
}

func (a *idleAction) FinishCanceling(delta time.Duration) {
	a.state = ActionStateFailure
}

func (a *idleAction) Reset() {
	a.state = ActionStateRequested
}

// IdleAction idles for a random amount of time selected uniformly between the
// min and max duration.
type IdleAction struct {
	// MinDuration is the minimum duration to idle for.
	MinDuration time.Duration

	// MaxDuration is the maximum duration to idle for.
	MaxDuration time.Duration
}

func (b IdleAction) Build() Action {
	return &idleAction{
		minDuration: b.MinDuration,
		maxDuration: b.MaxDuration,
	}
}
