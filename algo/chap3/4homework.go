package main

/*
1. Description: Input numerator A and denominator B. Simplify the fraction A / B and print the result.

idea + pseudo code:
c = gcd(a,b)
a /= c
b /= c
return a/b

2. Description: Coprime integers

idea:
- A prime integer is a positive whole number greater than 1
that has only two distinct positive divisors: 1 and itself
=> Coprime integers => gcd(a,b) = 1
*/

func homework4ShortenFraction(a, b int) (int, int) {
	c := gcdModulo(a, b)
	a /= c
	b /= c

	if a < 0 && b < 0 {
		return -a, -b
	}
	return a, b
}

func homework4TwoCoprimeIntegers(a, b int) bool {
	return gcdModulo(a, b) == 1
}
