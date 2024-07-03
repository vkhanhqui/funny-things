package service

import (
	"testing"
	"time"
)

func TestPubSub(t *testing.T) {
	ps := NewPubSub()

	go ps.Pub(50)
	go ps.Sub("id1")
	go ps.Sub("id2")

	time.Sleep(1 * time.Millisecond)
	ps.Stop(2)
}
