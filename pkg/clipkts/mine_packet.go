package clipkts

// MinePacket is used for AI to mine resources which are near a unit they
// control.
type MinePacket struct {
	// Type is always mine
	Type string `json:"type" mapstructure:"type"`

	// MiningUID is the UID of the object doing the mining
	MiningUID string `json:"mining_uid" mapstructure:"mining_uid"`

	// MinedUID is the UID of the object being mined
	MinedUID string `json:"mined_uid" mapstructure:"mined_uid"`
}

func (p *MinePacket) GetType() string {
	return "mine"
}

func (p *MinePacket) PrepareForMarshal() {
	p.Type = p.GetType()
}

func init() {
	registerPacketParser("mine", func(parsed map[string]interface{}) (Packet, error) {
		return parseSinglePacketOfType(parsed, &MinePacket{})
	})
}
