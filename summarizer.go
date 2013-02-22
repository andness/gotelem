package gotelem

import (
	"fmt"
	"math"
	"time"
)

// A summarizer computes avg/sum/count/min/max. It subscribes to
// observations and publishes the derived values. Typically it is
// subscribed to by a HTTPPublisher which collects up to X time worth
// of observations.
type SlidingWindowSummarizer struct {
	name    string
	suffix  string
	window  *slidingWindow
	sum     float64
	count   int64
	min     float64
	max     float64
	avg     float64
	timeNow func() time.Time
}

func suffix(d time.Duration) string {
	if d < time.Minute {
		return d.String()
	} else if d < time.Hour {
		return fmt.Sprintf("%dM", int(d/time.Minute))
	} else {
		return fmt.Sprintf("%dH", int(d/time.Hour))
	}
	panic("Ha ha ha")
}

func NewSlidingWindowSummarizer(name string, maxAge time.Duration) *SlidingWindowSummarizer {
	// TODO: Limit window size to a clean multiple of time.Minute
	summarizer := &SlidingWindowSummarizer{
		name:    name,
		suffix:  suffix(maxAge),
		window:  NewSlidingWindow(maxAge),
		min:     math.MaxFloat64,
		max:     -math.MaxFloat64,
		avg:     math.NaN(),
		timeNow: time.Now}
	return summarizer
}

func (s *SlidingWindowSummarizer) Summarize() []*Observation {
	now := s.timeNow().UTC()
	return []*Observation{
		&Observation{now, s.name + ":" + s.suffix + "_MIN", s.min},
		&Observation{now, s.name + ":" + s.suffix + "_MAX", s.max},
		&Observation{now, s.name + ":" + s.suffix + "_SUM", s.sum},
		&Observation{now, s.name + ":" + s.suffix + "_AVG", s.avg},
		&Observation{now, s.name + ":" + s.suffix + "_COUNT", float64(s.count)}}
}

func minMaxObservation(s []*Observation) (min, max float64) {
	min = math.MaxFloat64
	max = -math.MaxFloat64
	for _, o := range s {
		if o.Value < min {
			min = o.Value
		}
		if o.Value > max {
			max = o.Value
		}
	}
	return
}

// Every time we receive a new observation we immediately update our
// internals. This is done without having to scan the entire list of
// observations each time.
func (s *SlidingWindowSummarizer) Update(o *Observation) {

	// TODO: If an already expired item is added we get into a weird
	// state We should perhaps expire items when Summarize() is called
	// to?
	expired := s.window.Add(o)

	expiredMin, expiredMax := minMaxObservation(expired)

	// We can avoid recomputing min/max if the new observed value is
	// smaller than the current min/max.

	// And unless the current min/max was in the expired set, we don't
	// need to recompute.

	// But when we do recompute we might as well find both min and max
	// since we have to rescan everything anyway.

	// This does not account for the (highly unlikely) case of an
	// expired item being added to the window.
	recomputeMinMax := false
	if o.Value < s.min {
		s.min = o.Value
	} else if expiredMin == s.min {
		recomputeMinMax = true
	}

	if o.Value > s.max {
		s.max = o.Value
	} else if expiredMax == s.max {
		recomputeMinMax = true
	}
	if recomputeMinMax {
		s.min, s.max = s.window.minmax()
	}

	s.sum += o.Value
	for _, e := range expired {
		s.sum -= e.Value
	}

	s.count += int64(1 - len(expired))

	s.avg = math.NaN()
	if s.count != 0 {
		s.avg = s.sum / float64(s.count)
	}

}
