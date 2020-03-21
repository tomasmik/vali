package vali

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

// TODO Consider exposing this to the package user
type cmp struct {
	float  func(have float64, exp []interface{}) (bool, error)
	int    func(have int64, exp []interface{}) (bool, error)
	uint   func(have uint64, exp []interface{}) (bool, error)
	time   func(have time.Time, exp []interface{}) (bool, error)
	string func(have string, exp []interface{}) (bool, error)
}

func newComparison(c *cmp) *cmp {
	return &cmp{
		float:  c.float,
		int:    c.int,
		uint:   c.uint,
		time:   c.time,
		string: c.string,
	}
}

func (c *cmp) do(s interface{}, o []interface{}) (bool, error) {
	if o == nil || len(o) == 0 {
		return false, errors.New("nothing to compare against")
	}
	if s == nil {
		return false, errors.New("can't compare if base value nil")
	}
	switch s.(type) {
	case float32, float64:
		have, _ := GetFloat(s)
		if c.float != nil {
			return c.float(have, o)
		}
		return true, nil
	case int, int8, int16, int32, int64:
		have, _ := GetInt(s)
		if c.int != nil {
			return c.int(have, o)
		}
		return true, nil
	case uint, uint8, uint16, uint32, uint64:
		have, _ := GetUInt(s)
		if c.uint != nil {
			return c.uint(have, o)
		}
		return true, nil
	case time.Time:
		have, _ := s.(time.Time)
		if c.time != nil {
			return c.time(have, o)
		}
		return true, nil
	case string:
		have, _ := s.(string)
		if c.string != nil {
			return c.string(have, o)
		}
		return true, nil
	default:
		return false, fmt.Errorf("can't check max of type %v", reflect.TypeOf(s).String())
	}
}
