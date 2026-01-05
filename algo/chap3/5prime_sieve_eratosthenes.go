package main

func isPrimeSieveEratosThenes(n int) bool {
	if n <= 1 {
		return false
	}

	isPrime := primeSieveEratosThenes(n)
	return isPrime[n]
}

func primeSieveEratosThenes(n int) []bool {
	isPrime := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		isPrime[i] = true
	}

	for i := 2; i*i <= n; i++ {
		if isPrime[i] {
			for j := i * i; j <= n; j += i {
				isPrime[j] = false
			}
		}
	}

	return isPrime
}
