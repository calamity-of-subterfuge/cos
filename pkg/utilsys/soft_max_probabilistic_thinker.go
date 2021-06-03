package utilsys

import (
	"math"
	"math/rand"
)

type softMaxProbabilisticThinker struct {
	threshold float64
}

func (t softMaxProbabilisticThinker) Select(children []ScoredAction) int {
	// We should be in the perfect case for numerical stability
	// with this function, so long as the scores are really within
	// 0.1 and 1, so long as we're careful about the rolling sum
	// for actually selecting a result.

	sawThreshold := false
	sum := 0.0
	compensation := 0.0 // Kahan summation algorithm (greatly improves numerical precision)
	indexesInConsideration := make([]int, 0, len(children))
	rollingSum := make([]float64, 0, len(children))

	for idx, child := range children {
		score := child.Scorer.Score()

		if sawThreshold && score < t.threshold {
			continue
		}

		if !sawThreshold && score >= t.threshold {
			sawThreshold = true
			sum = 0.0
			compensation = 0.0
			indexesInConsideration = indexesInConsideration[:0]
			rollingSum = rollingSum[:0]
		}

		exponentiated := math.Exp(score)

		indexesInConsideration = append(indexesInConsideration, idx)

		compensatedExponentiated := exponentiated - compensation
		newSum := sum + compensatedExponentiated
		compensation = (newSum - sum) - compensatedExponentiated
		sum = newSum

		rollingSum = append(rollingSum, sum)
	}

	seed := rand.Float64() * sum

	for idx, psum := range rollingSum {
		if seed < psum {
			return idx
		}
	}

	panic("shouldn't get here")
}

// NewSoftMaxProbabilisticThinker produces a Thinker which selects a child
// randomly from the children with a probability proportional to e^(score). If
// there are any children whose score meets or exceeds the given threshold, then
// all children whose score is below the threshold are ignored. Note that this
// thinker is sensitive to the scores of children - in particular, you should
// consider it a soft requirement that scores are within the normal 0-1 range
// when using this thinker, with 0.1-1 being the most common.
func NewSoftMaxProbabilisticThinker(threshold float64, children []ScoredActionBuilder) ActionBuilder {
	return NewThinkerBuilder(softMaxProbabilisticThinker{threshold: threshold}, children)
}
