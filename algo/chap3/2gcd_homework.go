package main

// Find gcd of 3 numbers a, b and c

func gcd3Nums(a, b, c int) int {
	return gcdModulo(gcdModulo(a, b), c)
}
