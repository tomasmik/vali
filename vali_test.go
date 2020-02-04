package vali

import (
	"testing"
)

func TestMoreThan(t *testing.T) {
	more := 10
	less := 1
	type mock struct {
		First  int  `vali:"more_than=5"`
		Second int  `vali:"more_than=2"`
		Third  *int `vali:"more_than=5"`
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
			name: "test 'more_than' all values are above the required amount",
			args: args{
				s: &mock{
					First:  more,
					Second: more,
					Third:  &more,
				},
			},
			wantErr: false,
		},
		{
			name: "test 'more_than' first is less second and third is more",
			args: args{
				s: &mock{
					First:  less,
					Second: more,
					Third:  &more,
				},
			},
			wantErr: true,
		},
		{
			name: "test 'more_than' first and second are less",
			args: args{
				s: &mock{
					First:  less,
					Second: less,
					Third:  &more,
				},
			},
			wantErr: true,
		},
		{
			name: "test 'more_than' first and second are more third is less",
			args: args{
				s: &mock{
					First:  more,
					Second: more,
					Third:  &less,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := New(nil).Validate(tt.args.s); (err != nil) != tt.wantErr {
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := New(nil).Validate(tt.args.s); (err != nil) != tt.wantErr {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := New(nil).Validate(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("Vali.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
