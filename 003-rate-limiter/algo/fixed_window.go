package algo

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

func (fw *FixedWindow) SendRequest() bool {
	fw.resetWindow()

	if fw.CurReqNum < fw.MaxReqNum {
		fw.CurReqNum += 1
		return true
	}
	return false
}

func (fw *FixedWindow) resetWindow() {
	now := time.Now()
	if now.After(fw.WindowEnd) {
		fw.WindowEnd = now.Add(fw.TimeFrame)
		fw.CurReqNum = 0
	}
}
