package srvpkts

// TeamResourceChangedPacket informs a player that he amount of resources their
// team has has changed
type TeamResourceChangedPacket struct {
	// Type is team-resource-changed
	Type string `mapstructure:"type" json:"type"`

	// GameTime is the time at which this packet was sent
	GameTime float64 `mapstructure:"game_time" json:"game_time"`

	// Resources contains a subset (not necessarily strict) of the teams
	// resources; for each key in this map, the value is the amount
	// of that resource the team has
	Resources map[string]int `mapstructure:"resources" json:"resources"`
}

func (p *TeamResourceChangedPacket) GetType() string {
	return "team-resource-changed"
}

func (p *TeamResourceChangedPacket) PrepareForMarshal() {
	p.Type = p.GetType()
}

func init() {
	registerPacketParser("team-resource-changed", func(parsed map[string]interface{}) (Packet, error) {
		return parseSinglePacketOfType(parsed, &TeamResourceChangedPacket{})
	})
}
