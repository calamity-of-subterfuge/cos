package srvpkts

// ChatAuthorRemovedPacket informs a client that someone left range
// for communicating with them
type ChatAuthorRemovedPacket struct {
	// chat-author-removed
	Type string `mapstructure:"type" json:"type"`

	// GameTime this packet was sent
	GameTime float64 `mapstructure:"game_time" json:"game_time"`

	// UID is the identifier for the chat author
	UID string `mapstructure:"uid" json:"uid"`
}

func (p *ChatAuthorRemovedPacket) GetType() string {
	return "chat-author-removed"
}

func (p *ChatAuthorRemovedPacket) PrepareForMarshal() {
	p.Type = p.GetType()
}

func init() {
	registerPacketParser("chat-author-removed", func(parsed map[string]interface{}) (Packet, error) {
		return parseSinglePacketOfType(parsed, &ChatAuthorRemovedPacket{})
	})
}
