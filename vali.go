package vali

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
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

		valis, err := v.extractTags(tag)
		if err != nil {
			return newTagErr(tag, err)
		}

		for tag, vali := range valis {
			with, err := v.extractValues(val, i, tag)
			if err != nil {
				return newTagErr(tag, err)
			}

			cmp := derefInterface(val.Field(i).Interface())
			if cmp == nil && !v.cfg.IgnoreNilPointer {
				return newStErr(val.Field(i).String(), tag, errors.New("value is nil"))
			}

			if err := vali(cmp, with); err != nil {
				return newStErr(val.Field(i).String(), tag, err)
			}
		}
	}

	return nil
}

// SetTag allows to to create a new tag and use it for validation.
// Current tag that has the same name will get over written.
func (v *Vali) SetTag(tag string, TagFunc TagFunc) {
	v.tags[tag] = TagFunc
}

func (v *Vali) extractTags(tags string) (tags, error) {
	fns := make(map[string]TagFunc, 0)
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

func (v *Vali) extractValues(val reflect.Value, currentField int, tag string) ([]interface{}, error) {
	var fields []interface{}
	if !strings.Contains(tag, equalsSep) {
		return fields, nil
	}
	pastEqParts := strings.Split(tag, equalsSep)
	if len(pastEqParts) != 2 {
		return nil, fmt.Errorf("only 1 equalSep `%s` sign is allowed", equalsSep)
	}

	pasteq := pastEqParts[1]
	for _, f := range strings.Split(pasteq, valueSep) {
		if !strings.HasPrefix(f, pointerToField) {
			in, err := strconv.ParseInt(f, 10, 64)
			if err == nil {
				fields = append(fields, in)
				continue
			}
			fl, err := strconv.ParseFloat(f, 64)
			if err == nil {
				fields = append(fields, fl)
				continue
			}

			fields = append(fields, derefInterface(f))
			continue
		}

		f = strings.TrimPrefix(f, pointerToField)
		for j := 0; j < val.NumField(); j++ {
			if val.Type().Field(currentField).Name == f {
				return nil, errors.New("cant point to yourself")
			}

			if val.Type().Field(j).Name == f {
				in, _ := getInterface(val.Field(j))
				fields = append(fields, in)
			}
		}
	}

	return fields, nil
}
