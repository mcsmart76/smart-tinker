package timerwheel

import (
	"testing"
	"time"
)

func TestNewClosedTimerWheel(t *testing.T) {
	var w *ClosedTimerWheel

	w = NewClosedTimerWheel(0, 1000)
	if w != nil {
		t.Fatalf("Expected nil, got %v", w)
	}

	w = NewClosedTimerWheel(-10 * time.Millisecond, 1000)
	if w != nil {
		t.Fatalf("Expected nil, got %v", w)
	}

	w = NewClosedTimerWheel(5 * time.Millisecond, 0)
	if w != nil {
		t.Fatalf("Expected nil, got %v", w)
	}

	w = NewClosedTimerWheel(5 * time.Millisecond, -10)
	if w != nil {
		t.Fatalf("Expected nil, got %v", w)
	}

	w = NewClosedTimerWheel(5 * time.Millisecond, 1000)
	if w == nil {
		t.Fatal("Expected not nil")
	}
}

func TestClosedTimerWheelMaxTimeout(t *testing.T) {
	w := NewClosedTimerWheel(5 * time.Millisecond, 1000)
	if w == nil {
		t.Fatal("Expected not nil")
	}
	if w.MaxTimeout() != 5 * time.Second {
		t.Fatalf("Expected 5 seconds, got %v", w.MaxTimeout())
	}
}

func TestClosedTimerWheelSchedule(t *testing.T) {
	w := NewClosedTimerWheel(5 * time.Millisecond, 10)
	ch := make(chan int)
	handle1 := w.Schedule(30 * time.Millisecond, func() {
		ch <- 2
	})
	handle2 := w.Schedule(10 * time.Millisecond, func() {
		ch <- 1
	})
	if handle1 >= handle2 {
		t.Fatalf("Expected handle1 (%v) < handle2 (%v)", handle1, handle2)
	}

	timeout := time.After(2 * w.MaxTimeout())
	firstTimeoutReceived := false
	for {
		select {
		case v := <-ch:
			// TODO(mattsmart): Convert to switch
			if v == 1 {
				if firstTimeoutReceived {
					t.Fatalf("Got first timeout late")
				}
				firstTimeoutReceived = true
			} else if v == 2 && !firstTimeoutReceived {
				t.Fatalf("Got second timeout early")
			} else {
				t.Fatalf("Did not expect %v, error in test", v)
			}
		case <-timeout:
			t.Fatal("Timed out!")
		}
	}
}

func TestClosedTimerWheelSlot(t *testing.T) {
	w := NewClosedTimerWheel(5 * time.Millisecond, 10)
	if v := w.slot(9 * time.Millisecond); v != 1 {
		t.Fatalf("Got %v, expected 1", v)
	}
	if v := w.slot(10 * time.Millisecond); v != 2 {
		t.Fatalf("Got %v, expected 2", v)
	}
	if v := w.slot(11 * time.Millisecond); v != 2 {
		t.Fatalf("Got %v, expected 2", v)
	}
	if v := w.slot(25 * time.Millisecond); v != 5 {
		t.Fatalf("Got %v, expected 5", v)
	}
	if v := w.slot(100 * time.Millisecond); v != 0 {
		t.Fatalf("Got %v, expected 0", v)
	}

	if v := w.slot(0); v != 0 {
		t.Fatalf("Got %v, expected 0", v)
	}
	w.tick()
	w.tick()
	w.tick()
	if v := w.slot(0); v != 3 {
		t.Fatalf("Got %v, expected 3", v)
	}
}
