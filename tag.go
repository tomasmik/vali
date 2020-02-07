package vali

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type tag struct {
	name string
	args []interface{}
}

func extractTags(mainStruct reflect.Value, fieldIndex int) []tag {
	tgs := make([]tag, 0)
	vtag := mainStruct.Type().Field(fieldIndex).Tag.Get(valiTag)
	// Dont validate fields which have no tags
	if vtag == "" || vtag == "-" {
		return tgs
	}

	for _, t := range strings.Split(vtag, tagSep) {
		parts := strings.Split(t, equalsSep)

		tg := tag{
			name: parts[0],
			args: make([]interface{}, 0),
		}

		if !strings.Contains(t, equalsSep) || len(parts) == 1 {
			tgs = append(tgs, tg)
			continue
		}

		for _, f := range strings.Split(parts[1], valueSep) {
			if !strings.HasPrefix(f, pointerToField) {
				in, err := strconv.ParseInt(f, 10, 64)
				if err == nil {
					tg.args = append(tg.args, in)
					continue
				}
				fl, err := strconv.ParseFloat(f, 64)
				if err == nil {
					tg.args = append(tg.args, fl)
					continue
				}

				tg.args = append(tg.args, derefInterface(f))
				continue
			}

			f = strings.TrimPrefix(f, pointerToField)
			for j := 0; j < mainStruct.NumField(); j++ {
				// This used to check if the struct field pointer is pointing to itself
				// I am not sure if we should allow this or not. If the user wants to
				// he can point to himself I guess...
				// if mainStruct.Type().Field(currentField).Name == f {
				// 	continue
				// }

				if mainStruct.Type().Field(j).Name == f {
					in, _ := getInterface(mainStruct.Field(j))
					tg.args = append(tg.args, in)
				}
			}
		}
		tgs = append(tgs, tg)
	}

	return tgs
}

func validateTags(m map[string]struct{}) error {
	count := 0
	for k := range m {
		switch k {
		case optionalTag, requiredWithoutTag, requiredTag:
			count++
		}
	}

	if count > 1 {
		return fmt.Errorf("a field can only have one of: %s, %s, %s", optionalTag, requiredTag, requiredWithoutTag)
	}

	return nil
}

func tagSliceToMap(tgsl []tag) map[string]struct{} {
	m := map[string]struct{}{}
	for _, f := range tgsl {
		m[f.name] = struct{}{}
	}
	return m
}
