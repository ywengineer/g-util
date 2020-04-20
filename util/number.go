package util

import (
	"errors"
	"reflect"
	"strconv"
)

//Bool
//Int
//Int8
//Int16
//Int32
//Int64
//Uint
//Uint8
//Uint16
//Uint32
//Uint64
//Uintptr
func Int2String(n interface{}) (string, error) {
	if n == nil {
		return "", errors.New("nil value detected")
	}
	t := reflect.TypeOf(n).Kind()
	switch t {
	case reflect.Bool:
		if n.(bool) {
			return "1", nil
		} else {
			return "0", nil
		}
	case reflect.Int:
		return strconv.Itoa(n.(int)), nil
	case reflect.Int8:
		return strconv.FormatUint(uint64(n.(int)), 10), nil
	case reflect.Int16:
		return strconv.FormatUint(uint64(n.(int16)), 10), nil
	case reflect.Int32:
		return strconv.FormatUint(uint64(n.(int32)), 10), nil
	case reflect.Int64:
		return strconv.FormatUint(uint64(n.(int64)), 10), nil
	case reflect.Uint:
		return strconv.FormatUint(uint64(n.(uint)), 10), nil
	case reflect.Uint8:
		return strconv.FormatUint(uint64(n.(uint8)), 10), nil
	case reflect.Uint16:
		return strconv.FormatUint(uint64(n.(uint16)), 10), nil
	case reflect.Uint32:
		return strconv.FormatUint(uint64(n.(uint32)), 10), nil
	case reflect.Uint64:
		return strconv.FormatUint(n.(uint64), 10), nil
	case reflect.Uintptr:
		return Int2String(reflect.ValueOf(n).Elem().Interface())
	default:
		return "", errors.New("detect an non-numeric type : " + t.String())
	}
}
