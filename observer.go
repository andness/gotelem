package gotelem

import (
	"io"
	"time"
)

type Observer struct {
	name string
	*Sampler
	broadcaster
	summarizers []*SlidingWindowSummarizer
	timeNow     func() time.Time
}

func (o *Observer) Observe(value float64) {
	obs := &Observation{o.timeNow().UTC(), o.name, value}
	for _, s := range o.summarizers {
		s.Update(obs)
	}
	o.broadcast(obs)
}

func NewObserver(name string, samplingInterval time.Duration, summarizerWindows []time.Duration, httpPublisher *HTTPPublisher, log io.Writer) (observer *Observer) {
	observer = &Observer{
		name:    name,
		timeNow: time.Now}
	// TODO(go1.1)
	// :( http://code.google.com/p/go/issues/detail?id=2280
	sample := func(t time.Time) {
		observer.sample(t)
	}
	if samplingInterval != 0 {
		// TODO: There is no sense in having summarizers without a sampler since they won't be sampled but the argument list as it is allows you to specify this
		observer.Sampler = NewSampler(samplingInterval, sample)
		observer.summarizers = make([]*SlidingWindowSummarizer, len(summarizerWindows))
		for i, windowSize := range summarizerWindows {
			observer.summarizers[i] = NewSlidingWindowSummarizer(name, windowSize)
		}
	}
	if httpPublisher != nil {
		observer.AddReceiver(httpPublisher)
	}
	if log != nil {
		observer.AddReceiver(newLogger(log))
	}
	return
}

func (o *Observer) sample(t time.Time) {
	for _, s := range o.summarizers {
		for _, obs := range s.Summarize() {
			o.broadcast(obs)
		}
	}
}
