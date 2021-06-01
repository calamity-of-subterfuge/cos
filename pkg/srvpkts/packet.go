package srvpkts

// Packet describes any of the packets within this package, which is useful for
// grabbing the Type of the packet without a large switch statement.
type Packet interface {
	// GetType returns the canonical unique identifier for the packet
	GetType() string

	// Prepare ensures that the default values for this packet have been
	// set correctly.
	PrepareForMarshal()
}
