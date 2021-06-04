package utilsys

import (
	"math"
	"math/rand"
)

type softMaxProbabilisticThinker struct {
	threshold   float64
	scoreFactor float64
}

func (t softMaxProbabilisticThinker) Select(children []ScoredAction) int {
	// We should be in the perfect case for numerical stability with this
	// function, so long as the scores are really within 0.1 and 1, so long as
	// we're careful about the rolling sum for actually selecting a result, and
	// we use the identity softmax(x) = softmax(x + c) where x is a vector and c
	// is a vector of all a single value, and we set c to be near the maximum
	// value in x. x is the rescaled score vector and hence the max value should
	// be less than the factor.

	stabilityShift := t.scoreFactor - 1
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

		exponentiated := math.Exp(score*t.scoreFactor - stabilityShift)

		indexesInConsideration = append(indexesInConsideration, idx)

		compensatedExponentiated := exponentiated - compensation
		newSum := sum + compensatedExponentiated
		compensation = (newSum - sum) - compensatedExponentiated
		sum = newSum

		rollingSum = append(rollingSum, sum)
	}

	if sum == 0 { // everything underflowed :/
		return indexesInConsideration[rand.Intn(len(indexesInConsideration))]
	}

	seed := rand.Float64() * sum

	for idx, psum := range rollingSum {
		if seed < psum {
			return indexesInConsideration[idx]
		}
	}

	panic("shouldn't get here")
}

// NewSoftMaxProbabilisticThinker produces a Thinker which selects a child
// randomly from the children with a probability proportional to
// e^(score*factor). If there are any children whose score meets or exceeds the
// given threshold, then all children whose score is below the threshold are
// ignored.
//
// A factor of 1 makes this a pure soft-max function. A factor of 0 makes this
// a completely random choice. A higher factor reduces the amount of randomness.
// The factor is typically between 5 and 30
func NewSoftMaxProbabilisticThinker(threshold float64, factor float64, children []ScoredActionBuilder) ActionBuilder {
	if factor < 0 {
		panic("factor cannot be negative")
	}
	return NewThinkerBuilder(
		softMaxProbabilisticThinker{scoreFactor: factor, threshold: threshold},
		children,
	)
}
