package main

import (
	"time"
)

type FixedWindow struct {
	WindowEnd time.Time
	MaxReqNum int
	CurReqNum int
	TimeFrame time.Duration
}

func NewFixedWindow(maxRequestNum int, timeframe time.Duration) *FixedWindow {
	return &FixedWindow{
		WindowEnd: time.Now().Add(timeframe),
		MaxReqNum: maxRequestNum,
		CurReqNum: 0,
		TimeFrame: timeframe,
	}
}

func (fw *FixedWindow) SendRequests(reqNum int) bool {
	fw.ResetWindow()

	if fw.CurReqNum < fw.MaxReqNum {
		fw.CurReqNum += reqNum
		return true
	}
	return false
}

func (fw *FixedWindow) ResetWindow() {
	now := time.Now()
	if now.After(fw.WindowEnd) {
		fw.WindowEnd = now.Add(fw.TimeFrame)
		fw.CurReqNum = 0
	}
}

// func main() {
// 	fw := NewFixedWindow(10, time.Minute) // Max N requests per timeframe

// 	for i := 0; i < 10; i++ {
// 		if fw.SendRequests(2) { // send 2 requests per time
// 			fmt.Printf("Request %d: Success\n", i+1)
// 		} else {
// 			fmt.Printf("Request %d: Limit Exceeded\n", i+1)
// 		}
// 		time.Sleep(time.Second)
// 	}
// }
