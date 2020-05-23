package vali

import (
	"fmt"
	"reflect"
)

type tagError struct {
	field string
	tag   string
	err   error
}

func newTagError(field, tag string, err error) error {
	return &tagError{
		field: field,
		tag:   tag,
		err:   err,
	}
}

func (t *tagError) Error() string {
	return fmt.Sprintf("field: '%s', failed '%s' tag with an error: '%v'", t.field, t.tag, t.err)
}

func typeMismatch(i, o interface{}) error {
	return fmt.Errorf("argument with type %v cant be compared to value of type %v", reflect.TypeOf(o), reflect.TypeOf(i))
}
