package vali

import "fmt"

type TagErr struct {
	Tag string
	Err error
}

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

func newTagErr(tag string, err error) *TagErr {
	return &TagErr{
		Tag: tag,
		Err: err,
	}
}

func (e *StErr) Error() string {
	return fmt.Sprintf("field: '%s', failed '%s' tag with an error: '%v'", e.Field, e.Tag, e.Err)
}

func (e *TagErr) Error() string {
	return fmt.Sprintf("tag '%s' failed with error: '%v'", e.Tag, e.Err)
}
