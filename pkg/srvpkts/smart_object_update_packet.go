package srvpkts

// SmartObjectUpdatePacket informs the client about an update to a smart object
type SmartObjectUpdatePacket struct {
	GameObjectUpdatePacket

	// CurrentHealth is how much health this smart object has. This may be
	// a rounded representation of the objects true health.
	CurrentHealth int `mapstructure:"current_health" json:"current_health"`

	// Additional depends on the UnitType of this smart object.
	Additional interface{} `mapstructure:"additional" json:"additional"`
}

func (p *SmartObjectUpdatePacket) GetType() string {
	return "smart-object-update"
}

func (p *SmartObjectUpdatePacket) PrepareForMarshal() {
	p.Type = p.GetType()
}

func init() {
	registerPacketParser("smart-object-update", func(parsed map[string]interface{}) (Packet, error) {
		return parseSinglePacketOfType(parsed, &SmartObjectUpdatePacket{})
	})
}
