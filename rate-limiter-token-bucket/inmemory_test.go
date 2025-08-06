package funnythings

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestInMemory_Sweep(t *testing.T) {
	userID1 := "1234"
	userID2 := "5678"
	m := NewInMemory(3, time.Second, 1)
	m.Allow(userID1)
	m.Allow(userID2)

	if !m.Exists(userID1) {
		t.Fatal("123")
	}
	if !m.Exists(userID2) {
		t.Fatal("456")
	}

	time.Sleep(5*time.Second + 10*time.Millisecond)

	// removed
	if m.Exists(userID1) {
		t.Fatal("123")
	}
	if m.Exists(userID2) {
		t.Fatal("456")
	}
}

func TestInMemory_WithoutSleep(t *testing.T) {
	userID1 := "1234"
	userID2 := "5678"
	m := NewInMemory(3, time.Second, 1)
	for range 3 {
		if !m.Allow(userID1) {
			t.Fatal("123")
		}
		if !m.Allow(userID2) {
			t.Fatal("123")
		}
	}

	if m.Allow(userID1) {
		t.Fatal("456")
	}
	if m.Allow(userID2) {
		t.Fatal("456")
	}
}

func TestInMemory_WithSleep(t *testing.T) {
	userID1 := "1234"
	userID2 := "5678"
	m := NewInMemory(3, time.Second, 1)
	for range 3 {
		if !m.Allow(userID1) {
			t.Fatal("123")
		}
		if !m.Allow(userID2) {
			t.Fatal("123")
		}
	}

	// tokens = 0
	time.Sleep(time.Second + 10*time.Millisecond)

	// tokens = 1
	if !m.Allow(userID1) {
		t.Fatal("456")
	}
	if !m.Allow(userID2) {
		t.Fatal("456")
	}

	// tokens = 0
	if m.Allow(userID1) {
		t.Fatal("678")
	}
	if m.Allow(userID2) {
		t.Fatal("678")
	}
}

func TestInMemory_Concurrency(t *testing.T) {
	userID1 := "1234"
	userID2 := "5678"
	m := NewInMemory(3, time.Second, 1)

	expectedAllow := int64(6)
	var actualAllow atomic.Int64

	var wg sync.WaitGroup

	for range 10 {
		wg.Add(1)
		go func() {
			if m.Allow(userID1) {
				actualAllow.Add(1)
			}
			if m.Allow(userID2) {
				actualAllow.Add(1)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if actualAllow.Load() != expectedAllow {
		t.Fatal(actualAllow.Load())
	}
}
