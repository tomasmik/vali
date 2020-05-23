package vali

import (
	"errors"
	"fmt"
	"reflect"
)

const (
	// valiTag is the default string that is used
	// to tag interfaces values for validation.
	valiTag = "vali"
	// equals is a sign used to add values to a tag.
	equalsSep = "="
	// valueSep is used to seperate between values for a tag
	valueSep = ","
	// tagSep is used to seperate tags from each other
	tagSep = "|"
	// pointerToField to a field is used when you want to compare
	// one field to another.
	// Example:
	/*
	 type mock struct {
	 Str  string `vali:"-"`
	 Str2 string `vali:"required_without=*Str"`
	 }
	*/
	// The *Str points to another struct value, so there values
	// will be compared.
	pointerToField = "*"
	// dive in to a slice, validating the contents
	// Example:
	/*
	 type mock struct {
	 Str2 []string `vali:"min=2|>|one_of=a,b,c"`
	 }
	*/
	dive = ">"
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
	tags   tags
	types  types
	tgName string
}

// ErrSkipFurther is an error that can be used as a return value
// when validating struct field. If it's returned further validation
// is skipped
var ErrSkipFurther = errors.New("skip further")

// New returns a new validator instance,
// with the default predefined types.
func New() *Vali {
	return &Vali{
		tgName: valiTag,
		types:  map[reflect.Type]TypeFunc{},
		tags: map[string]TagFunc{
			requiredTag:        required,
			requiredWithoutTag: required_without,
			maxTag:             max,
			minTag:             min,
			oneofTag:           oneof,
			noneofTag:          noneof,
			eqTag:              eq,
			neqTag:             neq,
			dupsTag:            dups,
			optionalTag:        optional,
		},
	}
}

// NewEmpty returns a new vali validator without
// any predefined tags allowing the user to configure whatever he needs.
func NewEmpty() *Vali {
	return &Vali{
		tgName: valiTag,
		types:  map[reflect.Type]TypeFunc{},
		tags:   map[string]TagFunc{},
	}
}

// Validate accepts a struct and validates its according to the given tags.
// Validations are applied in this order:
// 1. Type validation if one is set.
// 2. Tag validation in order the tags were set.
//
// The return value `error` can be type asserted in to `*vali.AggErr`
// which allows to explore each error seprately.
// Example:
/*

        // In the real world you would define the vali.New() before calling
        // the `Validate()` method.
	err := vali.New().Validate(str)

*/
func (v *Vali) Validate(s interface{}) error {
	errs := newAggErr()

	if s == nil {
		return errs.addErr(errors.New("struct is nil"))
	}

	val := reflect.ValueOf(s)
	if val.Kind() != reflect.Ptr {
		return errs.addErr(fmt.Errorf("function only accepts pointer to structs; got %s", val.Kind()))
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
		return errs.addErr(fmt.Errorf("function only accepts structs; got %s", val.Kind()))
	}

	if fn, ok := v.types[val.Type()]; ok {
		if err := fn(orgVal.Interface()); err != nil {
			errs.addErr(err)
		}
	}

	for i := 0; i < val.NumField(); i++ {
		if !val.Field(i).CanSet() {
			continue
		}

		tags := extractTags(val, v.tgName, i)
		if len(tags) == 0 {
			continue
		}

		m := tagSliceToMap(tags)
		if err := validateTags(m); err != nil {
			errs.addErr(err)
			continue
		}

		// This is sort of hacky way to allow us to easily
		// convert between a "dive" validation and a single field
		// validation.
		// Maybe it should be improved in the future
		cmp := []interface{}{
			DerefInterface(val.Field(i).Interface()),
		}

		if derf, ok := derefReflectValue(val.Field(i)); ok {
			if derf.Kind() == reflect.Struct {
				ss := val.Field(i).Interface()
				if ers := v.Validate(&ss); ers != nil {
					errs.addErr(ers)
				}
			}
		}

		if err := v.validateField(val.Type().Field(i).Name, cmp, tags); err != nil {
			var b *bubbleErr
			var e *tagError

			if errors.As(err, &b) {
				return b.err
			} else if errors.As(err, &e) {
				errs.addErr(e)
			} else {
				return err
			}
		}
	}

	return errs.toError()
}

// SetTagValidation allows to create a new tag and use it for validation.
// Current tag that has the same name will get over written.
// You can return custom errors from custom tags by returning a BubbleErr.
// Example:
/*

	v.SetTagValidation("mockTg", func(s interface{}, o []interface{}) error {
		return nil
	})

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
// You can return custom errors from custom tags by returning a BubbleErr.
// Example:
/*

	vali.New().SetTypeValidation(&CustomMock{}, func(s interface{}) error {
		return nil
	})

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

// RenameTag can be used to change the default `valiTag`
// string to your own for marking struct fields for validation.
//
// Empty strings as `t` values are ignored
func (v *Vali) RenameTag(t string) {
	if t == "" {
		return
	}

	v.tgName = t
}

// validateField is a helper method which holds the validation code for a specific
// field. It calls itself recursively if it finds a dive tag validating
// the inside of a given `slice` or `array`.
func (v *Vali) validateField(field string, cmp []interface{}, tags []tag) error {
	for _, c := range cmp {
		for i, t := range tags {
			if t.name == dive {
				cmp, err := rebuildCmpSlice(cmp[0])
				if err != nil {
					return newTagError(field, t.name, err)
				}
				return v.validateField(field, cmp, tags[i+1:])
			}
			fn, ok := v.tags[t.name]
			if !ok {
				// no such tag
				// TODO consider throwing an error here
				continue
			}

			if err := fn(c, t.args); err != nil {
				if errors.Is(err, ErrSkipFurther) {
					return nil
				}

				var b *bubbleErr
				if errors.As(err, &b) {
					return b
				}
				return newTagError(field, t.name, err)
			}
		}
	}
	return nil
}

// rebuildCmpSlice is a helper function to rebuild the comparison
// slice if possible
func rebuildCmpSlice(val interface{}) ([]interface{}, error) {
	derf := interfaceToReflectVal(val)
	switch derf.Kind() {
	case reflect.Array, reflect.Slice:
	default:
		return nil, errors.New("value is not a slice, can't use it")
	}

	newcmp := []interface{}{}
	for j := 0; j < derf.Len(); j++ {
		newcmp = append(
			newcmp,
			DerefInterface(derf.Index(j).Interface()))
	}
	return newcmp, nil
}
