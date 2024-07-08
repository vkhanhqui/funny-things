package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPubSub(t *testing.T) {
	ps := NewPubSub()

	go ps.Pub(50)
	go ps.Sub("id1")
	go ps.Sub("id2")

	time.Sleep(1 * time.Millisecond)
	ps.Stop(2)
	assert.Equal(t, 1225, ps.GetTotal())
}
