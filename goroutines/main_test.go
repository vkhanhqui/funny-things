package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestIntro(t *testing.T) {
	words := []string{
		"one",
		"two",
		"three",
	}
	for i, w := range words {
		fmt.Println(i, w)
	}
}

func TestFirstGoroutines(t *testing.T) {
	words := []string{
		"one",
		"two",
		"three",
	}
	for _, w := range words {
		go fmt.Println(w)
	}
}

func TestTrySleeping(t *testing.T) {
	words := []string{
		"one",
		"two",
		"three",
	}
	for _, w := range words {
		go fmt.Println(w)
	}
	time.Sleep(time.Second)
}

func TestFunctionExecutedMoreThanSleepingTime(t *testing.T) {
	printAny := func(value any) {
		time.Sleep(2 * time.Second) // assume we do a heavy processing which takes 2 seconds
		fmt.Println(value)
	}
	words := []string{
		"one",
		"two",
		"three",
	}
	for _, w := range words {
		go printAny(w)
	}
	time.Sleep(time.Second)
}

func TestWaitGroup(t *testing.T) {
	var wg sync.WaitGroup

	printAny := func(value any) {
		time.Sleep(2 * time.Second) // assume we do a heavy processing which takes 2 seconds
		fmt.Println(value)
		wg.Done() // make this function complete by decreasing the WaitGroup by one

	}
	words := []string{
		"one",
		"two",
		"three",
	}

	wg.Add(len(words)) // define the number of goroutines that will run
	for _, w := range words {
		go printAny(w)
	}
	wg.Wait()
}

func TestReproduceRaceCondition(t *testing.T) {
	var wg sync.WaitGroup

	name := "Qui"
	updateName := func(newName string) {
		fmt.Printf("Start updating from '%s' to '%s'\n", name, newName)
		name = newName
		wg.Done()
	}

	wg.Add(3)
	go updateName("Khanh Qui")
	go updateName("Khanh Qui Vo")
	go updateName("Steve")
	wg.Wait()
	fmt.Printf("Final result: '%s'\n", name)
}

func TestFixRaceCondition(t *testing.T) {
	var wg sync.WaitGroup

	name := "Qui"
	updateName := func(newName string) {
		fmt.Printf("Start updating from '%s' to '%s'\n", name, newName)
		name = newName
		wg.Done()
	}

	wg.Add(1)
	go updateName("Khanh Qui")
	wg.Wait()

	wg.Add(1)
	go updateName("Khanh Qui Vo")
	wg.Wait()

	wg.Add(1)
	go updateName("Steve")
	wg.Wait()
	fmt.Printf("Final result: '%s'\n", name)
}

func TestMutex(t *testing.T) {
	var wg sync.WaitGroup
	var lock sync.Mutex

	name := "Qui"
	updateName := func(newName string) {
		lock.Lock()
		fmt.Printf("Start updating from '%s' to '%s'\n", name, newName)
		name = newName
		lock.Unlock()
		wg.Done()
	}

	wg.Add(3)
	go updateName("Khanh Qui")
	go updateName("Khanh Qui Vo")
	go updateName("Steve")
	wg.Wait()
	fmt.Printf("Final result: '%s'\n", name)
}
