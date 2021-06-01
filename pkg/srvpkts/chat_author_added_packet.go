package srvpkts

// ChatAuthorAddedPacket informs a client that there is a new person
// able to hear their messages
type ChatAuthorAddedPacket struct {
	// chat-author-added
	Type string `mapstructure:"type" json:"type"`

	// GameTime this packet was sent
	GameTime float64 `mapstructure:"game_time" json:"game_time"`

	ChatAuthorSync
}

func (p *ChatAuthorAddedPacket) GetType() string {
	return "chat-author-added"
}

func (p *ChatAuthorAddedPacket) PrepareForMarshal() {
	p.Type = p.GetType()
}

func init() {
	registerPacketParser("chat-author-added", func(parsed map[string]interface{}) (Packet, error) {
		return parseSinglePacketOfType(parsed, &ChatAuthorAddedPacket{})
	})
}
