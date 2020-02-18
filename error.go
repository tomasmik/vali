package vali

import (
	"fmt"
	"reflect"
)

// AggErr is a struct which allows the
// Validate func to stack errors in to a slice
// but return a single error.
type AggErr struct {
	Sl []error
}

func newAggErr() *AggErr {
	return &AggErr{
		Sl: make([]error, 0),
	}
}

func (e *AggErr) addErr(err ...error) *AggErr {
	e.Sl = append(e.Sl, err...)
	return e
}

func (e *AggErr) Error() string {
	var s string
	for i, err := range e.Sl {
		s = err.Error()
		if len(e.Sl)-1 != i {
			s += "\n"
		}
	}
	return s
}

func (e *AggErr) toError() error {
	if len(e.Sl) == 0 {
		return nil
	}
	return e
}

func tagError(field, tag string, err error) error {
	return fmt.Errorf("field: '%s', failed '%s' tag with an error: '%v'", field, tag, err)
}

func typeMismatch(i, o interface{}) error {
	return fmt.Errorf("argument with type %v cant be compared to value of type %v", reflect.TypeOf(o), reflect.TypeOf(i))
}
