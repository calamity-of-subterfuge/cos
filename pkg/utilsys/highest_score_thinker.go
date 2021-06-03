package utilsys

import (
	"math/rand"
)

type highestScoreThinker struct{}

func (t highestScoreThinker) Select(children []ScoredAction) int {
	currentHighScore := children[0].Scorer.Score()
	currentBestIndexes := []int{0}
	for idx := 1; idx < len(children); idx++ {
		score := children[idx].Scorer.Score()

		if score > currentHighScore {
			currentHighScore = score
			currentBestIndexes = currentBestIndexes[:1]
			currentBestIndexes[0] = idx
		} else if score == currentHighScore {
			currentBestIndexes = append(currentBestIndexes, idx)
		}
	}

	if len(currentBestIndexes) == 1 {
		return currentBestIndexes[0]
	} else {
		return currentBestIndexes[rand.Intn(len(currentBestIndexes))]
	}
}

// NewHighestScoreThinker produces a Thinker which performs whichever action
// from the given list of scored actions has the highest score. In the event
// of ties, it chooses uniformly at random from the ties.
func NewHighestScoreThinker(actions []ScoredActionBuilder) ActionBuilder {
	return NewThinkerBuilder(highestScoreThinker{}, actions)
}
