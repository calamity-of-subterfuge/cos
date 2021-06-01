package clipkts

// SendLocalMessagePacket is the client telling the server that they
// wish to send some text to nearby players
type SendLocalMessagePacket struct {
	// Type is send-local-message
	Type string `json:"type"`

	// Text is the text of the message; will be treated as untrusted.
	// Anything after the 4096th character is ignored.
	Text string `json:"text"`
}

func (p *SendLocalMessagePacket) GetType() string {
	return "send-local-message"
}

func (p *SendLocalMessagePacket) PrepareForMarshal() {
	p.Type = p.GetType()
}

func init() {
	registerPacketParser("send-local-message", func(parsed map[string]interface{}) (Packet, error) {
		var scPacket SendLocalMessagePacket
		_, err := parseSinglePacketOfType(parsed, &scPacket)
		if err != nil {
			return nil, err
		}
		if len(scPacket.Text) > 4096 {
			scPacket.Text = scPacket.Text[:4096]
		}
		return &scPacket, nil
	})
}
