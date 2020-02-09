package vali

import (
	"fmt"
	"reflect"
)

// StErr is returned by the validation func `Validate()`
type StErr struct {
	Field string
	Tag   string
	Err   error
}

func newStErr(field, tag string, err error) *StErr {
	return &StErr{
		Field: field,
		Tag:   tag,
		Err:   err,
	}
}

// Error prints the error as a string
func (e *StErr) Error() string {
	return fmt.Sprintf("field: '%s', failed '%s' tag with an error: '%v'", e.Field, e.Tag, e.Err)
}

func typeMismatch(i, o interface{}) error {
	return fmt.Errorf("argument with type %v cant be compared to value of type %v", reflect.TypeOf(o), reflect.TypeOf(i))
}
