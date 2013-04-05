package gotelem

import (
	"io"
	"time"
)

type Factory struct {
	Logger            io.Writer
	SamplingInterval  time.Duration
	SummarizerWindows []time.Duration
	HTTPPublisher     *HTTPPublisher
}

func (f *Factory) NewCounter(name string) (c *Counter) {
	return NewCounter(name, f.SamplingInterval, f.SummarizerWindows, f.HTTPPublisher, f.Logger)
}

func (f *Factory) NewObserver(name string) (o *Observer) {
	return NewObserver(name, f.SamplingInterval, f.SummarizerWindows, f.HTTPPublisher, f.Logger)
}

func (f *Factory) NewCallbackObserver(callback func(time.Time) []*Observation) *CallbackObserver {
	return NewCallbackObserver(callback, f.SamplingInterval, f.SummarizerWindows, f.HTTPPublisher, f.Logger)
}
