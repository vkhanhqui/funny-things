package main

import (
	"reflect"
	"testing"
)

func TestPrimeFactorsOfNBasic(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want []int
	}{
		{"n <= 1", 1, []int{}},
		{"prime number", 13, []int{13}},
		{"small composite", 12, []int{2, 2, 3}},
		{"power of prime", 16, []int{2, 2, 2, 2}},
		{"two primes", 21, []int{3, 7}},
		{"large composite", 84, []int{2, 2, 3, 7}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := primeFactorsOfN(tt.n)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("primeFactorsOfN(%d) = %v; want %v",
					tt.n, got, tt.want)
			}
		})
	}
}

func TestCountPrimesInRangeBasic(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{"single prime", 2, 2, 1},
		{"single composite", 4, 4, 0},
		{"small range", 1, 10, 4}, // 2,3,5,7
		{"range with no primes", 14, 16, 0},
		{"range including zero", 0, 10, 4},
		{"range including negatives", -10, 10, 4},
		{"larger range", 10, 30, 6}, // 11,13,17,19,23,29
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := countPrimesInRange(tt.a, tt.b)
			if got != tt.want {
				t.Fatalf(
					"countPrimesInRange(%d, %d) = %d; want %d",
					tt.a, tt.b, got, tt.want,
				)
			}
		})
	}
}

func TestPrimeAtK(t *testing.T) {
	tests := []struct {
		name string
		k    int
		want int
	}{
		{"first prime", 1, 2},
		{"second prime", 2, 3},
		{"third prime", 3, 5},
		{"fifth prime", 5, 11},
		{"tenth prime", 10, 29},
		{"twentieth prime", 20, 71},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := primeAtKSqrt(tt.k)
			if got != tt.want {
				t.Fatalf("primeAtKSqrt(%d) = %d; want %d", tt.k, got, tt.want)
			}

			got = primeAtKSieveEratosthenes(tt.k)
			if got != tt.want {
				t.Fatalf("primeAtKSieveEratosthenes(%d) = %d; want %d", tt.k, got, tt.want)
			}
		})
	}
}
