package vali

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

const (
	// requiredTag can be used to tag a struct field making
	// it fail the validation if the value is nil or default value.
	// Only 1 required tag is allowed and a field cannot have
	// a mix of required, required_without, optional.
	requiredTag = "required"
	// requiredWihtouTag can be used to tag a struct field making
	// it fail the validation if the value nil or default value and the field pointer is nil.
	// Only 1 required_without tag is allowed and a field cannot have
	// a mix of required, required_without, optional.
	requiredWithoutTag = "required_without"
	// optionalTag can be used to tag a struct field making
	// it not fail validation if it's empty or nil.
	// Only 1 optional tag is allowed and a field cannot have
	// a mix of required, required_without, optional.
	//
	// This is a special tag as it's not a part of the tags map
	optionalTag = "optional"
	// maxTag can be used to tag a struct field making
	// it fail validation if the field is more than max.
	maxTag = "max"
	// maxTag can be used to tag a struct field making
	// it fail validation if the field is less than min.
	minTag = "min"
	// oneofTag can be used to tag a struct field making
	// it fail validation if the field does not have one of the
	// given values.
	oneofTag = "one_of"
	// eqTag can be used to tag a struct field making
	// it fail validation if the field is not equal to a given value or field.
	eqTag = "eq"
	// neqTag can be used to tag a struct field making
	// it fail validation if the field is equal to a given value or field.
	neqTag = "neq"
)

func required(s interface{}, o []interface{}) error {
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
		if s == nil {
			return errors.New("value is nil")
		}

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

func min(s interface{}, o []interface{}) error {
	comparison := newComparison(&cmp{
		float: func(have float64, exp []interface{}) (bool, error) {
			val, _ := getFloat(s)
			more, ok := getFloat(exp[0])
			if !ok {
				return false, typeMismatch(s, exp[0])
			}
			return val > more, nil
		},
		int: func(have int64, exp []interface{}) (bool, error) {
			val, _ := getInt(s)
			more, ok := getInt(exp[0])
			if !ok {
				return false, typeMismatch(s, exp[0])
			}
			return val > more, nil
		},
		uint: func(have uint64, exp []interface{}) (bool, error) {
			val, _ := getUInt(s)
			more, ok := getUIntFallback(exp[0])
			if !ok {
				return false, typeMismatch(s, exp[0])
			}
			return val > more, nil
		},
		time: func(have time.Time, exp []interface{}) (bool, error) {
			more, ok := exp[0].(time.Time)
			if !ok {
				return false, typeMismatch(s, exp[0])
			}
			return !have.Before(more), nil
		},
		string: func(have string, exp []interface{}) (bool, error) {
			more, ok := getInt(exp[0])
			if !ok {
				return false, typeMismatch(s, exp[0])
			}

			return !(int64(len(have)) < more), nil
		},
	})

	ok, err := comparison.do(s, o)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("%v is less than %v", s, o[0])
	}

	return nil
}

func max(s interface{}, o []interface{}) error {
	comparison := newComparison(&cmp{
		float: func(have float64, exp []interface{}) (bool, error) {
			val, _ := getFloat(s)
			more, ok := getFloat(exp[0])
			if !ok {
				return false, typeMismatch(s, exp[0])
			}
			return val <= more, nil
		},
		int: func(have int64, exp []interface{}) (bool, error) {
			val, _ := getInt(s)
			more, ok := getInt(exp[0])
			if !ok {
				return false, typeMismatch(s, exp[0])
			}
			return val <= more, nil
		},
		uint: func(have uint64, exp []interface{}) (bool, error) {
			val, _ := getUInt(s)
			more, ok := getUIntFallback(exp[0])
			if !ok {
				return false, typeMismatch(s, exp[0])
			}
			return val <= more, nil
		},
		time: func(have time.Time, exp []interface{}) (bool, error) {
			more, ok := exp[0].(time.Time)
			if !ok {
				return false, typeMismatch(s, exp[0])
			}
			return have.Before(more), nil
		},
		string: func(have string, exp []interface{}) (bool, error) {
			more, ok := getInt(exp[0])
			if !ok {
				return false, typeMismatch(s, exp[0])
			}

			return (int64(len(have)) < more), nil
		},
	})

	ok, err := comparison.do(s, o)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("%v is more than %v", s, o[0])
	}

	return nil
}

func oneof(s interface{}, o []interface{}) error {
	comparison := newComparison(&cmp{
		float: func(have float64, exp []interface{}) (bool, error) {
			for _, arg := range exp {
				f, ok := getFloat(arg)
				if !ok {
					return false, typeMismatch(s, arg)
				}
				if have == f {
					return true, nil
				}
			}
			return false, nil
		},
		int: func(have int64, exp []interface{}) (bool, error) {
			for _, arg := range exp {
				f, ok := getInt(arg)
				if !ok {
					return false, typeMismatch(s, arg)
				}
				if have == f {
					return true, nil
				}
			}
			return false, nil
		},
		uint: func(have uint64, exp []interface{}) (bool, error) {
			for _, arg := range exp {
				f, ok := getUIntFallback(arg)
				if !ok {
					return false, typeMismatch(s, arg)

				}
				if have == f {
					return true, nil
				}
			}
			return false, nil
		},
		string: func(have string, exp []interface{}) (bool, error) {
			for _, arg := range o {
				f := getString(arg)
				if have == f {
					return true, nil
				}
			}
			return false, nil
		},
	})

	ok, err := comparison.do(s, o)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("must have at least one of oneof %v", o)
	}

	return nil
}

func eq(s interface{}, o []interface{}) error {
	comparison := newComparison(&cmp{
		float: func(have float64, exp []interface{}) (bool, error) {
			f, ok := getFloat(exp[0])
			if !ok {
				return false, typeMismatch(s, exp[0])
			}
			return have == f, nil
		},
		int: func(have int64, exp []interface{}) (bool, error) {
			f, ok := getInt(exp[0])
			if !ok {
				return false, typeMismatch(s, exp[0])
			}
			return have == f, nil
		},
		uint: func(have uint64, exp []interface{}) (bool, error) {
			f, ok := getUIntFallback(exp[0])
			if !ok {
				return false, typeMismatch(s, exp[0])

			}
			return have == f, nil
		},
		string: func(have string, exp []interface{}) (bool, error) {
			return have == getString(exp[0]), nil
		},
	})

	ok, err := comparison.do(s, o)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("%v is not equal to %v", s, o[0])
	}

	return nil
}

func neq(s interface{}, o []interface{}) error {
	comparison := newComparison(&cmp{
		float: func(have float64, exp []interface{}) (bool, error) {
			f, ok := getFloat(exp[0])
			if !ok {
				return false, typeMismatch(s, exp[0])
			}
			return have != f, nil
		},
		int: func(have int64, exp []interface{}) (bool, error) {
			f, ok := getInt(exp[0])
			if !ok {
				return false, typeMismatch(s, exp[0])
			}
			return have != f, nil
		},
		uint: func(have uint64, exp []interface{}) (bool, error) {
			f, ok := getUIntFallback(exp[0])
			if !ok {
				return false, typeMismatch(s, exp[0])

			}
			return have != f, nil
		},
		string: func(have string, exp []interface{}) (bool, error) {
			return have != getString(exp[0]), nil
		},
	})

	ok, err := comparison.do(s, o)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("%v is equal to %v", s, o[0])
	}

	return nil
}
