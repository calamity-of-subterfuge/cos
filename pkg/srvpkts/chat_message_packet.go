package srvpkts

// ChatMessagePacket tells a player about a new message from someone
// nearby, or possibly about a server message
type ChatMessagePacket struct {
	// chat-message
	Type string `mapstructure:"type" json:"type"`

	// GameTime this packet was sent using our game time measurement,
	// which is consistent across all players but doesn't necessarily
	// correspond to anything real
	GameTime float64 `mapstructure:"game_time" json:"game_time"`

	// Time is the unix time that the message was canonically received
	Time float64 `mapstructure:"time" json:"time"`

	// AuthorUID is the uid of the chat author which sent this message
	AuthorUID string `mapstructure:"author_uid" json:"author_uid"`

	// Text that should be displayed. Must be treated as untrusted by
	// the client.
	Text string `mapstructure:"text" json:"text"`
}

func (p *ChatMessagePacket) GetType() string {
	return "chat-message"
}

func (p *ChatMessagePacket) PrepareForMarshal() {
	p.Type = p.GetType()
}

func init() {
	registerPacketParser("chat-message", func(parsed map[string]interface{}) (Packet, error) {
		return parseSinglePacketOfType(parsed, &ChatMessagePacket{})
	})
}
