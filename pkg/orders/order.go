package orders

// Order describes any of the orders within this package, which is useful for
// grabbing the Type of the package without a large switch statement.
type Order interface {
	// GetType returns the canonical unique identifier for the order
	GetType() string

	// PrepareForMarshal ensures that the default values for this order have
	// been set correctly.
	PrepareForMarshal()
}
