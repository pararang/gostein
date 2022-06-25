package stein

import (
	"testing"
)

func Test_removePrefix(t *testing.T) {
	type args struct {
		s      string
		prefix string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should remove prefix",
			args: args{
				s:      "prefix_string",
				prefix: "prefix_",
			},
			want: "string",
		},
		{
			name: "should return s if it doesn't have prefix",
			args: args{
				s:      "string",
				prefix: "prefix_",
			},
			want: "string",
		},
		{
			name: "should return s if prefix is empty",
			args: args{
				s:      "string",
				prefix: "",
			},
			want: "string",
		},
		{
			name: "should remove slash prefix",
			args: args{
				s:      "/prefix_string",
				prefix: "/",
			},
			want: "prefix_string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removePrefix(tt.args.s, tt.args.prefix); got != tt.want {
				t.Errorf("removePrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_removeSuffix(t *testing.T) {
	type args struct {
		s      string
		suffix string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should remove suffix",
			args: args{
				s:      "string_suffix",
				suffix: "_suffix",
			},
			want: "string",
		},
		{
			name: "should return s if it doesn't have suffix",
			args: args{
				s:      "string",
				suffix: "_suffix",
			},
			want: "string",
		},
		{
			name: "should return s if suffix is empty",
			args: args{
				s:      "string",
				suffix: "",
			},
			want: "string",
		},
		{
			name: "should remove slash suffix",
			args: args{
				s:      "string/",
				suffix: "/",
			},
			want: "string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeSuffix(tt.args.s, tt.args.suffix); got != tt.want {
				t.Errorf("removeSuffix() = %v, want %v", got, tt.want)
			}
		})
	}
}

