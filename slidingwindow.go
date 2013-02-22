package gotelem

import (
	"math"
	"time"
)

// A simple sliding window implementation that will keep items until
// they are older than the maxAge parameter. Note that there is just a
// single Add() method, and no way to list items. The reason for this
// is that we don't want to iterate over the items to compute
// summaries, instead we maintain the summaries by diffing the
// incoming and outgoing items every time we Add(). If we had an Items
// function we could get inconsistencies. What if we call Add() 10
// times, but before we call Items() one or more of those items are
// expired. This would lead to inconsistency between our continously
// maintained summary and the summary we would get if we summarized
// over Items().
//
// The implementation is *not* thread-safe.
type slidingWindow struct {
	items    []*Observation
	maxAge   time.Duration
	timeNow  func() time.Time
	oldestAt int
	insertAt int
}

// Creates a new sliding window which keeps Observations until their
// timestamp is stricly less than than time.Now() - maxAge. Uses
// DefaultTimeKeeper as the time source.
func NewSlidingWindow(maxAge time.Duration) *slidingWindow {
	return &slidingWindow{
		maxAge:  maxAge,
		items:   make([]*Observation, 0, 16),
		timeNow: time.Now}
}

// Adds an item to the window and returns any items that were evicted
// because they were too old.
func (w *slidingWindow) Add(o *Observation) (expired []*Observation) {
	w.items = append(w.items, o)
	//println("Add", o.Value)
	expired = w.findExpired()
	w.compact()
	return
}

// Returns the min and max observation values in the window
func (w *slidingWindow) minmax() (min float64, max float64) {
	min = math.MaxFloat64
	for _, o := range w.items[w.oldestAt:] {
		if o.Value < min {
			min = o.Value
		}
		if o.Value > max {
			max = o.Value
		}
	}
	return
}

func (w *slidingWindow) findExpired() (expired []*Observation) {
	oldestAcceptable := w.timeNow().Add(-w.maxAge)
	//println("acc", oldestAcceptable.String())
	for ; w.oldestAt < len(w.items) && w.items[w.oldestAt].Timestamp.Before(oldestAcceptable); w.oldestAt++ {
		//println("time", w.items[w.oldestAt].Timestamp.String(), "val", w.items[w.oldestAt].Value)
		expired = append(expired, w.items[w.oldestAt])
	}
	return
}

// Move the sliding window elements to the start of the items so that
// we don't leak memory. Compacting only happens when the oldestAt
// index is at least halfway to the capacity of the items. Since we
// generally deal with data being added at a fixed rate this means the
// items should stabilize at a sufficient size with no more
// allocations needed.
func (w *slidingWindow) compact() {
	if w.oldestAt < cap(w.items)-w.oldestAt {
		return
	}
	//   0   1   2   3   4   5   6   7   8   9  10  11  12]
	// [ x   x   x   x   x   x   7   -   -   -   -   -   -]
	// cap = 12, oldestAt = 6, len = 7
	items := len(w.items)
	copy(w.items, w.items[w.oldestAt:])     // [7 x x x x x 7 - - - - -]
	w.items = w.items[0 : items-w.oldestAt] // [0:1] -> [7]
	w.oldestAt = 0
}
