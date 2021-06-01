package srvpkts

// GameObjectUpdatePacket informs a client about a small update to a game object
type GameObjectUpdatePacket struct {
	// game-object-update
	Type string `mapstructure:"type" json:"type"`

	// GameTime this packet was sent
	GameTime float64 `mapstructure:"game_time" json:"game_time"`

	// UID of the object that changed
	UID string `mapstructure:"uid" json:"uid"`

	// Position is the new position of the object in game units
	Position Vector `mapstructure:"position" json:"position"`

	// Velocity is the new change in position of the object in game units per
	// second
	Velocity Vector `mapstructure:"velocity" json:"velocity"`

	// Rotation is the new rotation of the object in radians
	Rotation float64 `mapstructure:"rotation" json:"rotation"`

	// AngularVelocity is the new change in rotation of the object in radians
	// per second.
	AngularVelocity float64 `mapstructure:"angular_velocity" json:"angular_velocity"`

	// Animation is the new animation of the object
	Animation string `mapstructure:"animation" json:"animation"`

	// AnimationPlaying is true if the new animation should play, false if it
	// should be frozen on an arbitrary frame.
	AnimationPlaying bool `mapstructure:"animation_playing" json:"animation_playing"`

	// AnimationLooping is true if the new animation should start over when it
	// finishes and false if it should not.
	AnimationLooping bool `mapstructure:"animation_looping" json:"animation_looping"`
}

func (p *GameObjectUpdatePacket) GetType() string {
	return "game-object-update"
}

func (p *GameObjectUpdatePacket) PrepareForMarshal() {
	p.Type = p.GetType()
}

func init() {
	registerPacketParser("game-object-update", func(parsed map[string]interface{}) (Packet, error) {
		return parseSinglePacketOfType(parsed, &GameObjectUpdatePacket{})
	})
}
