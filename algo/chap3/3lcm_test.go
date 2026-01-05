package main

import "testing"

func TestLCMBasicCases(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{"both positive", 12, 18, 36},
		{"coprime", 7, 13, 91},
		{"one divides the other", 6, 3, 6},
		{"same numbers", 10, 10, 10},
		{"with zero a", 0, 5, 0},
		{"with zero b", 5, 0, 0},
		{"both zero", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lcm(tt.a, tt.b)
			if got != tt.want {
				t.Fatalf("lcm(%d, %d) = %d; want %d",
					tt.a, tt.b, got, tt.want)
			}
		})
	}
}
