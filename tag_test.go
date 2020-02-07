package vali

import (
	"reflect"
	"testing"
)

func Test_extractTags(t *testing.T) {
	mock := struct {
		First int `json:"f" vali:"required"`
	}{
		First: 1,
	}
	mock2 := struct {
		First int `json:"f" vali:"min=2"`
	}{
		First: 1,
	}
	mock3 := struct {
		First int `json:"f" vali:"optional|min=2,3"`
	}{
		First: 1,
	}

	type args struct {
		mainStruct reflect.Value
		fieldIndex int
	}
	tests := []struct {
		name string
		args args
		want []tag
	}{
		{
			name: "should get the required tag",
			args: args{
				mainStruct: reflect.ValueOf(mock),
				fieldIndex: 0,
			},
			want: []tag{
				tag{name: requiredTag, args: []interface{}{}},
			},
		},
		{
			name: "should get the min tag and a single arg 2",
			args: args{
				mainStruct: reflect.ValueOf(mock2),
				fieldIndex: 0,
			},
			want: []tag{
				tag{name: minTag, args: []interface{}{int64(2)}},
			},
		},
		{
			name: "should get two tags: optional with no args and min with args 2,3",
			args: args{
				mainStruct: reflect.ValueOf(mock3),
				fieldIndex: 0,
			},
			want: []tag{
				tag{name: optionalTag, args: []interface{}{}},
				tag{name: minTag, args: []interface{}{int64(2), int64(3)}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractTags(tt.args.mainStruct, tt.args.fieldIndex); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateTags(t *testing.T) {
	type args struct {
		m map[string]struct{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "map only has 'optional', should not error",
			args: args{
				m: map[string]struct{}{
					optionalTag: struct{}{},
				},
			},
			wantErr: false,
		},
		{
			name: "map only has 'required', should not error",
			args: args{
				m: map[string]struct{}{
					requiredTag: struct{}{},
				},
			},
			wantErr: false,
		},
		{
			name: "map only has 'required_without', should not error",
			args: args{
				m: map[string]struct{}{
					requiredWithoutTag: struct{}{},
				},
			},
			wantErr: false,
		},
		{
			name: "map has optional and required, should error",
			args: args{
				m: map[string]struct{}{
					optionalTag: struct{}{},
					requiredTag: struct{}{},
				},
			},
			wantErr: true,
		},
		{
			name: "map has optional and required_without, should error",
			args: args{
				m: map[string]struct{}{
					optionalTag:        struct{}{},
					requiredWithoutTag: struct{}{},
				},
			},
			wantErr: true,
		},
		{
			name: "map has required and required_without, should error",
			args: args{
				m: map[string]struct{}{
					requiredTag:        struct{}{},
					requiredWithoutTag: struct{}{},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateTags(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("validateTags() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_tagSliceToMap(t *testing.T) {
	type args struct {
		tgsl []tag
	}
	tests := []struct {
		name string
		args args
		want map[string]struct{}
	}{
		{
			name: "map has required and required_without, should error",
			args: args{
				tgsl: []tag{
					tag{name: optionalTag, args: []interface{}{}},
					tag{name: minTag, args: []interface{}{int64(2), int64(3)}},
				},
			},
			want: map[string]struct{}{
				optionalTag: struct{}{},
				minTag:      struct{}{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tagSliceToMap(tt.args.tgsl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tagSliceToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
