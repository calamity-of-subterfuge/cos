package clipkts

// SendCommandPacket is the client telling the server that they
// wish to execute a text command. We allow clients to decide
// what counts as a "command" so that we can support many different
// formats (e.g, slash, hyphen, whatever). The first character of
// the command is ignored for the purpose of parsing.
type SendCommandPacket struct {
	// Type is send-command
	Type string `json:"type"`

	// Text is the text of the command to execute; the first character
	// is ignored (so "/foo the bar" is parsed as "foo the bar"). Anything
	// after the 4096'th character is ignored.
	Text string `json:"text"`
}

func (p *SendCommandPacket) GetType() string {
	return "send-command"
}

func (p *SendCommandPacket) PrepareForMarshal() {
	p.Type = p.GetType()
}

func init() {
	registerPacketParser("send-command", func(parsed map[string]interface{}) (Packet, error) {
		var scPacket SendCommandPacket
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
