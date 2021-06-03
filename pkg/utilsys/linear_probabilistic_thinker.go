package utilsys

import "math/rand"

type linearProbabilisticThinker struct {
	threshold float64
}

func (t linearProbabilisticThinker) Select(children []ScoredAction) int {
	sawThreshold := false
	sum := 0.0
	compensation := 0.0 // Kahan summation algorithm (greatly improves numerical precision)
	indexesInConsideration := make([]int, 0, len(children))
	rollingSumScoresInConsideration := make([]float64, 0, len(children))

	for idx, child := range children {
		score := child.Scorer.Score()

		if score < t.threshold && sawThreshold {
			continue
		}

		if !sawThreshold && score >= t.threshold {
			sawThreshold = true
			sum = 0
			compensation = 0
			indexesInConsideration = indexesInConsideration[:0]
			rollingSumScoresInConsideration = rollingSumScoresInConsideration[:0]
		}

		indexesInConsideration = append(indexesInConsideration, idx)

		compensatedScore := score - compensation
		newSum := sum + compensatedScore
		compensation = (newSum - sum) - compensatedScore
		sum = newSum

		rollingSumScoresInConsideration = append(rollingSumScoresInConsideration, sum)
	}

	seed := rand.Float64() * sum

	for idx, psum := range rollingSumScoresInConsideration {
		if seed < psum {
			return idx
		}
	}

	panic("shouldn't get here")
}

// NewLinearProbabilisticThinker produces a new Thinker which selects which
// child randomly, where the probability of a child being selected is
// proportional to its score. If any of the children scores meet or exceed the
// threshold, then all children below the threshold score are ignored.
func NewLinearProbabilisticThinker(threshold float64, children []ScoredActionBuilder) ActionBuilder {
	return NewThinkerBuilder(linearProbabilisticThinker{threshold: threshold}, children)
}
