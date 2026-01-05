package main

import "testing"

func TestHomework4ShortenFraction(t *testing.T) {
	tests := []struct {
		name  string
		a, b  int
		wantA int
		wantB int
	}{
		{"already simplified", 3, 4, 3, 4},
		{"common divisor", 6, 8, 3, 4},
		{"large common divisor", 48, 18, 8, 3},
		{"numerator zero", 0, 5, 0, 1},
		{"negative numerator", -6, 8, -3, 4},
		{"negative denominator", 6, -8, 3, -4},
		{"both negative", -6, -8, 3, 4},
		{"same numbers", 10, 10, 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotA, gotB := homework4ShortenFraction(tt.a, tt.b)
			if gotA != tt.wantA || gotB != tt.wantB {
				t.Fatalf(
					"shorten(%d, %d) = (%d, %d); want (%d, %d)",
					tt.a, tt.b, gotA, gotB, tt.wantA, tt.wantB,
				)
			}
		})
	}
}

func TestHomework4TwoCoprimeIntegers(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want bool
	}{
		{"coprime small", 7, 13, true},
		{"coprime one is one", 1, 99, true},
		{"not coprime", 12, 18, false},
		{"same numbers > 1", 10, 10, false},
		{"one zero", 0, 1, true},
		{"zero and non-coprime", 0, 6, false},
		{"both zero", 0, 0, false},
		{"negative values coprime", -7, 13, true},
		{"negative values not coprime", -12, 18, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := homework4TwoCoprimeIntegers(tt.a, tt.b)
			if got != tt.want {
				t.Fatalf(
					"homework4TwoCoprimeIntegers(%d, %d) = %v; want %v",
					tt.a, tt.b, got, tt.want,
				)
			}
		})
	}
}
