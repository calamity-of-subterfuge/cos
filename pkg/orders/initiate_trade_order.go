package orders

import "errors"

// InitiateTradeOrder is an order sent to a tent to initiate a
// resource trade with another team. If the other team accepts using
// an RespondTradeOrder on a tent then the trade goes through. If they
// reject the trade order with a RespondTradeOrder then the trade is
// canceled.
//
// The list of active trades on the tent are visible in the tents
// Additional details. See unitdets.Tent for the format.
type InitiateTradeOrder struct {
	// Type is always initiate-trade
	Type string `json:"type" mapstructure:"type"`

	// Team that you want to trade with, cannot be your own team.
	Team int `json:"team" mapstructure:"team"`

	// Offer of resources you will give to the other team. Keys are the uids
	// of the resources that you want to give and the values are the amounts
	// of those resources.
	Offer map[string]int `json:"offer" mapstructure:"offer"`

	// Request of resources you will receive from the other team. Keys are the
	// uids of the resources that you want to receive and the values are the
	// amounts of those resources.
	Request map[string]int `json:"request" mapstructure:"request"`
}

func (o *InitiateTradeOrder) GetType() string {
	return "initiate-trade"
}

func (o *InitiateTradeOrder) PrepareForMarshal() {
	o.Type = o.GetType()
}

func init() {
	registerOrderParser("initiate-trade", func(m map[string]interface{}) (Order, error) {
		var res InitiateTradeOrder
		_, err := parseSingleOrderOfType(m, &res)
		if err != nil {
			return nil, err
		}

		if res.Offer == nil {
			res.Offer = make(map[string]int)
		}

		if res.Request == nil {
			res.Request = make(map[string]int)
		}

		if len(res.Offer) == 0 && len(res.Request) == 0 {
			return nil, errors.New("cannot initiate a trade with neither an offer nor a request")
		}

		return &res, nil
	})
}
