package client

import "github.com/calamity-of-subterfuge/cos/v2/pkg/srvpkts"

// Resource describes a resource from the perspective of the client.
type Resource struct {
	// UID of this resource
	UID string

	// SheetURL where the PixiJS spritesheet containing this resource
	// can be found
	SheetURL string

	// Animation within the sheet for this resource
	Animation string

	// How much of this resource we have
	Amount int
}

// Sync this resource with the given information
func (r *Resource) Sync(sync *srvpkts.ResourceSync) *Resource {
	r.UID = sync.UID
	r.SheetURL = sync.SheetURL
	r.Animation = sync.Animation
	return r
}
