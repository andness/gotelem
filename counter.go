package gotelem

import (
	"io"
	"sync/atomic"
	"time"
)

type Counter struct {
	name string
	*Sampler
	broadcaster
	countSummarizers []*SlidingWindowSummarizer
	deltaSummarizers []*SlidingWindowSummarizer
	rateUnit         string
	count            int64
	prevSample       int64
}

func NewCounter(name string, samplingInterval time.Duration, summarizerWindows []time.Duration, httpPublisher *HTTPPublisher, log io.Writer) (counter *Counter) {
	counter = &Counter{
		name:     name,
		rateUnit: rateUnit(samplingInterval)}
	// TODO(go1.1)
	// We'll need this until Go 1.1 allows us to pass methods around
	// just like funcs
	sample := func(t time.Time) {
		counter.sample(t)
	}
	if httpPublisher != nil {
		counter.AddReceiver(httpPublisher)
	}
	if log != nil {
		counter.AddReceiver(newLogger(log))
	}
	if samplingInterval != 0 {
		counter.Sampler = NewSampler(samplingInterval, sample)
		counter.countSummarizers, counter.deltaSummarizers = counter.makeSummarizers(summarizerWindows)
	}
	return
}

func (c *Counter) Inc() {
	c.count = atomic.AddInt64(&c.count, 1)
}

func (c *Counter) Dec() {
	c.count = atomic.AddInt64(&c.count, -1)
}

func (c *Counter) makeSummarizers(windows []time.Duration) (countSummarizers, deltaSummarizers []*SlidingWindowSummarizer) {
	countSummarizers = make([]*SlidingWindowSummarizer, len(windows))
	deltaSummarizers = make([]*SlidingWindowSummarizer, len(windows))
	for i, w := range windows {
		countSummarizers[i] = NewSlidingWindowSummarizer(c.name, w)
		deltaSummarizers[i] = NewSlidingWindowSummarizer(c.name+"/"+c.rateUnit, w)
	}
	return
}

func (c *Counter) sample(t time.Time) {
	sampledCount := c.count
	delta := sampledCount - c.prevSample
	c.prevSample = sampledCount

	observation := &Observation{Timestamp: t, Name: c.name, Value: float64(sampledCount)}
	deltaObservation := &Observation{Timestamp: t, Name: c.name + "/" + c.rateUnit, Value: float64(delta)}

	//c.httpPublisher.Add(observation)
	//c.logObservation(observation)
	//c.httpPublisher.Add(deltaObservation)
	//c.logObservation(deltaObservation)
	c.broadcast(observation)
	c.broadcast(deltaObservation)
	for _, s := range c.countSummarizers {
		s.Update(observation)
		for _, obs := range s.Summarize() {
			//c.httpPublisher.Add(obs)
			//c.logObservation(obs)
			c.broadcast(obs)
		}
	}
	for _, s := range c.deltaSummarizers {
		s.Update(deltaObservation)
		for _, obs := range s.Summarize() {
			//c.httpPublisher.Add(obs)
			//c.logObservation(obs)
			c.broadcast(obs)
		}
	}
}

func rateUnit(interval time.Duration) (unit string) {
	switch interval {
	case time.Nanosecond:
		unit = "ns"
	case time.Microsecond:
		unit = "us"
	case time.Millisecond:
		unit = "ms"
	case time.Second:
		unit = "sec"
	case time.Minute:
		unit = "min"
	case time.Hour:
		unit = "hour"
	default:
		unit = interval.String()
	}
	return
}
