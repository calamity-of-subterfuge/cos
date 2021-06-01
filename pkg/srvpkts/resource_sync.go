package srvpkts

// ResourceSync describes everything a client needs to know to render a resource
// they've never seen before
type ResourceSync struct {
	// UID is the unique identifier for this resource
	UID string `map_structure:"uid" json:"uid"`

	// SheetURL is the URL of the spritesheet (JSON) where the
	// icon for this resource can be found
	SheetURL string `map_structure:"sheet_url" json:"sheet_url"`

	// Animation is the name within the sheet for the icon for this
	// resource. Currently this is always a single-icon animation but we
	// leave room for real animations later
	Animation string `map_structure:"animation" json:"animation"`

	// Name is the display name for this resource
	Name string `map_structure:"name" json:"name"`
}
