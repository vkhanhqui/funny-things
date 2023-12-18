package algo

import (
	"time"
)

type SlidingWindow struct {
	Timeframe time.Duration
	MaxReqNum int
	Requests  []time.Time
}

func NewSlidingWindow(maxReqNum int, windowSize time.Duration) *SlidingWindow {
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

func (sw *SlidingWindow) SendRequest() bool {
	now := time.Now()
	sw.SlideWindow(now)

	if len(sw.Requests) >= sw.MaxReqNum {
		return false
	}

	sw.Requests = append(sw.Requests, now)
	return true
}
