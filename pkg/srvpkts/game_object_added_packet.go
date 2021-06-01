package srvpkts

// GameObjectAddedPacket is a top-level packet that can be sent to a client
// to inform them a new game object came into view.
type GameObjectAddedPacket struct {
	// game-object-added
	Type string `mapstructure:"type" json:"type"`

	// GameTime this packet was sent
	GameTime float64 `mapstructure:"game_time" json:"game_time"`

	// Object is the object that just entered view.
	Object GameObjectSync `mapstructure:"object" json:"object"`
}

func (p *GameObjectAddedPacket) GetType() string {
	return "game-object-added"
}

func (p *GameObjectAddedPacket) PrepareForMarshal() {
	p.Type = p.GetType()
}

func init() {
	registerPacketParser("game-object-added", func(parsed map[string]interface{}) (Packet, error) {
		return parseSinglePacketOfType(parsed, &GameObjectAddedPacket{})
	})
}
