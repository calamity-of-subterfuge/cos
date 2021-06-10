package orders

import "errors"

// WithdrawTradeOrder is sent by the client to a tent to request that
// a trade offer be withdrawn. This works only on offers that were
// initiated by the clients team, to accept or reject offers that came
// from other teams use RespondTradeOrder
type WithdrawTradeOrder struct {
	// Type is always withdraw-trade
	Type string `json:"type" mapstructure:"type"`

	// UID of the trade offer to withdraw
	UID string `json:"uid" mapstructure:"uid"`
}

func (o *WithdrawTradeOrder) GetType() string {
	return "withdraw-trade"
}

func (o *WithdrawTradeOrder) PrepareForMarshal() {
	o.Type = o.GetType()
}

func init() {
	registerOrderParser("withdraw-trade", func(m map[string]interface{}) (Order, error) {
		var res WithdrawTradeOrder
		_, err := parseSingleOrderOfType(m, &res)
		if err != nil {
			return nil, err
		}

		if res.UID == "" {
			return nil, errors.New("missing UID")
		}

		return &res, nil
	})
}
