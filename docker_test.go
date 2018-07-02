package main

import (
	"reflect"
	"testing"
)

func Test_listToMap(t *testing.T) {
	type args struct {
		env []string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "creates a map from a list of env vars",
			args: args{
				env: []string{"FOO=BAR", "VERSION=1", "SOMEKEY=SOMEVAL"},
			},
			want: map[string]string{"FOO": "BAR", "VERSION": "1", "SOMEKEY": "SOMEVAL"},
		},
		{
			name: "handles empty env vars",
			args: args{
				env: []string{"FOO=BAR", "EMPTY"},
			},
			want: map[string]string{"FOO": "BAR", "EMPTY": ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := listToMap(tt.args.env); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("listToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filter(t *testing.T) {
	type args struct {
		env       map[string]string
		whitelist []string
		blacklist []string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "filters and masks a map of variables",
			args: args{
				env: map[string]string{
					"FOO":  "BAR",
					"KEY":  "VALUE",
					"MORE": "STUFF",
				},
				whitelist: []string{"KEY"},
				blacklist: []string{"FOO"},
			},
			want: map[string]string{
				"KEY":  "VALUE",
				"MORE": "*****",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filter(tt.args.env, tt.args.whitelist, tt.args.blacklist); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filter() = %v, want %v", got, tt.want)
			}
		})
	}
}
