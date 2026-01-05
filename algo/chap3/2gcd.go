package main

// Euclid algo is the most optimised algo
// to calculate "Greatest common divisor" (aka "GCD" or "UCLN")

func gcdSubtraction(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}

	for a != b {
		if a > b {
			a -= b
		} else {
			b -= a
		}
	}
	return a
}

func gcdModulo(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}

	for b != 0 {
		a, b = b, a%b
	}
	return a
}
