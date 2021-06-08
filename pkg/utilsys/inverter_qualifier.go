package utilsys

type inverterQualifier struct {
	scorer Scorer
}

func (s *inverterQualifier) Attached(world, actor interface{}) {
	s.scorer.Attached(world, actor)
}

func (s *inverterQualifier) Score() float64 {
	return 1 - s.scorer.Score()
}

// InverterQualifier inverts the score of the child, i.e., returns
// 1 - Scorer.Score()
type InverterQualifier struct {
	// Scorer is the scorer to invert
	Scorer ScorerBuilder
}

func (b InverterQualifier) Build() Scorer {
	return &inverterQualifier{
		scorer: b.Scorer.Build(),
	}
}
