package service

import (
	"fmt"
	"sync"
)

type PubSub interface {
	Pub(n int)
	Sub(id string)
	Stop(subNum int)
}

type pubSub struct {
	numbers chan int
	total   int
	quit    chan any
	wg      sync.WaitGroup
	mux     sync.Mutex
}

func NewPubSub() PubSub {
	return &pubSub{
		numbers: make(chan int),
		quit:    make(chan any),
		wg:      sync.WaitGroup{},
		mux:     sync.Mutex{},
	}
}

func (ps *pubSub) Pub(n int) {
	ps.wg.Add(1)
	fmt.Println("start pub")
	for i := 0; i < n; i++ {
		ps.numbers <- i
		fmt.Println("published", i)
	}
	close(ps.numbers)
	ps.wg.Done()
}

func (ps *pubSub) Sub(id string) {
	fmt.Println("start sub ", id)
	ps.wg.Add(1)
	for {
		select {
		case n := <-ps.numbers:
			ps.inc(id, n)
			fmt.Println(id, "total", ps.getTotal())
		case <-ps.quit:
			fmt.Println("quit")
			ps.wg.Done()
			return
		}
	}
}

func (ps *pubSub) Stop(subNum int) {
	for i := 0; i < subNum; i++ {
		ps.quit <- 1
	}
	ps.wg.Wait()
}

func (ps *pubSub) getTotal() int {
	ps.mux.Lock()
	defer ps.mux.Unlock()
	return ps.total
}

func (ps *pubSub) inc(id string, n int) {
	ps.mux.Lock()
	defer ps.mux.Unlock()
	fmt.Println(id, "inc", n)
	ps.total += n
}
