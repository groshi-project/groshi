package internal

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// todo:
// complex64
// complex128
// array
// map
// slice of anything except []uint8
// (are they really needed?)

var booleanTrue = []string{"true", "on", "enable", "1", "yes", "ok"}
var booleanFalse = []string{"false", "off", "disable", "0", "no"}

func ParseBool(x string) (bool, error) {
	x = strings.ToLower(x)
	for _, trueString := range booleanTrue {
		if x == trueString {
			return true, nil
		}
	}
	for _, falseString := range booleanFalse {
		if x == falseString {
			return false, nil
		}
	}
	return false, errors.New("this string can't be represented as boolean value")
}

func ParseInt(x string, bitSize int, intKind reflect.Kind) (any, error) {
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

func ParseUnsignedInt(x string, bitSize int, uintKind reflect.Kind) (any, error) {
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

func ParseFloat(x string, bitSize int, floatKind reflect.Kind) (any, error) {
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

func ParseBytes(x string) ([]byte, error) {
	return []byte(x), nil
}

func Parse(value string, targetType reflect.Type) (any, error) {
	targetTypeKind := targetType.Kind()
	switch targetTypeKind {

	// boolean
	case reflect.Bool:
		return ParseBool(value)

	// signed int:
	case reflect.Int:
		return ParseInt(value, 64, targetTypeKind)
	case reflect.Int8:
		return ParseInt(value, 8, targetTypeKind)
	case reflect.Int16:
		return ParseInt(value, 16, targetTypeKind)
	case reflect.Int32:
		return ParseInt(value, 32, targetTypeKind)
	case reflect.Int64:
		return ParseInt(value, 64, targetTypeKind)

	// unsigned int:
	case reflect.Uint:
		return ParseUnsignedInt(value, 64, targetTypeKind)
	case reflect.Uint8:
		return ParseUnsignedInt(value, 8, targetTypeKind)
	case reflect.Uint16:
		return ParseUnsignedInt(value, 16, targetTypeKind)
	case reflect.Uint32:
		return ParseUnsignedInt(value, 32, targetTypeKind)
	case reflect.Uint64:
		return ParseUnsignedInt(value, 64, targetTypeKind)

	// float:
	case reflect.Float32:
		return ParseFloat(value, 32, targetTypeKind)
	case reflect.Float64:
		return ParseFloat(value, 64, targetTypeKind)

	// string:
	case reflect.String:
		return value, nil

	// slice:
	case reflect.Slice:
		switch targetType.String() {
		case "[]uint8":
			return ParseBytes(value)
		default:
			return nil, errors.New("slice of unimplemented type")
		}
	default:
		return nil, fmt.Errorf("unimplemented type %v (of kind %v)", targetType, targetTypeKind)
	}
}
