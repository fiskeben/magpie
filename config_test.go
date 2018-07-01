package main

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func Test_parseConfig(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *ConfigMap
		wantErr bool
	}{
		{
			name: "parses some lines",
			args: args{
				r: strings.NewReader(`foo=bar
somekey=value
anotherkey=alsovalue`),
			},
			want: func() *ConfigMap {
				m := ConfigMap(map[string]string{"foo": "bar", "somekey": "value", "anotherkey": "alsovalue"})
				return &m
			}(),
			wantErr: false,
		},
		{
			name: "skips empty lines",
			args: args{
				r: strings.NewReader(`foo=bar
somekey=value

anotherkey=alsovalue`),
			},
			want: func() *ConfigMap {
				m := ConfigMap(map[string]string{"foo": "bar", "somekey": "value", "anotherkey": "alsovalue"})
				return &m
			}(),
			wantErr: false,
		},
		{
			name: "skips comments",
			args: args{
				r: strings.NewReader(`foo=bar
somekey=value
# anotherkey is important
anotherkey=alsovalue`),
			},
			want: func() *ConfigMap {
				m := ConfigMap(map[string]string{"foo": "bar", "somekey": "value", "anotherkey": "alsovalue"})
				return &m
			}(),
			wantErr: false,
		},
		{
			name: "fails if line isn't on the form key=val",
			args: args{
				r: strings.NewReader(`foo=bar
somekey=value
badline
anotherkey=alsovalue`),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseConfig(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
