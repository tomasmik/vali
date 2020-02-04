package vali

import "reflect"

func getInt(v reflect.Value) (int64, bool) {
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return 0, false
		}
		return getInt(v.Elem())
	case reflect.Int, reflect.Uint, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Int(), true
	}
	return 0, false
}

func getFloat(v reflect.Value) (float64, bool) {
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return 0.0, false
		}
		return getFloat(v.Elem())
	case reflect.Float32, reflect.Float64:
		return v.Float(), true
	}
	return 0.0, false
}
