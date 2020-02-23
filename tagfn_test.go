package vali

import (
	"testing"
	"time"
)

func TestOptional(t *testing.T) {
	type mock struct {
		First int `vali:"optional"`
	}
	type mock2 struct {
		First *int `vali:"optional|max=2"`
	}
	type mock3 struct {
		First int `vali:"required|optional"`
	}
	type mock4 struct {
		First int `vali:"optional|min=2"`
	}
	type args struct {
		s interface{}
	}
	one := 1
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test 'optional' using mock, value is empty and optional, should not error",
			args: args{
				s: &mock{
					First: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "test 'optional' using mock2, value is nil and optional, should not error",
			args: args{
				s: &mock2{
					First: nil,
				},
			},
			wantErr: false,
		},
		{
			name: "test 'optional' using mock2, value not nil and optional, should not error",
			args: args{
				s: &mock2{
					First: &one,
				},
			},
			wantErr: false,
		},
		{
			name: "test 'optional' using mock3, value is required and optional, should error",
			args: args{
				s: &mock3{
					First: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "test 'optional' using mock4, value is optional, but is not empty and is more than min, should not error",
			args: args{
				s: &mock4{
					First: 3,
				},
			},
			wantErr: false,
		},
		{
			name: "test 'optional' using mock4, value is optional, but is not empty and below min, should error",
			args: args{
				s: &mock4{
					First: 1,
				},
			},
			wantErr: true,
		},
	}

	v := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := v.Validate(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("Vali.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNeq(t *testing.T) {
	type mock struct {
		First int `vali:"neq=1"`
	}
	type mock2 struct {
		First float64 `vali:"neq=1.0"`
	}
	type mock3 struct {
		First string `vali:"neq=a"`
	}
	type mock4 struct {
		First string `vali:"neq=1"`
	}
	type mock5 struct {
		First uint `vali:"neq=1"`
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
			name: "test 'neq' using mock, value is not equal, should error",
			args: args{
				s: &mock{
					First: 5,
				},
			},
			wantErr: false,
		},
		{
			name: "test 'neq' using mock, value is equal, should not error",
			args: args{
				s: &mock{
					First: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "test 'neq' using mock2, value is not equal, should error",
			args: args{
				s: &mock2{
					First: 5.0,
				},
			},
			wantErr: false,
		},
		{
			name: "test 'neq' using mock2, value is equal, should not error",
			args: args{
				s: &mock2{
					First: 1.0,
				},
			},
			wantErr: true,
		},
		{
			name: "test 'neq' using mock3, value is not equal, should error",
			args: args{
				s: &mock3{
					First: "xxx",
				},
			},
			wantErr: false,
		},
		{
			name: "test 'neq' using mock3, value if equal, should not error",
			args: args{
				s: &mock3{
					First: "a",
				},
			},
			wantErr: true,
		},
		{
			name: "test 'neq' using mock4, value if not equal, should not error",
			args: args{
				s: &mock3{
					First: "2",
				},
			},
			wantErr: false,
		},
		{
			name: "test 'neq' using mock5, value is equal, should error",
			args: args{
				s: &mock5{
					First: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "test 'neq' using mock5, value is not equal, should not error",
			args: args{
				s: &mock5{
					First: 2,
				},
			},
			wantErr: false,
		},
	}

	v := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := v.Validate(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("Vali.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEq(t *testing.T) {
	type mock struct {
		First int `vali:"eq=1"`
	}
	type mock2 struct {
		First float64 `vali:"eq=1.0"`
	}
	type mock3 struct {
		First string `vali:"eq=a"`
	}
	type mock4 struct {
		First string `vali:"eq=1"`
	}
	type mock5 struct {
		First uint `vali:"eq=1"`
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
			name: "test 'eq' using mock, value is not equal, should error",
			args: args{
				s: &mock{
					First: 5,
				},
			},
			wantErr: true,
		},
		{
			name: "test 'eq' using mock, value is equal, should not error",
			args: args{
				s: &mock{
					First: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "test 'eq' using mock2, value is not equal, should error",
			args: args{
				s: &mock2{
					First: 5.0,
				},
			},
			wantErr: true,
		},
		{
			name: "test 'eq' using mock2, value is equal, should not error",
			args: args{
				s: &mock2{
					First: 1.0,
				},
			},
			wantErr: false,
		},
		{
			name: "test 'eq' using mock3, value is not equal, should error",
			args: args{
				s: &mock3{
					First: "xxx",
				},
			},
			wantErr: true,
		},
		{
			name: "test 'eq' using mock3, value is equal, should not error",
			args: args{
				s: &mock3{
					First: "a",
				},
			},
			wantErr: false,
		},
		{
			name: "test 'eq' using mock4, value is equal, should not error",
			args: args{
				s: &mock4{
					First: "1",
				},
			},
			wantErr: false,
		},
		{
			name: "test 'eq' using mock5, value is equal, should not error",
			args: args{
				s: &mock5{
					First: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "test 'eq' using mock5, value is not equal, should error",
			args: args{
				s: &mock5{
					First: 2,
				},
			},
			wantErr: true,
		},
	}
	v := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := v.Validate(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("Vali.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOneOf(t *testing.T) {
	type mock struct {
		First int `vali:"one_of=1,2,3"`
	}
	type mock2 struct {
		First float64 `vali:"one_of=1.0,2.0,3.0"`
	}
	type mock3 struct {
		First string `vali:"one_of=a,b,c"`
	}
	type mock4 struct {
		First string `vali:"one_of=a,4,c"`
	}
	type mock5 struct {
		First uint `vali:"one_of=1,2,3"`
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
			name: "test 'one_of' using mock, value does not have one of, should error",
			args: args{
				s: &mock{
					First: 5,
				},
			},
			wantErr: true,
		},
		{
			name: "test 'one_of' using mock, value does have one of, should not error",
			args: args{
				s: &mock{
					First: 2,
				},
			},
			wantErr: false,
		},
		{
			name: "test 'one_of' using mock2, value does not have one of, should error",
			args: args{
				s: &mock2{
					First: 5.0,
				},
			},
			wantErr: true,
		},
		{
			name: "test 'one_of' using mock2, value does have one of, should not error",
			args: args{
				s: &mock2{
					First: 2.0,
				},
			},
			wantErr: false,
		},
		{
			name: "test 'one_of' using mock3, value does not have one of, should error",
			args: args{
				s: &mock3{
					First: "xxx",
				},
			},
			wantErr: true,
		},
		{
			name: "test 'one_of' using mock3, value does have one of, should not error",
			args: args{
				s: &mock3{
					First: "b",
				},
			},
			wantErr: false,
		},
		{
			name: "test 'one_of' using mock5, value does have one of, should not error",
			args: args{
				s: &mock5{
					First: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "test 'one_of' using mock5, value does have one of, should error",
			args: args{
				s: &mock5{
					First: 5,
				},
			},
			wantErr: true,
		},
		{
			name: "test 'one_of' using mock4, type mismatch, should error",
			args: args{
				s: &mock3{
					First: "x",
				},
			},
			wantErr: true,
		},
	}

	v := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := v.Validate(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("Vali.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMax(t *testing.T) {
	more := 10
	less := 1
	moreF := 10.0
	lessF := 1.0
	type mock struct {
		First  int     `vali:"max=5"`
		Second float64 `vali:"max=2.0"`
		Third  *int    `vali:"max=5"`
	}
	type mock2 struct {
		First  *int
		Second int `vali:"max=*First"`
	}
	type mock3 struct {
		First  time.Time
		Second time.Time `vali:"max=*First"`
	}
	type mock4 struct {
		First string `vali:"max=5"`
	}
	type mock5 struct {
		Fn func() `vali:"max=5"`
	}
	type mock6 struct {
		First *int `vali:"max=5"`
	}
	type mock7 struct {
		First uint `vali:"max=5"`
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
			name: "test 'max' using mock, all values are below the max amount",
			args: args{
				s: &mock{
					First:  less,
					Second: lessF,
					Third:  &less,
				},
			},
			wantErr: false,
		},
		{
			name: "test 'max' using mock, first is more, second and third is less",
			args: args{
				s: &mock{
					First:  more,
					Second: lessF,
					Third:  &less,
				},
			},
			wantErr: true,
		},
		{
			name: "test 'max' using mock, first and second are more",
			args: args{
				s: &mock{
					First:  more,
					Second: moreF,
					Third:  &less,
				},
			},
			wantErr: true,
		},
		{
			name: "test 'max' using mock, first and second are less third is more",
			args: args{
				s: &mock{
					First:  less,
					Second: lessF,
					Third:  &more,
				},
			},
			wantErr: true,
		},
		{
			name: "test 'max' using mock2, first is not validated, second is less than first",
			args: args{
				s: &mock2{
					First:  &more,
					Second: less,
				},
			},
			wantErr: false,
		},
		{
			name: "test 'max' using mock2, first is not validated, second is more than first",
			args: args{
				s: &mock2{
					First:  &less,
					Second: more,
				},
			},
			wantErr: true,
		},
		{
			name: "test 'max' using mock3, first is not validated, second is less than first",
			args: args{
				s: &mock3{
					First:  time.Now(),
					Second: time.Now().Add(-1 * time.Hour),
				},
			},
			wantErr: false,
		},
		{
			name: "test 'max' using mock3, first is not validated, second is more than first",
			args: args{
				s: &mock3{
					First:  time.Now(),
					Second: time.Now().Add(1 * time.Hour),
				},
			},
			wantErr: true,
		},
		{
			name: "test 'max' using mock4, string is above max",
			args: args{
				s: &mock4{
					First: "aaaaaaa",
				},
			},
			wantErr: true,
		},
		{
			name: "test 'max' using mock4, string is below max",
			args: args{
				s: &mock4{
					First: "a",
				},
			},
			wantErr: false,
		},
		{
			name: "test 'max' using mock6,field values are nil, should error",
			args: args{
				s: &mock6{
					First: nil,
				},
			},
			wantErr: true,
		},
		{
			name: "test 'max' using mock2,pointing to a nil struct value, should error",
			args: args{
				s: &mock2{
					First:  nil,
					Second: more,
				},
			},
			wantErr: true,
		},
		{
			name: "test 'max' using mock5, can't test func should error",
			args: args{
				s: &mock5{
					Fn: func() {},
				},
			},
			wantErr: true,
		},
		{
			name: "test 'max' using mock7, comparing uint with tag, is less, should not return an error",
			args: args{
				s: &mock7{
					First: 4,
				},
			},
			wantErr: false,
		},
		{
			name: "test 'max' using mock7, comparing uint with tag, is more, should return an error",
			args: args{
				s: &mock7{
					First: 6,
				},
			},
			wantErr: true,
		},
	}

	v := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := v.Validate(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("Vali.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMin(t *testing.T) {
	more := 10
	moreF := 10.0
	less := 1
	lessF := 1.0
	type mock struct {
		First  int     `vali:"min=5"`
		Second float64 `vali:"min=2.0"`
		Third  *int    `vali:"min=5"`
	}
	type mock2 struct {
		First  *int
		Second int `vali:"min=*First"`
	}
	type mock3 struct {
		First  time.Time
		Second time.Time `vali:"min=*First"`
	}
	type mock4 struct {
		First string `vali:"min=5"`
	}
	type mock5 struct {
		Fn func() `vali:"min=5"`
	}
	type mock6 struct {
		First *int `vali:"min=5"`
	}
	type mock7 struct {
		First uint `vali:"min=5"`
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
			name: "test 'min' using mock, all values are above the required amount",
			args: args{
				s: &mock{
					First:  more,
					Second: moreF,
					Third:  &more,
				},
			},
			wantErr: false,
		},
		{
			name: "test 'min' using mock, first is less second and third is more",
			args: args{
				s: &mock{
					First:  less,
					Second: moreF,
					Third:  &more,
				},
			},
			wantErr: true,
		},
		{
			name: "test 'min' using mock, first and second are less",
			args: args{
				s: &mock{
					First:  less,
					Second: lessF,
					Third:  &more,
				},
			},
			wantErr: true,
		},
		{
			name: "test 'min' using mock, first and second are more third is less",
			args: args{
				s: &mock{
					First:  more,
					Second: moreF,
					Third:  &less,
				},
			},
			wantErr: true,
		},
		{
			name: "test 'min' using mock2, first is not validated, second is more than first",
			args: args{
				s: &mock2{
					First:  &less,
					Second: more,
				},
			},
			wantErr: false,
		},
		{
			name: "test 'min' using mock2, first is not validated, second is less than first",
			args: args{
				s: &mock2{
					First:  &more,
					Second: less,
				},
			},
			wantErr: true,
		},
		{
			name: "test 'min' using mock3, first is not validated, second is more than first",
			args: args{
				s: &mock3{
					First:  time.Now(),
					Second: time.Now().Add(1 * time.Hour),
				},
			},
			wantErr: false,
		},
		{
			name: "test 'min' using mock3, first is not validated, second is less than first",
			args: args{
				s: &mock3{
					First:  time.Now(),
					Second: time.Now().Add(-1 * time.Hour),
				},
			},
			wantErr: true,
		},
		{
			name: "test 'min' using mock4, string is above min",
			args: args{
				s: &mock4{
					First: "aaaaaaa",
				},
			},
			wantErr: false,
		},
		{
			name: "test 'min' using mock4, string is below min",
			args: args{
				s: &mock4{
					First: "a",
				},
			},
			wantErr: true,
		},
		{
			name: "test 'min' using mock6,field values are nil, should error",
			args: args{
				s: &mock6{
					First: nil,
				},
			},
			wantErr: true,
		},
		{
			name: "test 'min' using mock2,pointing to a nil struct value, should error",
			args: args{
				s: &mock2{
					First:  nil,
					Second: more,
				},
			},
			wantErr: true,
		},
		{
			name: "test 'min' using mock5, can't test func should error",
			args: args{
				s: &mock5{
					Fn: func() {},
				},
			},
			wantErr: true,
		},
		{
			name: "test 'min' using mock7, comparing uint with tag, is less, should return an error",
			args: args{
				s: &mock7{
					First: 4,
				},
			},
			wantErr: true,
		},
		{
			name: "test 'min' using mock7, comparing uint with tag, is more, should not return an error",
			args: args{
				s: &mock7{
					First: 6,
				},
			},
			wantErr: false,
		},
	}

	v := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := v.Validate(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("Vali.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRequiredWithout(t *testing.T) {
	str := "str"
	strEmpty := ""
	type mock struct {
		Str  string `vali:"-"`
		Str2 string `vali:"required_without=*Str"`
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
			name: "test 'required_without', both values not empty no error",
			args: args{
				s: &mock{
					Str:  str,
					Str2: str,
				},
			},
			wantErr: false,
		},
		{
			name: "test 'required_without', Str empty Str2 required not empty no error",
			args: args{
				s: &mock{
					Str:  strEmpty,
					Str2: str,
				},
			},
			wantErr: false,
		},
		{
			name: "test 'required_without', Str empty Str2 required empty should error",
			args: args{
				s: &mock{
					Str:  strEmpty,
					Str2: strEmpty,
				},
			},
			wantErr: true,
		},
	}

	v := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := v.Validate(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("Vali.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRequired(t *testing.T) {
	str := "str"
	strEmpty := ""
	type mock struct {
		Str    string  `vali:"required"`
		PtrStr *string `vali:"required"`
	}
	type mockfn struct {
		Fn func() `vali:"required"`
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
			name: "test required, Str empty, PtrStr not empty",
			args: args{
				s: &mock{
					Str:    strEmpty,
					PtrStr: &str,
				},
			},
			wantErr: true,
		},
		{
			name: "test required, Str not empty, PtrStr empty",
			args: args{
				s: &mock{
					Str:    str,
					PtrStr: &strEmpty,
				},
			},
			wantErr: true,
		},
		{
			name: "test required, Str not empty, PtrStr nil",
			args: args{
				s: &mock{
					Str:    str,
					PtrStr: nil,
				},
			},
			wantErr: true,
		},
		{
			name: "test required, Str not empty, PtrStr not empty",
			args: args{
				s: &mock{
					Str:    str,
					PtrStr: &str,
				},
			},
			wantErr: false,
		},
		{
			name: "test required, function not empty, should not error",
			args: args{
				s: &mockfn{
					Fn: func() {},
				},
			},
			wantErr: false,
		},
		{
			name: "test required, function nil, should error",
			args: args{
				s: &mockfn{
					Fn: nil,
				},
			},
			wantErr: true,
		},
	}

	v := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := v.Validate(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("Vali.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
