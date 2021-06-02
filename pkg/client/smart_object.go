package client

import (
	"log"

	"github.com/calamity-of-subterfuge/cos/pkg/srvpkts"
)

// SmartObject describes a "fancy" non-player unit in the world. These are
// often controllable by someone and have special orders. The UnitType of
// the smart object acts as the enum for the special behavior of the
// object.
type SmartObject struct {
	// GameObject for this smart object
	GameObject *GameObject

	// UnitType that this smart object is acting as
	UnitType string

	// CurrentHealth of this unit, may be a rounded representation of the
	// units true health
	CurrentHealth int

	// MaxHealth of this unit
	MaxHealth int

	// ControllingTeam is the team which controls this object.
	ControllingTeam int

	// ControllingRole is the role required to control this object. See core.Role
	// and core.RoleToName.
	ControllingRole string

	// Additional contains the SyncInfo on the Unit controlling this smart object,
	// and depends on the type of unit.
	Additional SmartObjectAdditional
}

// Sync this smart object with the given sync information
func (o *SmartObject) Sync(sync *srvpkts.SmartObjectSync) *SmartObject {
	o.GameObject = (&GameObject{}).Sync(&sync.GameObjectSync)
	o.UnitType = sync.UnitType
	o.CurrentHealth = sync.CurrentHealth
	o.MaxHealth = sync.MaxHealth
	o.ControllingTeam = sync.ControllingTeam
	o.ControllingRole = sync.ControllingRole
	add, err := ParseSmartObjectAdditional(sync)
	if err != nil {
		log.Fatalf("while syncing smart object with %v: %v", sync, err)
	}
	o.Additional = add
	return o
}

// Update this smart object with the given update packet
func (o *SmartObject) Update(update *srvpkts.SmartObjectUpdatePacket) *SmartObject {
	o.GameObject.Update(&update.GameObjectUpdatePacket)
	o.CurrentHealth = update.CurrentHealth
	o.Additional.Update(update)
	return o
}
