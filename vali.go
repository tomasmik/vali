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

// types map stores all the types that we can validate.
// These types can be set by the package user.
type types map[reflect.Type]TypeFunc

// TagFunc is a func that is used to validate a field `s`
// using data providing in slice `o`.
//
// `s` can be nil if `Config.IgnoreNilPointer` is set to true
// Slice `o` can be empty or nil and doesn't have to be used if not needed.
type TagFunc func(s interface{}, o []interface{}) error

// Type func is a func that is can be set and used
// to validate the given type `s`
type TypeFunc func(s interface{}) error

// Vali is a struct that holds all the configuration
// for the validation tool.
//
// It can be extended using its public methods by adding
// additional tag validation funcs or type validation funcs.
// Just make sure you add them before actually using the `Validate` method
// as no thread safety exists, so using and editting validation funcs
// will result in a race condition.
type Vali struct {
	tags  tags
	types types
}

// New returns a new validator instance.
func New() *Vali {
	return &Vali{
		types: map[reflect.Type]TypeFunc{},
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
// Validations are applied in this order:
// 1. Type validation if one is set.
// 2. Tag validation in order the tags were set.
// It returns a slice of errors, if the errors consists of any errors
// that are not of type *StErr consider the validation a failure
// and fix the errors before validating a struct again.
// Example:
/*
`
        // In the real world you would define the vali.New() before calling
        // the `Validate()` method.
	err := vali.New().Validate(str)
`
*/
func (v *Vali) Validate(s interface{}) []error {
	var errs []error
	if s == nil {
		return append(errs, errors.New("struct is nil"))
	}

	val := reflect.ValueOf(s)
	if val.Kind() != reflect.Ptr {
		return append(errs, fmt.Errorf("function only accepts pointer to structs; got %s", val.Kind()))
	}

	orgVal := val
	for val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		if val.IsNil() {
			break
		}
		val = val.Elem()
	}

	// we only accept structs
	if val.Kind() != reflect.Struct {
		return append(errs, fmt.Errorf("function only accepts structs; got %s", val.Kind()))
	}

	if fn, ok := v.types[val.Type()]; ok {
		if err := fn(orgVal.Interface()); err != nil {
			errs = append(errs, err)
		}
	}

	for i := 0; i < val.NumField(); i++ {
		if !val.Field(i).CanSet() {
			continue
		}

		cmp := derefInterface(val.Field(i).Interface())
		if reflect.ValueOf(cmp).Kind() == reflect.Struct {
			ss := val.Field(i).Interface()
			if ers := v.Validate(&ss); ers != nil {
				errs = append(errs, ers...)
			}
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

		for _, t := range tags {
			fn, ok := v.tags[t.name]
			if !ok {
				// no such tag
				// TODO consider throwing an error here
				continue
			}

			if err := fn(cmp, t.args); err != nil {
				errs = append(errs, newStErr(val.Type().Field(i).Name, t.name, err))
			}
		}
	}

	return errs
}

// SetTagValidation allows to create a new tag and use it for validation.
// Current tag that has the same name will get over written.
// Example:
/*
`
	v.SetTagValidation("mockTg", func(s interface{}, o []interface{}) error {
		return nil
	})
`
*/
func (v *Vali) SetTagValidation(tag string, fn TagFunc) {
	if fn == nil || tag == "" {
		return
	}

	v.tags[tag] = fn
}

// SetTypeValidation allows to create new validation funcs for types.
// *T will get rendered down to T, so *T and T will have the same
// validation type func set.
// Setting validation func will override any previous type validation funcs.
// Example:
/*
`
	vali.New().SetTypeValidation(&CustomMock{}, func(s interface{}) error {
		return nil
	})
`
*/
func (v *Vali) SetTypeValidation(typ interface{}, fn TypeFunc) {
	if typ == nil || fn == nil {
		return
	}

	val := reflect.ValueOf(typ)
	for val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		if val.IsNil() {
			return
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return
	}

	v.types[val.Type()] = fn
}
