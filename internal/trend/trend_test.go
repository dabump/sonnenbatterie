package trend

import "testing"

func Test_calculate(t *testing.T) {
	type args struct {
		values []int
	}
	tests := []struct {
		name string
		args args
		want Trend
	}{
		{
			name: "Upward trend",
			args: args{
				values: []int{4, 3, 2, 1},
			},
			want: Upward,
		},
		{
			name: "Downward trend",
			args: args{
				values: []int{1, 2, 3, 4},
			},
			want: Downward,
		},
		{
			name: "No trend - Equalise",
			args: args{
				values: []int{1, 2, 2, 1},
			},
			want: NoTrend,
		},
		{
			name: "No trend - Too few values",
			args: args{
				values: []int{1},
			},
			want: NoTrend,
		},
		{
			name: "No trend - No values",
			args: args{
				values: []int{},
			},
			want: NoTrend,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Calculate(tt.args.values); got != tt.want {
				t.Errorf("calculate() = %v, want %v", got, tt.want)
			}
		})
	}
}
