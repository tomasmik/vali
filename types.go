package vali

import (
	"fmt"
	"reflect"
)

func getInt(s interface{}) (int64, bool) {
	switch s.(type) {
	case int, int8, int16, int32, int64:
		return interfaceToReflectVal(s).Int(), true
	default:
		return 0, false
	}
}

func getUInt(s interface{}) (uint64, bool) {
	switch s.(type) {
	case uint, uint8, uint16, uint32, uint64:
		return interfaceToReflectVal(s).Uint(), true
	default:
		return 0, false
	}
}

func getFloat(s interface{}) (float64, bool) {
	switch s.(type) {
	case float32, float64:
		return interfaceToReflectVal(s).Float(), true
	default:
		return 0.0, false
	}
}

func getString(s interface{}) string {
	return fmt.Sprintf("%v", interfaceToReflectVal(s))
}

func getInterface(v reflect.Value) (interface{}, bool) {
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return nil, false
		}

		return getInterface(v.Elem())
	default:
		return v.Interface(), true
	}
}

func interfaceToReflectVal(s interface{}) reflect.Value {
	v, ok := s.(reflect.Value)
	if !ok {
		v = reflect.ValueOf(s)
	}
	return v
}

func derefInterface(s interface{}) interface{} {
	in, _ := getInterface(interfaceToReflectVal(s))
	return in
}
