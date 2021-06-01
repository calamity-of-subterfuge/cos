package srvpkts

// Vector describes a basic float-based 2d vector
type Vector struct {
	// X coordinate
	X float64 `json:"x" mapstructure:"x"`

	// Y coordinate
	Y float64 `json:"y" mapstructure:"y"`
}
