package orders

import "errors"

// RespondTradeOrder is sent by client to respond to a particular trade request
// on a hut. This only works for trade requests that are from another team; to
// withdraw a trade order use WithdrawTradeOrder
type RespondTradeOrder struct {
	// Type is always respond-trade
	Type string `json:"type" mapstructure:"type"`

	// UID of the trade offer which is being accepted or rejected
	UID string `json:"uid" mapstructure:"uid"`

	// Accepted is true if we are trying to accept the offer and false if we are
	// trying to reject the offer.
	Accepted bool `json:"accepted" mapstructure:"accepted"`
}

func (o *RespondTradeOrder) GetType() string {
	return "respond-trade"
}

func (o *RespondTradeOrder) PrepareForMarshal() {
	o.Type = o.GetType()
}

func init() {
	registerOrderParser("respond-trade", func(m map[string]interface{}) (Order, error) {
		var res RespondTradeOrder
		_, err := parseSingleOrderOfType(m, &res)
		if err != nil {
			return nil, err
		}

		if res.UID == "" {
			return nil, errors.New("missing UID")
		}

		return &res, err
	})
}
