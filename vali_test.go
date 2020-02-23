package vali

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestValidate(t *testing.T) {
	type Mock struct {
		First  string `vali:"eq=a"`
		Second int    `vali:"max=5"`
	}
	type mock2 struct {
		M *Mock `vali:"required"`
	}
	type mock3 struct {
		M *Mock `vali:"required"`
	}
	var mockFn func()
	var mockin interface{}
	ptToPt := &mockin
	type args struct {
		s interface{}
	}

	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "struct inside of a struct is valid, should not error",
			args: args{
				s: &mock2{
					M: &Mock{
						First:  "a",
						Second: 5,
					},
				},
			},
			want: nil,
		},
		{
			name: "struct inside of a struct is not valid, should error",
			args: args{
				s: &mock2{
					M: &Mock{
						First:  "b",
						Second: 5,
					},
				},
			},
			want: newAggErr().addErr(newAggErr().addErr(tagError("First", eqTag, errors.New("b is not equal to a")))),
		},
		{
			name: "struct is nil, should error",
			args: args{
				s: nil,
			},
			want: newAggErr().addErr(errors.New("struct is nil")),
		},
		{
			name: "argument is a func not a struct, should error",
			args: args{
				s: &mockFn,
			},
			want: newAggErr().addErr(fmt.Errorf("function only accepts structs; got %s", reflect.ValueOf(mockFn).Kind())),
		},
		{
			name: "argument is a func not a struct, should error",
			args: args{
				s: &ptToPt,
			},
			want: newAggErr().addErr(fmt.Errorf("function only accepts structs; got %s", reflect.Interface)),
		},
		{
			name: "only pointers to a struct are accepted, should error",
			args: args{
				s: mock2{
					M: &Mock{
						First:  "a",
						Second: 5,
					},
				},
			},
			want: newAggErr().addErr(fmt.Errorf("function only accepts pointer to structs; got %s", reflect.ValueOf(mock2{}).Kind())),
		},
	}

	v := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := v.Validate(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Vali.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetTagValidation(t *testing.T) {
	type mock3 struct {
		First string `vali:"eq_a"`
	}
	type args struct {
		s interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test 'neq' using mock3, value if equal, should not error",
			args: args{
				s: &mock3{
					First: "a",
				},
			},
			wantErr: false,
		},
		{
			name: "test 'neq' using mock4, value if not equal, should not error",
			args: args{
				s: &mock3{
					First: "2",
				},
			},
			wantErr: true,
		},
	}

	v := New()
	t.Run("tag validation should have defined tags and funcs", func(t *testing.T) {
		currTags := len(v.tags)
		v.SetTagValidation("", func(s interface{}, o []interface{}) error {
			return nil
		})
		v.SetTagValidation("a", nil)
		v.SetTagValidation("", nil)

		if len(v.tags) != currTags {
			t.Errorf("expected to find %d custom tags, found: %d", currTags, len(v.tags))
		}
	})

	t.Run("tag validation add one", func(t *testing.T) {
		v.SetTagValidation("eq_a", func(s interface{}, o []interface{}) error {
			str, ok := s.(string)
			if !ok {
				return errors.New("not a string")
			}
			if str != "a" {
				return errors.New("not eq to 'a'")
			}

			return nil
		})

		if _, ok := v.tags["eq_a"]; !ok {
			t.Error("expected to find tag `eq_a`")
		}
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := v.Validate(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("Vali.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValiSetTypeValidation(t *testing.T) {
	type CustomMock struct {
		First int `vali:"min=4"`
	}
	type fields struct {
		tags  tags
		types types
	}
	type args struct {
		s interface{}
	}

	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "struct is valid, should not error",
			args: args{
				s: &CustomMock{
					First: 5,
				},
			},
			want: nil,
		},
		{
			name: "struct is valid, should not error",
			args: args{
				s: &CustomMock{
					First: 2,
				},
			},
			want: newAggErr().addErr(tagError("First", minTag, fmt.Errorf("%d is less than %d", 2, 4))),
		},
		{
			name: "struct is valid, should not error",
			args: args{
				s: &CustomMock{
					First: 3,
				},
			},
			want: newAggErr().addErr(
				errors.New("m.First can't be 3"),
				tagError("First", minTag, fmt.Errorf("%d is less than %d", 3, 4))),
		},
	}

	v := New()
	t.Run("type validation only allows structs and funcs that are not nil", func(t *testing.T) {
		v.SetTypeValidation(nil, nil)
		v.SetTypeValidation(CustomMock{}, nil)
		v.SetTypeValidation(func() {}, nil)

		if len(v.types) != 0 {
			t.Errorf("only one type test should be registered")
		}
	})

	t.Run("type validation add one", func(t *testing.T) {
		v.SetTypeValidation(CustomMock{}, func(s interface{}) error {
			return nil
		})
		v.SetTypeValidation(&CustomMock{}, func(s interface{}) error {
			m, ok := s.(*CustomMock)
			if !ok {
				return errors.New("bad type")
			}
			if m.First == 3 {
				return errors.New("m.First can't be 3")
			}
			return nil
		})

		if len(v.types) != 1 {
			t.Errorf("one type test should be registered")
		}
		if _, ok := v.types[reflect.ValueOf(CustomMock{}).Type()]; !ok {
			t.Errorf("should have CustomMock type in the types map")
		}
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := v.Validate(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Vali.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
