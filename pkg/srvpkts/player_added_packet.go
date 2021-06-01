package srvpkts

// PlayerAddedPacket is a top-level packet that can be sent to a client
// to inform them a new player came into view.
type PlayerAddedPacket struct {
	// Type is player-added
	Type string `mapstructure:"type" json:"type"`

	// GameTime is the time at which this packet was sent
	GameTime float64 `mapstructure:"game_time" json:"game_time"`

	// Object is the player that just entered view.
	Object PlayerSync `mapstructure:"object" json:"object"`
}

func (p *PlayerAddedPacket) GetType() string {
	return "player-added"
}

func (p *PlayerAddedPacket) PrepareForMarshal() {
	p.Type = p.GetType()
}

func init() {
	registerPacketParser("player-added", func(parsed map[string]interface{}) (Packet, error) {
		return parseSinglePacketOfType(parsed, &PlayerAddedPacket{})
	})
}
