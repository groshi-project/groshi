package internal

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// todo (is this really needed?):
// slice of anything except []uint8
// array
// map
// complex64
// complex128

var booleanValues = map[string]bool{
	"true":   true,
	"on":     true,
	"enable": true,
	"1":      true,
	"yes":    true,
	"ok":     true,

	"false":   false,
	"off":     false,
	"disable": false,
	"0":       false,
	"no":      false,
}

func parseBool(x string) (bool, error) {
	value, found := booleanValues[strings.ToLower(x)]
	if !found {
		return false, errors.New("this string can't be represented as boolean value")
	}
	return value, nil
}

func parseInt(x string, bitSize int, intKind reflect.Kind) (any, error) {
	i64, err := strconv.ParseInt(x, 10, bitSize)
	if err != nil {
		return 0, err
	}
	switch intKind {
	case reflect.Int:
		return int(i64), nil
	case reflect.Int8:
		return int8(i64), nil
	case reflect.Int16:
		return int16(i64), nil
	case reflect.Int32:
		return int32(i64), nil
	case reflect.Int64:
		return i64, nil
	default:
		panic(fmt.Errorf("unknown int kind %v", intKind))
	}
}

func parseUnsignedInt(x string, bitSize int, uintKind reflect.Kind) (any, error) {
	ui64, err := strconv.ParseUint(x, 10, bitSize)
	if err != nil {
		return 0, err
	}
	switch uintKind {
	case reflect.Uint:
		return uint(ui64), nil
	case reflect.Uint8:
		return uint8(ui64), nil
	case reflect.Uint16:
		return uint16(ui64), nil
	case reflect.Uint32:
		return uint32(ui64), nil
	case reflect.Uint64:
		return ui64, nil
	default:
		panic(fmt.Errorf("unknown uint kind %v", uintKind))
	}
}

func parseFloat(x string, bitSize int, floatKind reflect.Kind) (any, error) {
	f64, err := strconv.ParseFloat(x, bitSize)
	if err != nil {
		return 0.0, err
	}
	switch floatKind {
	case reflect.Float32:
		return float32(f64), nil
	case reflect.Float64:
		return f64, nil
	default:
		panic(fmt.Errorf("unknown float kind %v", floatKind))
	}
}

func parseSliceOfBytes(x string) ([]byte, error) {
	return []byte(x), nil
}

func ParseValue(x string, targetType reflect.Type) (any, error) {
	targetTypeKind := targetType.Kind()
	switch targetTypeKind {

	// boolean:
	case reflect.Bool:
		return parseBool(x)

	// signed int:
	case reflect.Int:
		return parseInt(x, 64, targetTypeKind)
	case reflect.Int8:
		return parseInt(x, 8, targetTypeKind)
	case reflect.Int16:
		return parseInt(x, 16, targetTypeKind)
	case reflect.Int32:
		return parseInt(x, 32, targetTypeKind)
	case reflect.Int64:
		return parseInt(x, 64, targetTypeKind)

	// unsigned int:
	case reflect.Uint:
		return parseUnsignedInt(x, 64, targetTypeKind)
	case reflect.Uint8:
		return parseUnsignedInt(x, 8, targetTypeKind)
	case reflect.Uint16:
		return parseUnsignedInt(x, 16, targetTypeKind)
	case reflect.Uint32:
		return parseUnsignedInt(x, 32, targetTypeKind)
	case reflect.Uint64:
		return parseUnsignedInt(x, 64, targetTypeKind)

	// float:
	case reflect.Float32:
		return parseFloat(x, 32, targetTypeKind)
	case reflect.Float64:
		return parseFloat(x, 64, targetTypeKind)

	// string:
	case reflect.String:
		return x, nil

	// slice:
	case reflect.Slice:
		sliceType := targetType.String()
		if sliceType == "[]uint8" {
			return parseSliceOfBytes(x)
		}
		fallthrough
	default:
		return nil, fmt.Errorf("unimplemented type %v (of kind %v)", targetType, targetTypeKind)
	}
}
