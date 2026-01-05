package main

// BCNN(a, b) = a * b / UCLN(a, b)

func lcm(a, b int) int {
	if a == 0 || b == 0 {
		return 0
	}
	return a * b / gcdModulo(a, b)
}
