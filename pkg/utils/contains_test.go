package utils

import "testing"

func TestContains(t *testing.T) {
	type args struct {
		arr  interface{}
		item interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "case1",
			args: args{
				arr:  []int{1, 2, 3},
				item: 1,
			},
			want: true,
		},
		{
			name: "case2",
			args: args{
				arr:  []int{1, 2, 3},
				item: "1",
			},
			want: false,
		},
		{
			name: "case3",
			args: args{
				arr:  []uint{1, 2, 3},
				item: 1,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.arr, tt.args.item); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
