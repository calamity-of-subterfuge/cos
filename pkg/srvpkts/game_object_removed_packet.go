package srvpkts

// GameObjectRemovedPacket is used to inform the client they are no longer
// able to see the given game object, likely because it's just out view now
type GameObjectRemovedPacket struct {
	// game-object-removed
	Type string `mapstructure:"type" json:"type"`

	// GameTime this packet was sent
	GameTime float64 `mapstructure:"game_time" json:"game_time"`

	// UID of the object removed
	UID string `mapstructure:"uid" json:"uid"`
}

func (p *GameObjectRemovedPacket) GetType() string {
	return "game-object-removed"
}

func (p *GameObjectRemovedPacket) PrepareForMarshal() {
	p.Type = p.GetType()
}

func init() {
	registerPacketParser("game-object-removed", func(parsed map[string]interface{}) (Packet, error) {
		return parseSinglePacketOfType(parsed, &GameObjectRemovedPacket{})
	})
}
