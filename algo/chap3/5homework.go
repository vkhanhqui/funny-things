package main

/*
1. Count primes in range a and b

idea: use sieve eratosthenes method

2. Find the prime number at K position

idea: use sieve eratosthenes or sqrt method

3. Find prime factors of N

Example:
For N = 24
Prime Factors: 2, 2, 2, 3

Explain:
- 2 and 3 are primes
- and 2*2*2*3 = 24

idea + pseudo code:
- find primes in range 2 -> N
- while N % prime == 0 {
	n /= prime
	res.append(prime)
}

*/

func primeFactorsOfN(n int) []int {
	if n <= 1 {
		return []int{}
	}

	var primes []int
	isPrime := primeSieveEratosThenes(n)
	for i := 2; i <= n; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}

	var res []int
	for _, p := range primes {
		for n%p == 0 {
			n /= p
			res = append(res, p)
			if n == 1 {
				return res
			}
		}
	}
	return res
}

func countPrimesInRange(a, b int) int {
	if b <= 1 {
		return 0
	}

	isPrime := make([]bool, b+1)
	for i := 2; i <= b; i++ {
		isPrime[i] = true
	}

	for i := 2; i*i <= b; i++ {
		if isPrime[i] {
			for j := i * i; j <= b; j += i {
				isPrime[j] = false
			}
		}
	}

	start := a
	if a <= 1 {
		start = 2
	}

	res := 0
	for i := start; i <= b; i++ {
		if isPrime[i] {
			res++
		}
	}

	return res
}

func primeAtKSqrt(k int) int {
	MAX := 1_000_000
	for i := 2; i <= MAX; i++ {
		if isPrimeSqrt(i) {
			k -= 1
		}
		if k == 0 {
			return i
		}
	}
	return 0
}

func primeAtKSieveEratosthenes(k int) int {
	MAX := 1_000_000
	isPrime := make([]bool, MAX+1)
	for i := 2; i <= MAX; i++ {
		isPrime[i] = true
	}

	for i := 2; i*i <= MAX; i++ {
		if isPrime[i] {
			for j := i * i; j <= MAX; j += i {
				isPrime[j] = false
			}
		}
	}

	res := 0
	for i := 2; i <= MAX; i++ {
		if isPrime[i] {
			res++
		}
		if res == k {
			return i
		}
	}

	return 0
}
