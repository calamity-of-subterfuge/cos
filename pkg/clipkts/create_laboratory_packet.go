package clipkts

import "github.com/jakecoffman/cp"

// CreateLaboratoryPacket is the packet sent by the science AIs to the
// server when the science AI would like to create a laboratory nearby.
type CreateLaboratoryPacket struct {
	// Type is always create-laboratory
	Type string `json:"type" mapstructure:"type"`

	// Location that the laboratory should be placed at
	Location cp.Vector
}

func (p *CreateLaboratoryPacket) GetType() string {
	return "create-laboratory"
}

func (p *CreateLaboratoryPacket) PrepareForMarshal() {
	p.Type = p.GetType()
}

func init() {
	registerPacketParser("create-laboratory", func(parsed map[string]interface{}) (Packet, error) {
		return parseSinglePacketOfType(parsed, &CreateLaboratoryPacket{})
	})
}
