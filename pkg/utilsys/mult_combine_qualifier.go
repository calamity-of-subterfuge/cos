package utilsys

type multCombineQualifier struct {
	children []Scorer
}

func (s *multCombineQualifier) Attached(world, actor interface{}) {
	for _, child := range s.children {
		child.Attached(world, actor)
	}
}

func (s *multCombineQualifier) Score() float64 {
	res := 1.0
	for _, child := range s.children {
		res *= child.Score()
	}
	return res
}

// MultCombineQualifier produces a score from the children by multiplying
// their scores together.
type MultCombineQualifier struct {
	// Children are the children the score is built from
	Children []ScorerBuilder
}

func (b MultCombineQualifier) Build() Scorer {
	children := make([]Scorer, len(b.Children))
	for idx, childBuilder := range b.Children {
		children[idx] = childBuilder.Build()
	}
	return &multCombineQualifier{children: children}
}
