package utilsys

type fixedScorer struct {
	score float64
}

func (s *fixedScorer) Attached(world, actor interface{}) {}
func (s *fixedScorer) Score() float64                    { return s.score }

// FixedScorer is the simplest type of scorer which always returns
// the same value
type FixedScorer struct {
	// Score is the score that the fixed scorer returns for all actors
	// at all times in all worlds
	Score float64
}

// Build implements ScorerBuilder
func (s FixedScorer) Build() Scorer {
	return &fixedScorer{score: s.Score}
}
