package vali

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// defaultTags returns the list of all default tags
func defaultTags() tags {
	return map[string]tagFunc{
		"required":         required,
		"required_without": required_without,
		"less_than":        less_than,
		"more_than":        more_than,
		"optional":         optional,
	}
}

// required tag is used when default values can't be used
// in struct field values.
func required(s interface{}, o interface{}) error {
	v := s.(reflect.Value)
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

func required_without(s interface{}, o interface{}) error {
	sOK := required(s, nil) == nil

	if !sOK {
		v := o.([]reflect.Value)
		for _, f := range v {
			if err := required(f, nil); err != nil {
				return err
			}
		}
	}
	return nil
}

func optional(s interface{}, o interface{}) error {
	return nil
}

// TODO refactor this in to something readable
func more_than(s interface{}, o interface{}) error {
	v := s.(reflect.Value)
	if i, ok := getInt(v); ok {
		more_than, err := strconv.Atoi(o.([]string)[0])
		if err != nil {
			panic(err)
		}

		if i < int64(more_than) {
			return fmt.Errorf("%d is less than %d", i, more_than)
		}
		return nil
	}

	if i, ok := getFloat(v); ok {
		more_than, err := strconv.ParseFloat(o.([]string)[0], 64)
		if err != nil {
			panic(err)
		}

		if i < more_than {
			return fmt.Errorf("%f is less than %f", i, more_than)
		}
		return nil
	}

	return errors.New("value is not an int")
}

func less_than(s interface{}, o interface{}) error {
	v := s.(reflect.Value)
	switch v.Kind() {
	case reflect.Int:
		less_than, err := strconv.Atoi(o.([]string)[0])
		if err != nil {
			panic(err)
		}

		have := v.Int()
		if have > int64(less_than) {
			return fmt.Errorf("%d is more than %d", have, less_than)
		}
	case reflect.Float64:
		less_than, err := strconv.ParseFloat(o.([]string)[0], 64)
		if err != nil {
			panic(err)
		}

		have := v.Float()
		if have > less_than {
			return fmt.Errorf("%f is more than %f", have, less_than)
		}
	}
	return errors.New("value is not an int")
}
