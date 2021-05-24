package main

import "testing"

func TestEqual(t *testing.T) {
	type args struct {
		a []string
		b []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty",
			args: args{
				a: []string{},
				b: []string{},
			},
			want: true,
		},
		{
			name: "the same",
			args: args{
				a: []string{"a"},
				b: []string{"a"},
			},
			want: true,
		},
		{
			name: "the same multiple",
			args: args{
				a: []string{"a", "b"},
				b: []string{"a", "b"},
			},
			want: true,
		},
		{
			name: "different",
			args: args{
				a: []string{"a"},
				b: []string{"b"},
			},
			want: false,
		},
		{
			name: "different multiple",
			args: args{
				a: []string{"a", "b"},
				b: []string{"a", "b", "c"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Equal(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}
