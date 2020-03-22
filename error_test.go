package vali

import (
	"errors"
	"testing"
)

func TestAggErr_Error(t *testing.T) {
	type fields struct {
		Sl []error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "single error should only be one line",
			fields: fields{
				Sl: []error{errors.New("a")},
			},
			want: "a",
		},
		{
			name: "two errors should be in two seperate lines",
			fields: fields{
				Sl: []error{errors.New("a"), errors.New("b")},
			},
			want: `a
b`,
		},
		{
			name: "three errors should be in three seperate lines",
			fields: fields{
				Sl: []error{errors.New("a"), errors.New("b"), errors.New("c")},
			},
			want: `a
b
c`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &AggErr{}
			for _, f := range tt.fields.Sl {
				e.addErr(f)
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("AggErr.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
