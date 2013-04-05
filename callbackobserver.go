package gotelem

import (
	"io"
	"time"
)

// A CallbackObserver is very similar to a regular Observer except it
// doesn't offer an Observe method. Instead it is created with a
// callback function which is called every time the observer is
// sampled. The observatons returned from the callback are broadcast
// like with a regular observer.
type CallbackObserver struct {
	*Sampler
	broadcaster
}

func NewCallbackObserver(callback func(time.Time) []*Observation, samplingInterval time.Duration, summarizerWindows []time.Duration, httpPublisher receiver, log io.Writer) (observer *CallbackObserver) {
	observer = &CallbackObserver{}
	if samplingInterval != 0 {
		observer.Sampler = NewSampler(samplingInterval, observer.makeBackCaller(callback, summarizerWindows))
	}
	if httpPublisher != nil {
		observer.AddReceiver(httpPublisher)
	}
	if log != nil {
		observer.AddReceiver(newLogger(log))
	}
	return
}

func (o *CallbackObserver) makeBackCaller(callback func(time.Time) []*Observation, summarizerWindows []time.Duration) func(t time.Time) {
	// TODO(go1.1)
	// In Go 1.1 we apparently will be allowed to pass around methods
	// just like we pass around funcs. Then we can move the state to
	// the CallbackObserver itself and just pass this method to the
	// Sampler. But until that, I'll just use a closure.
	summarizers := make(map[string][]*SlidingWindowSummarizer)
	return func(t time.Time) {
		observations := callback(t)
		for _, obs := range observations {
			o.broadcast(obs)
			summers, summersCreated := summarizers[obs.Name]
			if !summersCreated {
				for _, window := range summarizerWindows {
					summers = append(summers, NewSlidingWindowSummarizer(obs.Name, window))
				}
				summarizers[obs.Name] = summers
			}
			for _, summer := range summers {
				summer.Update(obs)
				for _, sum := range summer.Summarize() {
					o.broadcast(sum)
				}
			}
		}
	}
}
