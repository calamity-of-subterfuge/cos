package orders

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/mitchellh/mapstructure"
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

func decodeHook(from reflect.Value, to reflect.Value) (interface{}, error) {
	if from.Kind() == reflect.String {
		switch to.Kind() {
		case reflect.Int:
			return strconv.Atoi(from.String())
		case reflect.Int8:
			a, e := strconv.ParseInt(from.String(), 10, 8)
			return int8(a), e
		case reflect.Int16:
			a, e := strconv.ParseInt(from.String(), 10, 16)
			return int16(a), e
		case reflect.Int32:
			a, e := strconv.ParseInt(from.String(), 10, 32)
			return int32(a), e
		case reflect.Int64:
			return strconv.ParseInt(from.String(), 10, 64)
		case reflect.Float32:
			a, e := strconv.ParseFloat(from.String(), 32)
			return float32(a), e
		case reflect.Float64:
			return strconv.ParseFloat(from.String(), 64)
		}
	}
	return from.Interface(), nil
}

func parseSingleOrderOfType(parsed map[string]interface{}, typ Order) (Order, error) {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Squash:     true,
		DecodeHook: decodeHook,
		Result:     typ,
	})
	if err != nil {
		return nil, fmt.Errorf("constructing decoder: %w", err)
	}

	err = decoder.Decode(parsed)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal order: %w", err)
	}
	return typ, nil
}
