package srvpkts

// ChatAuthorUpdatePacket informs a client that one of the people they
// can hear in chat changed something about how their messages are
// displayed
type ChatAuthorUpdatePacket struct {
	// chat-author-update
	Type string `mapstructure:"type" json:"type"`

	// GameTime this packet was sent
	GameTime float64 `mapstructure:"game_time" json:"game_time"`

	ChatAuthorSync
}

func (p *ChatAuthorUpdatePacket) GetType() string {
	return "chat-author-update"
}

func (p *ChatAuthorUpdatePacket) PrepareForMarshal() {
	p.Type = p.GetType()
}

func init() {
	registerPacketParser("chat-author-update", func(parsed map[string]interface{}) (Packet, error) {
		return parseSinglePacketOfType(parsed, &ChatAuthorUpdatePacket{})
	})
}
