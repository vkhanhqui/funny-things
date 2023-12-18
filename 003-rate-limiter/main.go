package main

import (
	"fmt"
	"rate-limiter/algo"
	"rate-limiter/app"
	"time"
)

func main() {
	// fixedWindow := algo.NewFixedWindow(5, 6*time.Second)
	// leakyBucket := algo.NewLeakyBucket(5, 1, 6*time.Second)
	// slidingWindow := algo.NewSlidingWindow(5, 6*time.Second)
	tokenBucket := algo.NewTokenBucket(5, 1, 6*time.Second)
	l := app.NewLimiter(tokenBucket)

	for i := 0; i < 10; i++ {
		if l.SendRequest("user_1") {
			fmt.Printf("Request %d: Success\n", i+1)
		} else {
			fmt.Println("Limit Exceeded")
		}

		time.Sleep(time.Second)
	}
}
