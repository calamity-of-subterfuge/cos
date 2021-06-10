package utils

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/mitchellh/mapstructure"
)

// DecodeHook is a valid mapstructure DecodeHook in order to correctly interpret
// json.Number, which is a subtype of string, as actual numbers in the target
// struct. This produces the correct error based on the target size (e.g.,
// the json number 256 can't be stored in a int8, so an error would be
// produced)
func DecodeHook(from reflect.Value, to reflect.Value) (interface{}, error) {
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
	} else if from.Kind() == reflect.Map && from.Type().Key().Kind() == reflect.String && from.Type().Elem().Kind() == reflect.String && to.Type().Key().Kind() == reflect.String {
		// from is map[string]string
		// to is map[string]SOMETHING

		switch to.Type().Elem().Kind() {
		case reflect.Int:
			res := make(map[string]int, from.Len())
			iter := from.MapRange()
			for iter.Next() {
				k := iter.Key()
				v := iter.Value()

				parsed, err := strconv.Atoi(v.String())
				if err != nil {
					return nil, err
				}
				res[k.String()] = parsed
			}
			return res, nil
		case reflect.Int8:
			res := make(map[string]int8, from.Len())
			iter := from.MapRange()
			for iter.Next() {
				k := iter.Key()
				v := iter.Value()

				parsed, err := strconv.ParseInt(v.String(), 10, 8)
				if err != nil {
					return nil, err
				}
				res[k.String()] = int8(parsed)
			}
			return res, nil
		case reflect.Int16:
			res := make(map[string]int16, from.Len())
			iter := from.MapRange()
			for iter.Next() {
				k := iter.Key()
				v := iter.Value()

				parsed, err := strconv.ParseInt(v.String(), 10, 16)
				if err != nil {
					return nil, err
				}
				res[k.String()] = int16(parsed)
			}
			return res, nil
		case reflect.Int32:
			res := make(map[string]int32, from.Len())
			iter := from.MapRange()
			for iter.Next() {
				k := iter.Key()
				v := iter.Value()

				parsed, err := strconv.ParseInt(v.String(), 10, 32)
				if err != nil {
					return nil, err
				}
				res[k.String()] = int32(parsed)
			}
			return res, nil
		case reflect.Int64:
			res := make(map[string]int64, from.Len())
			iter := from.MapRange()
			for iter.Next() {
				k := iter.Key()
				v := iter.Value()

				parsed, err := strconv.ParseInt(v.String(), 10, 64)
				if err != nil {
					return nil, err
				}
				res[k.String()] = parsed
			}
			return res, nil
		case reflect.Float32:
			res := make(map[string]float32, from.Len())
			iter := from.MapRange()
			for iter.Next() {
				k := iter.Key()
				v := iter.Value()

				parsed, err := strconv.ParseFloat(v.String(), 32)
				if err != nil {
					return nil, err
				}
				res[k.String()] = float32(parsed)
			}
			return res, nil
		case reflect.Float64:
			res := make(map[string]float64, from.Len())
			iter := from.MapRange()
			for iter.Next() {
				k := iter.Key()
				v := iter.Value()

				parsed, err := strconv.ParseFloat(v.String(), 64)
				if err != nil {
					return nil, err
				}
				res[k.String()] = parsed
			}
			return res, nil
		}
	}
	return from.Interface(), nil
}

// DecodeWithType is a convenience method to use a mapstructure Decoder
// using the DecodeHook to decode the given map with the given target type.
func DecodeWithType(parsed map[string]interface{}, typ interface{}) (interface{}, error) {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Squash:     true,
		DecodeHook: DecodeHook,
		Result:     typ,
	})
	if err != nil {
		return nil, fmt.Errorf("constructing decoder: %w", err)
	}

	err = decoder.Decode(parsed)
	if err != nil {
		return nil, fmt.Errorf("failed to decode: %w", err)
	}
	return typ, nil
}
