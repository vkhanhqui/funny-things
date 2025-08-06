package funnythings

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestTokenBucket_WithoutSleep(t *testing.T) {
	bucket := NewTokenBucket(3, time.Second, 1)
	for range 3 {
		if !bucket.Allow() {
			t.Fatal("123")
		}
	}

	if bucket.Allow() {
		t.Fatal("456")
	}
}

func TestTokenBucket_WithSleep(t *testing.T) {
	bucket := NewTokenBucket(3, time.Second, 1)
	for range 3 {
		if !bucket.Allow() {
			t.Fatal("123")
		}
	}

	// tokens = 0
	time.Sleep(time.Second + 10*time.Millisecond)

	// tokens = 1
	if !bucket.Allow() {
		t.Fatal("456")
	}

	// tokens = 0
	if bucket.Allow() {
		t.Fatal("678")
	}
}

func TestTokenBucket_Concurrency(t *testing.T) {
	bucket := NewTokenBucket(3, time.Second, 1)
	expectedAllow := int64(3)
	var actualAllow atomic.Int64

	var wg sync.WaitGroup

	for range 10 {
		wg.Add(1)
		go func() {
			if bucket.Allow() {
				actualAllow.Add(1)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if actualAllow.Load() != expectedAllow {
		t.Fail()
	}
}
