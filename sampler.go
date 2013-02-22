package gotelem

import (
	"time"
)

type Sampler struct {
	stop     chan bool
	ticker   *time.Ticker
	interval time.Duration
	callback func(time.Time)
}

func NewSampler(interval time.Duration, callback func(time.Time)) (sampler *Sampler) {
	sampler = &Sampler{
		stop:     make(chan bool),
		callback: callback,
		interval: interval}
	sampler.SetInterval(interval)
	return
}

func (s *Sampler) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
		s.stop <- true
	}
}

func (s *Sampler) SetInterval(interval time.Duration) {
	s.interval = interval
	s.setSamplingTicker(time.NewTicker(interval))
}

func (s *Sampler) setSamplingTicker(ticker *time.Ticker) {
	s.Stop()
	s.ticker = ticker
	s.stop = make(chan bool)
	go s.sampler(s.ticker.C, s.stop)
}

func (s *Sampler) sampler(ticker <-chan time.Time, stop <-chan bool) {
	for {
		// ticker.Stop() doesn't close the channel because this could
		// cause problems for users who do not check the return value
		// of the receive. I guess it's a sensible tradeoff, my case
		// of restarting the ticker is probably rare.
		select {
		case t := <-ticker:
			s.callback(t)
		case <-stop:
			return
		}
	}
}
