package vali

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
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
type tags map[string]tagFunc

// tagFunc is a func that is used to validate a field `s`
// when a tag is found.
type tagFunc func(s interface{}, o interface{}) error

// Vali is a struct that holds all the configuration
// for the validation tool.
//
// It can be extended using its public methods by adding
// additional tag validation funcs or type validation funcs.
// Just make sure you add them before actually using the `Validate` method.
type Vali struct {
	tags tags
	cfg  *Config
}

// Config is a config struct that can be given to the
// `New()` func when creating the validator.
type Config struct {
	IgnoreNilPointer bool
}

// New returns a new validator instance.
func New(cfg *Config) *Vali {
	if cfg == nil {
		cfg = &Config{
			IgnoreNilPointer: false,
		}
	}

	return &Vali{
		tags: defaultTags(),
		cfg:  cfg,
	}
}

// Validate accepts a struct and validates its according to the given tags.
func (v *Vali) Validate(s interface{}) error {
	if s == nil {
		return errors.New("struct is nil")
	}

	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// we only accept structs
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("function only accepts structs; got %s", val.Kind())
	}
	for i := 0; i < val.NumField(); i++ {
		// Ignore fields that are private
		if !val.Field(i).CanSet() {
			continue
		}

		tag := val.Type().Field(i).Tag.Get(valiTag)
		// Dont validate fields which have no tags
		if tag == "" || tag == "-" {
			continue
		}

		valis, err := v.extract(tag)
		if err != nil {
			return newTagErr(tag, err)
		}

		for tag, vali := range valis {
			var fields interface{}
			if strings.Contains(tag, equalsSep) {
				pastEqParts := strings.Split(tag, equalsSep)
				if len(pastEqParts) != 2 {
					return newTagErr(tag, fmt.Errorf("only 1 equalSep `%s` sign is allowed", equalsSep))
				}

				pasteq := pastEqParts[1]
				parts := strings.Split(pasteq, valueSep)
				if strings.Contains(tag, pointerToField) {
					allPointers := true
					for _, f := range parts {
						if !strings.Contains(f, pointerToField) {
							allPointers = false
						}
					}
					if !allPointers {
						return newTagErr(tag, errors.New("all values must be pointers if one is a pointer"))

					}

					var structFields []reflect.Value
					for _, f := range parts {
						f = strings.TrimPrefix(f, pointerToField)
						for j := 0; j < val.NumField(); j++ {
							if val.Type().Field(i).Name == f {
								return newTagErr(tag, errors.New("cant point to yourself"))
							}

							if val.Type().Field(j).Name == f {
								structFields = append(structFields, val.Field(i))
							}
						}
					}
					fields = structFields
				} else {
					fields = parts
				}
			}

			if err := vali(val.Field(i), fields); err != nil {
				return newStErr(val.Field(i).String(), tag, err)
			}
		}

	}

	return nil
}

func (v *Vali) extract(tags string) (tags, error) {
	fns := make(map[string]tagFunc, 0)
	for _, tag := range strings.Split(tags, tagSep) {

		parts := strings.Split(tag, equalsSep)
		fn, ok := v.tags[parts[0]]
		if !ok {
			return nil, fmt.Errorf("%s tag does not exist", tag)
		}

		fns[tag] = fn
	}

	return fns, nil
}
