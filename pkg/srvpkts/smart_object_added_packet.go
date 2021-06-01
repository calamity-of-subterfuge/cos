package srvpkts

// SmartObjectAddedPacket is a top-level packet that can be sent to a client
// to inform them a new smart object came into view.
type SmartObjectAddedPacket struct {
	// Type is smart-object-added
	Type string `mapstructure:"type" json:"type"`

	// GameTime is the time at which this packet was sent
	GameTime float64 `mapstructure:"game_time" json:"game_time"`

	// Object is the object that just entered view.
	Object SmartObjectSync `mapstructure:"object" json:"object"`
}

func (p *SmartObjectAddedPacket) GetType() string {
	return "smart-object-added"
}

func (p *SmartObjectAddedPacket) PrepareForMarshal() {
	p.Type = p.GetType()
}

func init() {
	registerPacketParser("smart-object-added", func(parsed map[string]interface{}) (Packet, error) {
		return parseSinglePacketOfType(parsed, &SmartObjectAddedPacket{})
	})
}
