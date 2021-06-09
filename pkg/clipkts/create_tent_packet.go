package clipkts

import "github.com/jakecoffman/cp"

// CreateTentPacket is sent by an economy AI in order to place the tent
// for their team.
type CreateTentPacket struct {
	// Type is always create-tent
	Type string `json:"type" mapstructure:"type"`

	// Location is where you are trying to place the tent
	Location cp.Vector
}

func (p *CreateTentPacket) GetType() string {
	return "create-tent"
}

func (p *CreateTentPacket) PrepareForMarshal() {
	p.Type = p.GetType()
}

func init() {
	registerPacketParser("create-tent", func(parsed map[string]interface{}) (Packet, error) {
		return parseSinglePacketOfType(parsed, &CreateTentPacket{})
	})
}
