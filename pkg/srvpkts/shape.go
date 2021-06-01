package srvpkts

// Shape describes a single collision shape on a game object.
type Shape struct {
	// ShapeType is currently always "polygon"
	ShapeType string `mapstructure:"shape_type" json:"shape_type"`

	// Mass is the mass of this shape in kg
	Mass float64 `mapstructure:"mass" json:"mass"`

	// Details is currently always a PolygonDetails
	Details interface{} `mapstructure:"details" json:"details"`
}

// PolygonDetails are the Details for the "polygon" Shape
type PolygonDetails struct {
	// Vertices is the vertices of the polygon relative to its center
	// of gravity
	Vertices []Vector `mapstructure:"vertices" json:"vertices"`

	// Radius is the rounding on the edges in game units. The actual shape
	// is this amount larger than the vertices would imply in any
	// direction. This rounding reduces collision oddities.
	Radius float64 `mapstructure:"radius" json:"radius"`
}
