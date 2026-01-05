package main

import "testing"

func TestIsPrimeSqrt(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want bool
	}{
		{"negative number", -7, false},
		{"zero", 0, false},
		{"one", 1, false},
		{"two", 2, true},
		{"three", 3, true},
		{"four", 4, false},
		{"small prime", 5, true},
		{"small composite", 9, false},
		{"even composite", 10, false},
		{"large prime", 29, true},
		{"perfect square", 49, false},
		{"large composite", 77, false},
		{"large prime near square", 97, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isPrimeSqrt(tt.n)
			if got != tt.want {
				t.Fatalf("isPrimeSqrt(%d) = %v; want %v", tt.n, got, tt.want)
			}
		})
	}
}
