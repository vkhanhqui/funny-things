package main

import (
	"log"
	"time"
	"worker-pattern/service"
)

func main() {
	log.SetFlags(log.Ltime)

	totalWorker := 5
	wp := service.NewWorker(totalWorker)
	wp.Run()

	type square struct {
		side   int
		result int
	}

	totalTask := 10
	squareC := make(chan square, totalTask)

	for i := 0; i < totalTask; i++ {
		side := i
		wp.AddTask(func() {
			log.Printf("Calculate square with side: %d", side)
			squareC <- square{side, side * side}
			time.Sleep(5 * time.Second)
		})
	}

	for i := 0; i < totalTask; i++ {
		res := <-squareC
		log.Printf("Square with side %d = %d", res.side, res.result)
	}
}
