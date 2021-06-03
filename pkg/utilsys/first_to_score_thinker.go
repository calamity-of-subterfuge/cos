package utilsys

type firstToScoreThinker struct {
	threshold float64
	fallback  Thinker
}

func (t firstToScoreThinker) Select(children []ScoredAction) int {
	for idx, child := range children {
		if child.Scorer.Score() >= t.threshold {
			return idx
		}
	}

	return t.fallback.Select(children)
}

// NewFirstToScoreThinker produces a Thinker which performs the
// first child whose score meets or exceeds the threshold, falling
// back to a HighestScoreThinker if no children meet or exceed the
// threshold.
func NewFirstToScoreThinker(threshold float64, children []ScoredActionBuilder) ActionBuilder {
	return NewThinkerBuilder(firstToScoreThinker{
		threshold: threshold,
		fallback:  highestScoreThinker{},
	}, children)
}
