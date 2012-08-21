package timerwheel

import (
	"time"
	"log"
)

type Handle int64

type Callback func()

type TimeWheeler interface {
	// Schedule a callback to run at a duration in the future.
	// Not reschedulable (yet).  Returns a handle for cancelling the timer.
	// TODO(mattsmart): Add error
	Schedule(duration time.Duration, cb Callback) Handle

	// Cancel a timer given its handle.
	Cancel(handle Handle)
}

// ClosedTimerWheel is a very simple implementation of TimeWheeler.  Callbacks
// can be scheduled at most granularity * capacity into the future.
type ClosedTimerWheel struct {
	// The next available handle.  Does not deal with wraparound.
	nextHandle Handle

	// The time granularity of this timer wheel.	
        granularity time.Duration

	// ticker.C is a channel that will be written to every granularity time.
	ticker *time.Ticker

	// Wheel of callbacks, each granularity in duration between each other.
	// cap(wheel == the capacity of the timer wheel.
	// TODO(mattsmart): s/Callback/struct with Handle and Callback/ ?
	wheel [][]Callback

	callbacks map[Handle]Callback

	// 0 <= currentSlot < capacity
	currentSlot int
}

// NewClosedTimerWheel returns a new ClosedTimerWheel with the given
// granularity and capacity.
func NewClosedTimerWheel(granularity time.Duration, capacity int) *ClosedTimerWheel {
	if granularity <= 0 {
		// TODO(mattsmart): Turn these into returned errors.
		log.Printf("Invalid granularity %v", granularity)
		return nil
	}
	if capacity <= 0 {
		log.Printf("Invalid capacity %v", capacity)
		return nil
	}
	return &ClosedTimerWheel{
		nextHandle: 1,
		currentSlot: 0,
		granularity: granularity,
		wheel: make([][]Callback, capacity),
		ticker: time.NewTicker(granularity),
	}
}

func (t *ClosedTimerWheel) Schedule(duration time.Duration, cb Callback) Handle {
	h := t.nextHandle
	t.nextHandle++
	return h
}

func (t *ClosedTimerWheel) Cancel(handle Handle) {
	_, ok := t.callbacks[handle]
	if !ok {
		log.Println("Invalid handle:", handle)
		return
	}
	delete(t.callbacks, handle)
}

func (t *ClosedTimerWheel) MaxTimeout() time.Duration {
	return t.granularity * time.Duration(cap(t.wheel))
}

func (t *ClosedTimerWheel) slot(duration time.Duration) int {
	return (t.currentSlot + int(duration / t.granularity)) % cap(t.wheel)
}

func (t *ClosedTimerWheel) tick() {
	t.currentSlot++
}
