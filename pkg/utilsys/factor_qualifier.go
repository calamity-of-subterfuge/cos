package utilsys

type factorQualifier struct {
	factor float64
	scorer Scorer
}

func (q *factorQualifier) Attached(world, actor interface{}) {
	q.scorer.Attached(world, actor)
}

func (q *factorQualifier) Score() float64 {
	return q.factor * q.scorer.Score()
}

type factorQualifierBuilder struct {
	factor float64
	scorer ScorerBuilder
}

func (b factorQualifierBuilder) Build() Scorer {
	return &factorQualifier{
		factor: b.factor,
		scorer: b.scorer.Build(),
	}
}

// NewFactorQualifier creates a new Qualifier that qualifies the score
// of the child by multiplying it by the given factor.
func NewFactorQualifier(factor float64, scorer ScorerBuilder) ScorerBuilder {
	return factorQualifierBuilder{factor: factor, scorer: scorer}
}
