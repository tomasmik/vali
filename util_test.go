package vali

import (
	"reflect"
	"testing"
)

func TestGetInt(t *testing.T) {
	i := int(1)
	i8 := int8(1)
	i16 := int16(1)
	i32 := int32(1)
	i64 := int64(1)

	i0 := int64(0)
	type args struct {
		s interface{}
	}
	tests := []struct {
		name  string
		args  args
		want  int64
		want1 bool
	}{
		{
			name: "input is an int, should return int64 and true",
			args: args{
				s: i,
			},
			want:  i64,
			want1: true,
		},
		{
			name: "input is an int8, should return int64 and true",
			args: args{
				s: i8,
			},
			want:  i64,
			want1: true,
		},
		{
			name: "input is an int16, should return int64 and true",
			args: args{
				s: i16,
			},
			want:  i64,
			want1: true,
		},
		{
			name: "input is an int32, should return int64 and true",
			args: args{
				s: i32,
			},
			want:  i64,
			want1: true,
		},
		{
			name: "input is an int64, should return int64 and true",
			args: args{
				s: i64,
			},
			want:  i64,
			want1: true,
		},
		{
			name: "input is a string, should return 0 and false",
			args: args{
				s: "asdf",
			},
			want:  i0,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetInt(tt.args.s)
			if got != tt.want {
				t.Errorf("GetInt() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetInt() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetUInt(t *testing.T) {
	i := uint(1)
	i8 := uint8(1)
	i16 := uint16(1)
	i32 := uint32(1)
	i64 := uint64(1)

	i0 := uint64(0)
	type args struct {
		s interface{}
	}
	tests := []struct {
		name  string
		args  args
		want  uint64
		want1 bool
	}{
		{
			name: "input is an uint, should return uint64 and true",
			args: args{
				s: i,
			},
			want:  i64,
			want1: true,
		},
		{
			name: "input is an uint8, should return uint64 and true",
			args: args{
				s: i8,
			},
			want:  i64,
			want1: true,
		},
		{
			name: "input is an uint16, should return uint64 and true",
			args: args{
				s: i16,
			},
			want:  i64,
			want1: true,
		},
		{
			name: "input is an uint32, should return uint64 and true",
			args: args{
				s: i32,
			},
			want:  i64,
			want1: true,
		},
		{
			name: "input is an uint64, should return uint64 and true",
			args: args{
				s: i64,
			},
			want:  i64,
			want1: true,
		},
		{
			name: "input is a string, should return 0 and false",
			args: args{
				s: "asdf",
			},
			want:  i0,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetUInt(tt.args.s)
			if got != tt.want {
				t.Errorf("GetUInt() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetUInt() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetUIntFallback(t *testing.T) {
	ui := uint(1)
	ui8 := uint8(1)
	ui16 := uint16(1)
	ui32 := uint32(1)
	ui64 := uint64(1)
	ui0 := uint64(0)

	i := int(1)
	i8 := int8(1)
	i16 := int16(1)
	i32 := int32(1)
	i64 := int64(1)
	type args struct {
		s interface{}
	}
	tests := []struct {
		name  string
		args  args
		want  uint64
		want1 bool
	}{
		{
			name: "input is an uint, should return uint64 and true",
			args: args{
				s: ui,
			},
			want:  ui64,
			want1: true,
		},
		{
			name: "input is an uint8, should return uint64 and true",
			args: args{
				s: ui8,
			},
			want:  ui64,
			want1: true,
		},
		{
			name: "input is an uint16, should return uint64 and true",
			args: args{
				s: ui16,
			},
			want:  ui64,
			want1: true,
		},
		{
			name: "input is an uint32, should return uint64 and true",
			args: args{
				s: ui32,
			},
			want:  ui64,
			want1: true,
		},
		{
			name: "input is an uint64, should return uint64 and true",
			args: args{
				s: ui64,
			},
			want:  ui64,
			want1: true,
		},
		{
			name: "input is a string, should return 0 and false",
			args: args{
				s: "asdf",
			},
			want:  ui0,
			want1: false,
		},
		{
			name: "input is an int, should return int64 and true",
			args: args{
				s: i,
			},
			want:  ui64,
			want1: true,
		},
		{
			name: "input is an int8, should return int64 and true",
			args: args{
				s: i8,
			},
			want:  ui64,
			want1: true,
		},
		{
			name: "input is an int16, should return int64 and true",
			args: args{
				s: i16,
			},
			want:  ui64,
			want1: true,
		},
		{
			name: "input is an int32, should return int64 and true",
			args: args{
				s: i32,
			},
			want:  ui64,
			want1: true,
		},
		{
			name: "input is an int64, should return int64 and true",
			args: args{
				s: i64,
			},
			want:  ui64,
			want1: true,
		},
		{
			name: "input is a string, should return 0 and false",
			args: args{
				s: "asdf",
			},
			want:  ui0,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetUIntFallback(tt.args.s)
			if got != tt.want {
				t.Errorf("GetUIntFallback() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetUIntFallback() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetFloat64(t *testing.T) {
	f32 := float32(1)
	f64 := float64(1)
	f0 := float64(0)
	type args struct {
		s interface{}
	}
	tests := []struct {
		name  string
		args  args
		want  float64
		want1 bool
	}{
		{
			name: "input is an float32, should return float64 and true",
			args: args{
				s: f32,
			},
			want:  f64,
			want1: true,
		},
		{
			name: "input is an float64, should return float64 and true",
			args: args{
				s: f64,
			},
			want:  f64,
			want1: true,
		},
		{
			name: "input is a string, should return 0 and false",
			args: args{
				s: "asdf",
			},
			want:  f0,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetFloat(tt.args.s)
			if got != tt.want {
				t.Errorf("GetFloat() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetFloat() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDerefInterface(t *testing.T) {
	i := "mock"
	type args struct {
		s interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "value is a pointer, should deref it",
			args: args{
				s: &i,
			},
		},
		{
			name: "value is not a pointer, should not panic",
			args: args{
				s: i,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := interfaceToReflectVal(DerefInterface(tt.args.s))
			if r.Kind() == reflect.Ptr {
				t.Error("DerefInterface() got = not pointer, want pointer")
			}
		})
	}
}
