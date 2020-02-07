package vali

import (
	"errors"
	"fmt"
	"reflect"
)

const (
	// valiTag is a tag we look for in a struct field.
	// If the field has this tag, then we validate it.
	valiTag = "vali"
	// equals is a sign used to add values to a tag.
	equalsSep = "="
	// valueSep is used to seperate between values for a tag
	valueSep = ","
	// tagSep is used to seperate tags from each other
	tagSep = "|"
	// pointerToField to a field is used when you want to compare
	// one field to another.
	// For example:
	//  type mock struct {
	//  Str  string `vali:"-"`
	//  Str2 string `vali:"required_without=*Str"`
	//  }
	// The *Str points to another struct value, so there values
	// will be compared.
	pointerToField = "*"
)

// tags if a type of map which holds all tag validation funcs
type tags map[string]TagFunc

// TagFunc is a func that is used to validate a field `s`
// using data providing in slice `o`.
//
// `s` can be nil if `Config.IgnoreNilPointer` is set to true
// Slice `o` can be empty or nil and doesn't have to be used if not needed.
type TagFunc func(s interface{}, o []interface{}) error

// Vali is a struct that holds all the configuration
// for the validation tool.
//
// It can be extended using its public methods by adding
// additional tag validation funcs or type validation funcs.
// Just make sure you add them before actually using the `Validate` method
// as no thread safety exists, so using and editting validation funcs
// will result in a race condition.
type Vali struct {
	tags tags
}

// New returns a new validator instance.
func New() *Vali {
	return &Vali{
		tags: map[string]TagFunc{
			requiredTag:        required,
			requiredWithoutTag: required_without,
			maxTag:             max,
			minTag:             min,
			oneofTag:           oneof,
			eqTag:              eq,
			neqTag:             neq,
			optionalTag:        optional,
		},
	}
}

// Validate accepts a struct and validates its according to the given tags.
func (v *Vali) Validate(s interface{}) []error {
	var errs []error

	if s == nil {
		return append(errs, errors.New("struct is nil"))
	}

	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// we only accept structs
	if val.Kind() != reflect.Struct {
		return append(errs, fmt.Errorf("function only accepts structs; got %s", val.Kind()))
	}

	for i := 0; i < val.NumField(); i++ {
		// Ignore fields that are private
		if !val.Field(i).CanSet() {
			continue
		}

		tags := extractTags(val, i)
		if len(tags) == 0 {
			continue
		}

		m := tagSliceToMap(tags)
		if err := validateTags(m); err != nil {
			errs = append(errs, err)
			continue
		}

		cmp := derefInterface(val.Field(i).Interface())
		if _, ok := m[optionalTag]; ok {
			// Nil is fine
			if cmp == nil {
				continue
			}

			if err := required(cmp, nil); err != nil {
				// Empty is fine
				continue
			}
		}

		if _, ok := m[optionalTag]; !ok && cmp == nil {
			errs = append(errs, newStErr(val.Field(i).String(), "", errors.New("value is nil")))
			continue
		}

		for _, t := range tags {
			fn, ok := v.tags[t.name]
			if !ok {
				// no such tag
				// TODO consider throwing an error here
				continue
			}

			if err := fn(cmp, t.args); err != nil {
				errs = append(errs, newStErr(val.Field(i).String(), t.name, err))
			}
		}
	}

	return errs
}

// SetTag allows to to create a new tag and use it for validation.
// Current tag that has the same name will get over written.
func (v *Vali) SetTag(tag string, TagFunc TagFunc) {
	v.tags[tag] = TagFunc
}
