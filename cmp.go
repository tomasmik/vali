package vali

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

// TODO Consider exposing this to the package user
type cmp struct {
	float    func(have float64, exp []interface{}) (bool, error)
	int      func(have int64, exp []interface{}) (bool, error)
	uint     func(have uint64, exp []interface{}) (bool, error)
	time     func(have time.Time, exp []interface{}) (bool, error)
	string   func(have string, exp []interface{}) (bool, error)
	duration func(have time.Duration, exp []interface{}) (bool, error)
	slice    func(have interface{}, exp []interface{}) (bool, error)
}

func (c *cmp) do(s interface{}, o []interface{}) (bool, error) {
	if o == nil {
		return false, errors.New("nothing to compare against, [o] is nil")
	}
	if s == nil {
		return false, errors.New("can't compare if base value [s] nil")
	}

	err := func(s string) error {
		return fmt.Errorf("comparison for %v is not set", s)
	}
	// Check using reflect for types that are not easy to cast using an interface
	switch reflect.ValueOf(s).Kind() {
	case reflect.Array, reflect.Slice:
		if c.slice != nil {
			return c.slice(s, o)
		}
		return false, err(reflect.TypeOf(s).String())
	}

	// Validate types that are easy to cast from an interface
	switch s.(type) {
	case float32, float64:
		have, _ := GetFloat(s)
		if c.float != nil {
			return c.float(have, o)
		}
		return false, err(reflect.TypeOf(s).String())
	case int, int8, int16, int32, int64:
		have, _ := GetInt(s)
		if c.int != nil {
			return c.int(have, o)
		}
		return false, err(reflect.TypeOf(s).String())
	case uint, uint8, uint16, uint32, uint64:
		have, _ := GetUInt(s)
		if c.uint != nil {
			return c.uint(have, o)
		}
		return false, err(reflect.TypeOf(s).String())
	case time.Time:
		have, _ := s.(time.Time)
		if c.time != nil {
			return c.time(have, o)
		}
		return false, err(reflect.TypeOf(s).String())
	case string:
		have, _ := s.(string)
		if c.string != nil {
			return c.string(have, o)
		}
		return false, err(reflect.TypeOf(s).String())
	case time.Duration:
		have, _ := s.(time.Duration)
		if c.duration != nil {
			return c.duration(have, o)
		}
		return false, err(reflect.TypeOf(s).String())
	default:
		return false, fmt.Errorf("comparing type %v to other values is not supported", reflect.TypeOf(s).String())
	}
}

func newOneOfCMP(s interface{}, o []interface{}) *cmp {
	return &cmp{
		float: func(have float64, exp []interface{}) (bool, error) {
			for _, arg := range exp {
				f, ok := GetFloat(arg)
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
				f, ok := GetInt(arg)
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
				f, ok := GetUIntFallback(arg)
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
				f := GetString(arg)
				if have == f {
					return true, nil
				}
			}
			return false, nil
		},
	}

}

func newMinCMP(s interface{}, o []interface{}) *cmp {
	return &cmp{
		float: func(have float64, exp []interface{}) (bool, error) {
			val, _ := GetFloat(s)
			more, ok := GetFloat(exp[0])
			if !ok {
				return false, typeMismatch(s, exp[0])
			}
			return val > more, nil
		},
		int: func(have int64, exp []interface{}) (bool, error) {
			val, _ := GetInt(s)
			more, ok := GetInt(exp[0])
			if !ok {
				return false, typeMismatch(s, exp[0])
			}
			return val > more, nil
		},
		uint: func(have uint64, exp []interface{}) (bool, error) {
			val, _ := GetUInt(s)
			more, ok := GetUIntFallback(exp[0])
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
			more, ok := GetInt(exp[0])
			if !ok {
				return false, typeMismatch(s, exp[0])
			}

			return !(int64(len(have)) < more), nil
		},
		slice: func(have interface{}, exp []interface{}) (bool, error) {
			sl := reflect.ValueOf(have)
			more, ok := GetUIntFallback(exp[0])
			if !ok {
				return false, typeMismatch(s, exp[0])
			}

			return sl.Len() > int(more), nil
		},
	}
}

func newEqualsCMP(s interface{}, o []interface{}) *cmp {
	return &cmp{
		float: func(have float64, exp []interface{}) (bool, error) {
			f, ok := GetFloat(exp[0])
			if !ok {
				return false, typeMismatch(s, exp[0])
			}
			return have == f, nil
		},
		int: func(have int64, exp []interface{}) (bool, error) {
			f, ok := GetInt(exp[0])
			if !ok {
				return false, typeMismatch(s, exp[0])
			}
			return have == f, nil
		},
		uint: func(have uint64, exp []interface{}) (bool, error) {
			f, ok := GetUIntFallback(exp[0])
			if !ok {
				return false, typeMismatch(s, exp[0])

			}
			return have == f, nil
		},
		string: func(have string, exp []interface{}) (bool, error) {
			return have == GetString(exp[0]), nil
		},
		slice: func(have interface{}, exp []interface{}) (bool, error) {
			sl := reflect.ValueOf(have)
			more, ok := GetUIntFallback(exp[0])
			if !ok {
				return false, typeMismatch(s, exp[0])
			}

			return sl.Len() == int(more), nil
		},
	}
}
