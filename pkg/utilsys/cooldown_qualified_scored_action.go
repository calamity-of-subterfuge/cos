package utilsys

import "time"

type cooldownQualifiedScorer struct {
	scorer         Scorer
	minCooldown    time.Duration
	maxCooldown    time.Duration
	lastFinishedAt *time.Time
}

func (s *cooldownQualifiedScorer) Attached(world, actor interface{}) {
	s.scorer.Attached(world, actor)
}

func (s *cooldownQualifiedScorer) Score() float64 {
	timeSinceLast := time.Since(*s.lastFinishedAt)

	if timeSinceLast < s.minCooldown {
		return 0
	} else if timeSinceLast < s.maxCooldown {
		progress := float64(timeSinceLast-s.minCooldown) / float64(s.maxCooldown)
		return progress * s.scorer.Score()
	} else {
		return s.scorer.Score()
	}
}

type cooldownQualifiedAction struct {
	action                      Action
	lastFinishedAt              *time.Time
	cooldownSuppressedOnSuccess bool
	cooldownSuppressedOnFailure bool
}

func (a *cooldownQualifiedAction) State() ActionState {
	return a.action.State()
}

func (a *cooldownQualifiedAction) Attached(world, actor interface{}) {
	a.action.Attached(world, actor)
	a.checkIfFinished()
}

func (a *cooldownQualifiedAction) Execute(delta time.Duration) {
	a.action.Execute(delta)
	a.checkIfFinished()
}

func (a *cooldownQualifiedAction) Cancel() {
	a.action.Cancel()
	a.checkIfFinished()
}

func (a *cooldownQualifiedAction) FinishCanceling(delta time.Duration) {
	a.action.FinishCanceling(delta)
	a.checkIfFinished()
}

func (a *cooldownQualifiedAction) Reset() {
	a.action.Reset()
	a.checkIfFinished()
}

func (a *cooldownQualifiedAction) checkIfFinished() {
	switch a.State() {
	case ActionStateFailure:
		if !a.cooldownSuppressedOnFailure {
			*a.lastFinishedAt = time.Now()
		}
	case ActionStateSuccess:
		if !a.cooldownSuppressedOnSuccess {
			*a.lastFinishedAt = time.Now()
		}
	}
}

// CooldownQualifiedScoredAction creates a dependency between the scorer
// and action provided. Specifically, the Score is set to 0 if it's been
// less than the MinCooldown since the Action completed, it is scaled
// linearly between 0 and 1 between the MinCooldown and MaxCooldown, and
// it is unmodified past the MaxCooldown.
type CooldownQualifiedScoredAction struct {
	ScoredAction ScoredActionBuilder

	MinCooldown time.Duration
	MaxCooldown time.Duration

	CooldownSuppressedOnSuccess bool
	CooldownSuppressedOnFailure bool
}

func (b CooldownQualifiedScoredAction) Build() ScoredAction {
	var finishedAt time.Time

	built := b.ScoredAction.Build()
	return ScoredAction{
		Scorer: &cooldownQualifiedScorer{
			scorer:         built.Scorer,
			minCooldown:    b.MinCooldown,
			maxCooldown:    b.MaxCooldown,
			lastFinishedAt: &finishedAt,
		},
		Action: &cooldownQualifiedAction{
			action:                      built.Action,
			cooldownSuppressedOnSuccess: b.CooldownSuppressedOnSuccess,
			cooldownSuppressedOnFailure: b.CooldownSuppressedOnFailure,
			lastFinishedAt:              &finishedAt,
		},
	}
}
