package vali

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

// defaultTags returns the list of all default tags
func defaultTags() tags {
	return map[string]TagFunc{
		"required":         required,
		"required_without": required_without,
		"max":              max,
		"min":              min,
		"one_of":           oneof,
		"eq":               eq,
		"neq":              neq,
		"optional":         optional,
	}
}

// required tag is used when default values can't be used
// in struct field values.
func required(s interface{}, o []interface{}) error {
	if s == nil {
		return errors.New("value is nil")
	}

	v := interfaceToReflectVal(s)
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return errors.New("pointer is nil")
		}

		return required(v.Elem(), o)
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			// Ignore unexported fields
			if v.Field(i).CanSet() {
				if err := required(v.Field(i), o); err != nil {
					return err
				}
			}
		}
	case reflect.Func, reflect.Map, reflect.Slice:
		if v.IsNil() {
			return fmt.Errorf("%s is nil", v.Type().String())
		}
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if err := required(v.Index(i), o); err != nil {
				return err
			}
		}
	default:
		z := reflect.Zero(v.Type())
		if v.Interface() == z.Interface() {
			return fmt.Errorf("empty %s", v.Type().String())
		}
	}

	return nil
}

func required_without(s interface{}, o []interface{}) error {
	sOK := required(s, nil) == nil

	if !sOK {
		for _, f := range o {
			if err := required(f, nil); err != nil {
				return err
			}
		}
	}
	return nil
}

func optional(s interface{}, o []interface{}) error {
	return nil
}

func min(s interface{}, o []interface{}) error {
	ov := o[0]
	if ov == nil {
		return errors.New("can't compare if arg is nil")
	}
	if s == nil {
		return errors.New("can't compare if base value nil")
	}

	switch s.(type) {
	case float32, float64:
		val, _ := getFloat(s)
		more, ok := getFloat(ov)
		if !ok {
			return typeMismatch(s, ov)
		}
		if val < more {
			return fmt.Errorf("%f is less than %f", val, more)
		}
	case int, int8, int16, int32, int64:
		val, _ := getInt(s)
		more, ok := getInt(ov)
		if !ok {
			return typeMismatch(s, ov)
		}
		if val < more {
			return fmt.Errorf("%d is less than %d", val, more)
		}
	case uint, uint8, uint16, uint32, uint64:
		val, _ := getUInt(s)
		more, ok := getUInt(ov)
		if !ok {
			return typeMismatch(s, ov)
		}
		if val < more {
			return fmt.Errorf("%d is less than %d", val, more)
		}
	case time.Time:
		val, _ := s.(time.Time)
		more, ok := ov.(time.Time)
		if !ok {
			return typeMismatch(s, ov)
		}
		if val.Before(more) {
			return fmt.Errorf("%s is less than %s", val.String(), more.String())
		}
	case string:
		val, _ := s.(string)
		more, ok := getInt(ov)
		if !ok {
			return typeMismatch(s, ov)
		}

		if int64(len(val)) < more {
			return fmt.Errorf("%s length is less than %d", val, more)
		}
	default:
		return fmt.Errorf("can't check min of type %v", reflect.TypeOf(s).String())
	}
	return nil
}

func max(s interface{}, o []interface{}) error {
	ov := o[0]
	if ov == nil {
		return errors.New("can't compare if arg is nil")
	}
	if s == nil {
		return errors.New("can't compare if base value nil")
	}

	switch s.(type) {
	case float32, float64:
		val, _ := getFloat(s)
		less, ok := getFloat(ov)
		if !ok {
			return typeMismatch(s, ov)
		}
		if val > less {
			return fmt.Errorf("%f is more than %f", val, less)
		}
	case int, int8, int16, int32, int64:
		val, _ := getInt(s)
		less, ok := getInt(ov)
		if !ok {
			return typeMismatch(s, ov)
		}
		if val > less {
			return fmt.Errorf("%d is more than %d", val, less)
		}
	case uint, uint8, uint16, uint32, uint64:
		val, _ := getUInt(s)
		less, ok := getUInt(ov)
		if !ok {
			return typeMismatch(s, ov)
		}
		if val > less {
			return fmt.Errorf("%d is more than %d", val, less)
		}
	case time.Time:
		val, _ := s.(time.Time)
		less, ok := ov.(time.Time)
		if !ok {
			return typeMismatch(s, ov)
		}
		if val.After(less) {
			return fmt.Errorf("%s is more than %s", val.String(), less.String())
		}
	case string:
		val, _ := s.(string)
		less, ok := getInt(ov)
		if !ok {
			return typeMismatch(s, ov)
		}

		if int64(len(val)) < less {
			return fmt.Errorf("%s length is more than %d", val, less)
		}
	default:
		return fmt.Errorf("can't check max of type %v", reflect.TypeOf(s).String())
	}
	return nil
}

