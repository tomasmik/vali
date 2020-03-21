package vali

import (
	"fmt"
	"reflect"
)

// GetInt is a safe way to convert an interface
// of any integer value to an int64 type integer.
func GetInt(s interface{}) (int64, bool) {
	switch s.(type) {
	case int, int8, int16, int32, int64:
		return interfaceToReflectVal(s).Int(), true
	default:
		return 0, false
	}
}

// GetInt is a safe way to convert an interface
// of any unsinged integer value to an uint64 type integer.
func GetUInt(s interface{}) (uint64, bool) {
	switch s.(type) {
	case uint, uint8, uint16, uint32, uint64:
		return interfaceToReflectVal(s).Uint(), true
	default:
		return 0, false
	}
}

// GetUInt is a safe way to convert an interface
// of any [unsinged] integer value to an uint64 type integer.
//
// It's different from GetUInt as if the given number is not
// a usinged integer it falls back to conveting it to an integer
// but still returning a uint.
// This can be useful when you're getting a integer value from a tag
// can will compare it to a uint. As the value from a tag will
// be a signed integer because of how GO works.
func GetUIntFallback(s interface{}) (uint64, bool) {
	switch s.(type) {
	case uint, uint8, uint16, uint32, uint64:
		return interfaceToReflectVal(s).Uint(), true
	default:
		val, ok := GetInt(s)
		if !ok {
			return 0, false
		}

		return uint64(val), true
	}
}

// GetFloat is a safe way to convert an interface
// of any floating point value to an float64 type floating point number.
func GetFloat(s interface{}) (float64, bool) {
	switch s.(type) {
	case float32, float64:
		return interfaceToReflectVal(s).Float(), true
	default:
		return 0.0, false
	}
}

// GetString convert any given interface to a string
// It uses the fmt packages `Sprintf` function
// which if speed is concerned is not the best fit.
func GetString(s interface{}) string {
	return fmt.Sprintf("%v", interfaceToReflectVal(s))
}

// DerefInterface is a safe way to dereference an interface.
// If the interface was not a pointer, it returns it back
// without panicking which makes the function possible to use
// when you dont know whetever or not the interface is a pointer.
func DerefInterface(s interface{}) interface{} {
	in, _ := getInterface(interfaceToReflectVal(s))
	return in
}

func interfaceToReflectVal(s interface{}) reflect.Value {
	v, ok := s.(reflect.Value)
	if !ok {
		v = reflect.ValueOf(s)
	}
	return v
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
