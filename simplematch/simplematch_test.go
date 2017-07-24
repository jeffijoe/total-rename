package simplematch

import (
	"reflect"
	"testing"
)

func TestNewMatcher(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name string
		args args
		want *Matcher
	}{
		{
			name: "case 1",
			args: args{
				pattern: "a|b|c",
			},
			want: &Matcher{
				splat: []string{
					"A",
					"B",
					"C",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMatcher(tt.args.pattern); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMatcher() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatcher_Matches(t *testing.T) {
	type fields struct {
		pattern string
	}
	type args struct {
		s string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			fields: fields{
				pattern: "a|b|cool",
			},
			args: args{
				s: "i match cause of a",
			},
			want: true,
		},
		{
			fields: fields{
				pattern: "a|b|cool",
			},
			args: args{
				s: "i cool",
			},
			want: true,
		},
		{
			fields: fields{
				pattern: "a|b|cool",
			},
			args: args{
				s: "hohoho",
			},
			want: false,
		},
		{
			fields: fields{
				pattern: "",
			},
			args: args{
				s: "hohoho",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMatcher(tt.fields.pattern)
			if got := m.Matches(tt.args.s); got != tt.want {
				t.Errorf("Matcher.Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}
