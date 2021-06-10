package orders

import (
	"errors"
	"fmt"

	"github.com/calamity-of-subterfuge/cos/pkg/utils"
)

type orderParser func(map[string]interface{}) (Order, error)

var orderParsersByType map[string]orderParser = make(map[string]orderParser)

func registerOrderParser(typ string, parser orderParser) {
	orderParsersByType[typ] = parser
}

// ParseOrder parses a single order from its map representation within a packet.
// This typically comes from the clipkts.IssueSmartObjectOrderPacket#Order
func ParseOrder(parsed map[string]interface{}) (Order, error) {
	orderTypeRaw, found := parsed["type"]
	if !found {
		return nil, errors.New("order missing type")
	}

	orderType, ok := orderTypeRaw.(string)
	if !ok {
		return nil, errors.New("order has type but it's not a string")
	}

	var parser orderParser
	parser, found = orderParsersByType[orderType]
	if !found {
		return nil, fmt.Errorf("unknown order type: %s", orderType)
	}

	return parser(parsed)
}

func parseSingleOrderOfType(parsed map[string]interface{}, typ Order) (Order, error) {
	_, err := utils.DecodeWithType(parsed, typ)
	if err != nil {
		return nil, err
	}
	return typ, err
}