func oneof(s interface{}, o []interface{}) error {
	if len(o) == 0 {
		return errors.New("no arguments passed")
	}

	switch s.(type) {
	case float32, float64:
		val, _ := getFloat(s)
		for _, arg := range o {
			f, ok := getFloat(arg)
			if !ok {
				return typeMismatch(s, arg)
			}
			if val == f {
				return nil
			}
		}
	case int, int8, int16, int32, int64:
		val, _ := getInt(s)
		for _, arg := range o {
			f, ok := getInt(arg)
			if !ok {
				return typeMismatch(s, arg)
			}
			if val == f {
				return nil
			}
		}
	case uint, uint8, uint16, uint32, uint64:
		val, _ := getUInt(s)
		for _, arg := range o {
			f, ok := getUInt(arg)
			if !ok {
				return typeMismatch(s, arg)
			}
			if val == f {
				return nil
			}
		}
	case string:
		val, _ := s.(string)
		for _, arg := range o {
			f := getString(arg)
			if val == f {
				return nil
			}
		}
	default:
		return fmt.Errorf("can't check oneof of type %v", reflect.TypeOf(s).String())
	}
	return fmt.Errorf("must have at least one of oneof %v", o)
}

func eq(s interface{}, o []interface{}) error {
	if len(o) == 0 {
		return errors.New("no arguments passed")
	}

	arg := o[0]
	switch s.(type) {
	case float32, float64:
		val, _ := getFloat(s)
		f, ok := getFloat(arg)
		if !ok {
			return typeMismatch(s, arg)
		}
		if val == f {
			return nil
		}
	case int, int8, int16, int32, int64:
		val, _ := getInt(s)
		f, ok := getInt(arg)
		if !ok {
			return typeMismatch(s, arg)
		}
		if val == f {
			return nil
		}
	case uint, uint8, uint16, uint32, uint64:
		val, _ := getUInt(s)
		f, ok := getUInt(arg)
		if !ok {
			return typeMismatch(s, arg)
		}
		if val == f {
			return nil
		}
	case string:
		val, _ := s.(string)
		f := getString(arg)
		if val == f {
			return nil
		}
	default:
		return fmt.Errorf("can't check eq of type %v", reflect.TypeOf(s).String())
	}
	return fmt.Errorf("%v is not equal to %v", s, arg)
}

func neq(s interface{}, o []interface{}) error {
	if len(o) == 0 {
		return errors.New("no arguments passed")
	}

	arg := o[0]
	switch s.(type) {
	case float32, float64:
		val, _ := getFloat(s)
		f, ok := getFloat(arg)
		if !ok {
			return typeMismatch(s, arg)
		}
		if val != f {
			return nil
		}
	case int, int8, int16, int32, int64:
		val, _ := getInt(s)
		f, ok := getInt(arg)
		if !ok {
			return typeMismatch(s, arg)
		}
		if val != f {
			return nil
		}
	case uint, uint8, uint16, uint32, uint64:
		val, _ := getUInt(s)
		f, ok := getUInt(arg)
		if !ok {
			return typeMismatch(s, arg)
		}
		if val != f {
			return nil
		}
	case string:
		val, _ := s.(string)
		f := getString(arg)
		if val != f {
			return nil
		}
	default:
		return fmt.Errorf("can't check eq of type %v", reflect.TypeOf(s).String())
	}

	return fmt.Errorf("%v is equal to %v", s, arg)
}
