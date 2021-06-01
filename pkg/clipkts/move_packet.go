package clipkts

import "fmt"

// MovePacket describes the client trying to move the game object with
// the given uid in the given direction.
type MovePacket struct {
	// Type is always move
	Type string `json:"type" mapstructure:"type"`

	// UID of the game object to move
	UID string `json:"uid" mapstructure:"uid"`

	// Direction to move the object in. If the magnitude is less than one then
	// move in the move force multiplied by this vector, otherwise move in this
	// direction at the maximum move force
	Direction Vector `json:"dir" mapstructure:"dir"`
}

func (p *MovePacket) GetType() string {
	return "move"
}

func (p *MovePacket) PrepareForMarshal() {
	p.Type = p.GetType()
}

func init() {
	registerPacketParser("move", func(parsed map[string]interface{}) (Packet, error) {
		var movePacket MovePacket
		_, err := parseSinglePacketOfType(parsed, &movePacket)
		if err != nil {
			return nil, err
		}
		if movePacket.UID == "" {
			return nil, fmt.Errorf("uid cannot be blank")
		}
		return &movePacket, nil
	})
}
