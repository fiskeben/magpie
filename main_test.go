package main

import "testing"

func Test_mask(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "masks a long string",
			args: args{
				value: "This is a long string to mask",
			},
			want: "Th*************************sk",
		},
		{
			name: "masks a short string",
			args: args{
				value: "12345",
			},
			want: "*****",
		},
		{
			name: "masks a string barely long enough",
			args: args{
				value: "123456",
			},
			want: "12**56",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mask(tt.args.value); got != tt.want {
				t.Errorf("mask() = %v, want %v", got, tt.want)
			}
		})
	}
}
