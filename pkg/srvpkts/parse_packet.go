package srvpkts

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/calamity-of-subterfuge/cos/pkg/utils"
)

type packetParser func(map[string]interface{}) (Packet, error)

var packetParsersByType map[string]packetParser = make(map[string]packetParser)

func registerPacketParser(typ string, parser packetParser) {
	packetParsersByType[typ] = parser
}

// ParsePacket attempts to parse the packet described by the given bytes as one
// of the packets in this package. The resulting interface is nil if error is
// not nil, otherwise its a pointer to one of the Packet structs and should be
// switched on by type.
func ParsePacket(raw []byte) ([]Packet, error) {
	var parsed interface{}
	err := json.Unmarshal(raw, &parsed)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if arr, ok := parsed.([]interface{}); ok {
		result := make([]Packet, 0, len(arr))
		for _, ele := range arr {
			var eleParsed map[string]interface{}
			eleParsed, ok = ele.(map[string]interface{})

			if !ok {
				return nil, fmt.Errorf("element in array is not a JSON object: %v", ele)
			}

			var interpreted Packet
			interpreted, err = ParseSinglePacket(eleParsed)
			if err != nil {
				return nil, fmt.Errorf("element in array not valid: %w", err)
			}

			result = append(result, interpreted)
		}
		return result, nil
	} else if asMap, ok := parsed.(map[string]interface{}); ok {
		result := make([]Packet, 1)

		var interpreted Packet
		interpreted, err = ParseSinglePacket(asMap)
		if err != nil {
			return nil, fmt.Errorf("failed interpret JSON object: %w", err)
		}

		result[0] = interpreted
		return result, nil
	}

	return nil, fmt.Errorf("valid json but not a json object or array of json objects")
}

// ParseSinglePacket parses a single packet which has already been
// converted to the map interpretation. Note that many of these may
// be sent within a single websocket message frame, formatted as a
// JSON array. In order to parse from the raw bytes of a websocket message
// frame, use ParsePacket.
func ParseSinglePacket(parsed map[string]interface{}) (Packet, error) {
	packetTypeRaw, found := parsed["type"]
	if !found {
		return nil, errors.New("packet missing type")
	}

	packetType, ok := packetTypeRaw.(string)
	if !ok {
		return nil, errors.New("packet has type but it's not a string")
	}

	var parser packetParser
	parser, found = packetParsersByType[packetType]
	if !found {
		return nil, fmt.Errorf("unknown packet type: %s", packetType)
	}

	return parser(parsed)
}

func parseSinglePacketOfType(parsed map[string]interface{}, typ Packet) (Packet, error) {
	_, err := utils.DecodeWithType(parsed, typ)
	if err != nil {
		return nil, err
	}
	return typ, nil
}
