package main

import (
	"time"
)

type SlidingWindow struct {
	Timeframe time.Duration
	MaxReqNum int
	Requests  []time.Time
}

func NewSlidingWindow(windowSize time.Duration, maxReqNum int) *SlidingWindow {
	return &SlidingWindow{
		Timeframe: windowSize,
		MaxReqNum: maxReqNum,
		Requests:  []time.Time{},
	}
}

func (sw *SlidingWindow) SlideWindow(time time.Time) {
	threshold := time.Add(-sw.Timeframe)
	for len(sw.Requests) > 0 && sw.Requests[0].Before(threshold) {
		sw.Requests = sw.Requests[1:]
	}
}

func (sw *SlidingWindow) AllowRequest() bool {
	now := time.Now()
	sw.SlideWindow(now)

	if len(sw.Requests) >= sw.MaxReqNum {
		return false
	}

	sw.Requests = append(sw.Requests, now)
	return true
}

// func main() {
// 	sw := NewSlidingWindow(time.Minute, 5) // Allow N requests per timeframe

// 	for i := 0; i < 10; i++ {
// 		if sw.AllowRequest() {
// 			fmt.Printf("Request %d: Success\n", i+1)
// 		} else {
// 			fmt.Printf("Request %d: Limit Exceeded\n", i+1)
// 		}
// 		time.Sleep(time.Second)
// 	}
// }
