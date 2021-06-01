package client

import (
	"github.com/calamity-of-subterfuge/cos/v2/pkg/srvpkts"
)

// SmartObjectAdditional describes additional information on a SmartObject
type SmartObjectAdditional interface {
	// Update the additional information using the given packet
	Update(update *srvpkts.SmartObjectUpdatePacket) error
}

// BlankSmartObjectAdditional is for objects with no additional information
type BlankSmartObjectAdditional struct{}

// Update is no-op
func (a *BlankSmartObjectAdditional) Update(update *srvpkts.SmartObjectUpdatePacket) error {
	return nil
}

// SmartObjectAdditionalParser parses a SmartObjectAdditional from
// a SmartObjectSync
type SmartObjectAdditionalParser func(*srvpkts.SmartObjectSync) (SmartObjectAdditional, error)

var smartObjectsByUnitType map[string]SmartObjectAdditionalParser = make(map[string]SmartObjectAdditionalParser)

// ParseSmartObjectAdditional parses a SmartObjectAdditional from the given
// SmartObjectSync based on its UnitType
func ParseSmartObjectAdditional(sync *srvpkts.SmartObjectSync) (SmartObjectAdditional, error) {
	parser, found := smartObjectsByUnitType[sync.UnitType]
	if !found {
		return &BlankSmartObjectAdditional{}, nil
	}

	return parser(sync)
}
