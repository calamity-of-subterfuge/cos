package clipkts

import "fmt"

// IssueSmartObjectOrderPacket describes the client trying to issue an order
// to a smart object by the given UID.
type IssueSmartObjectOrderPacket struct {
	// Type is always issue-smart-object-order
	Type string `json:"type" mapstructure:"type"`

	// UID is the uid of the smart object the order is being issued to
	UID string `json:"uid" mapstructure:"uid"`

	// Order contains the order being issued and should be processed based
	// on the UnitType of the smart object.
	Order map[string]interface{} `json:"order" mapstructure:"order"`
}

func (p *IssueSmartObjectOrderPacket) GetType() string {
	return "issue-smart-object-order"
}

func (p *IssueSmartObjectOrderPacket) PrepareForMarshal() {
	p.Type = p.GetType()
}

func init() {
	registerPacketParser("issue-smart-object-order", func(parsed map[string]interface{}) (Packet, error) {
		var issSOOrder IssueSmartObjectOrderPacket
		_, err := parseSinglePacketOfType(parsed, &issSOOrder)
		if err != nil {
			return nil, err
		}
		if issSOOrder.UID == "" {
			return nil, fmt.Errorf("uid must not be blank")
		}
		if issSOOrder.Order == nil {
			return nil, fmt.Errorf("order cannot be empty")
		}
		return &issSOOrder, nil
	})
}
