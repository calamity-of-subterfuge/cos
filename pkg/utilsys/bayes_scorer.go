package utilsys

import "math/big"

// BayesFactor describes something which can alter our prediction
// about something by a given factor. It is stateful and only
// used for a single actor in a single world.
type BayesFactor interface {
	// Attached is called once to tell the factor the world and actor
	// it is operating within. Typically the factor will store more
	// strictly typed representations.
	Attached(world, actor interface{})

	// Factor returns how much information was gained by this factor.
	// If this factor is not informative, this will return
	// `big.NewRat(1, 1)`. If this increases the probability of success
	// by a factor of 2, this returns `big.NewRat(2, 1)`. If it decreases
	// the odds of success by a factor of 2, this returns `big.NewRat(1, 2)`.
	Factor() *big.Rat
}

// BayesFactorBuilder acts as a constructor for BayesFactors, since we need
// one BayesFactor per actor. Typically stateless.
type BayesFactorBuilder interface {
	// Build a new BayesFactor not attached yet.
	Build() BayesFactor
}

type bayesScorer struct {
	prior   *big.Rat
	factors []BayesFactor
}

func (s *bayesScorer) Attached(world, actor interface{}) {
	for _, factor := range s.factors {
		factor.Attached(world, actor)
	}
}

func (s *bayesScorer) Score() float64 {
	res := big.NewRat(1, 1).Set(s.prior)
	for _, factor := range s.factors {
		res.Mul(res, factor.Factor())
	}

	num := res.Num()
	denom := res.Denom()

	resultScore := big.NewFloat(0)
	resultScore.SetInt(big.NewInt(0).Add(num, denom))
	resultScore.Quo(big.NewFloat(0).SetInt(num), resultScore)
	f, _ := resultScore.Float64()
	return f
}

// BayesScorer is a type of ScorerBuilder that assumes that the utility
// of the action is 1, but it only succeeds probabilistically. It has
// some general chance at success, such as 1 success per 4 failures. It
// also has a set of things which alter its odds of success based on the
// world, such as "succeeds twice as often when Y is researched" in order
// to produce the final probability of succeeds.
type BayesScorer struct {
	prior   *big.Rat
	factors []BayesFactorBuilder
}

func (s BayesScorer) Build() Scorer {
	builtFactors := make([]BayesFactor, len(s.factors))
	for idx, builder := range s.factors {
		builtFactors[idx] = builder.Build()
	}
	return &bayesScorer{
		prior:   s.prior,
		factors: builtFactors,
	}
}

// NewBayesScorer produces a ScorerBuilder that has a score of 1 on the action,
// but the action only succeeds probabilistically. The estimate of the odds of
// success is prior. Note that prior should NOT be interpreted as a fraction,
// e.g., 3/5. Instead, it's interpreted such that the numerator is the number of
// successes and the denominator is the number of failures. so "3/9" should be
// interpreted as 3 successes to 9 failures. A "1/1" prior means 1 success to
// 1 failure, aka a 50% chance of success.
func NewBayesScorer(prior *big.Rat, factors []BayesFactorBuilder) ScorerBuilder {
	if prior == nil {
		panic("prior cannot be nil")
	}
	if len(factors) == 0 {
		panic("factors cannot be empty")
	}

	return BayesScorer{
		prior:   prior,
		factors: factors,
	}
}
