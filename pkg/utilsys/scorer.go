package utilsys

// A Scorer is something which is capable of producing a measure of utility for
// something. There is one instance per actor in the world and it is assumed to
// be stateful.
type Scorer interface {
	// Attached is called once when the Scorer is first attached to the AI to
	// let it know which world it is operating in and the actor which the scorer
	// is determining the utility of the action for. It should store these with
	// a stricter type.
	Attached(world, actor interface{})

	// Score returns the current measure, typically a 0-1 value where 0 is the
	// least valuable and 1 is the most valuable.
	Score() float64
}

// ScorerBuilder builds Scorer's
type ScorerBuilder interface {
	// Build a new scorer and return it so it may be attached
	Build() Scorer
}
