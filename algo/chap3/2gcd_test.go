package main

import (
	"testing"
)

func TestGCD(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{"both positive", 48, 18, 6},
		{"coprime", 7, 13, 1},
		{"one is zero", 0, 5, 5},
		{"both zero", 0, 0, 0},
		{"same numbers", 10, 10, 10},
		{"negative a", -24, 18, 6},
		{"negative b", 24, -18, 6},
		{"both negative", -24, -18, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gcdModulo(tt.a, tt.b)
			if got != tt.want {
				t.Fatalf("gcd(%d, %d) = %d; want %d",
					tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestGCD3Nums(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		c    int
		want int
	}{
		{"all positive", 48, 18, 30, 6},
		{"coprime triple", 7, 13, 19, 1},
		{"one divides all", 6, 12, 18, 6},
		{"same numbers", 10, 10, 10, 10},
		{"with zero", 0, 12, 18, 6},
		{"two zeros", 0, 0, 9, 9},
		{"all zeros", 0, 0, 0, 0},
		{"negative values", -24, 18, -6, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gcd3Nums(tt.a, tt.b, tt.c)
			if got != tt.want {
				t.Fatalf(
					"gcd3Nums(%d, %d, %d) = %d; want %d",
					tt.a, tt.b, tt.c, got, tt.want,
				)
			}
		})
	}
}
