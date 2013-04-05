package gotelem

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"testing"
	"time"
)

// Creates a time.Now equivalent that increments it's internal time by
// `increment` on each invocation.
func makeMockTimeNow(initial time.Time, increment time.Duration) func() time.Time {
	t := initial
	return func() time.Time {
		t = t.Add(increment)
		return t.UTC()
	}
}

func TestSmokeTestObserver(t *testing.T) {
	rand.Seed(1337)
	factory := &Factory{
		Logger:            os.Stdout,
		SamplingInterval:  time.Second,
		SummarizerWindows: []time.Duration{time.Minute},
		HTTPPublisher:     DefaultHTTPPublisher}
	observer := factory.NewObserver("BAPI_Schedule_ExecTime")
	observer.timeNow = makeMockTimeNow(time.Now().UTC(), 100*time.Millisecond)

	// NewObserver must starts a sampler goroutine, give it a chance to start before we continue
	runtime.Gosched()

	// Now create the fake ticker chan and insert it into the observer
	fakeTickerChan := make(chan time.Time)
	observer.setSamplingTicker(&time.Ticker{C: fakeTickerChan})

	counter := factory.NewCounter("BAPI_Schedule_Calls")
	fakeTickerChanCounter := make(chan time.Time)
	counter.setSamplingTicker(&time.Ticker{C: fakeTickerChanCounter})

	for i := 1; i <= 100; i++ {
		observer.Observe(rand.Float64() * 100)
		counter.Inc()
		if i%10 == 0 {
			fakeTickerChan <- observer.timeNow()
			fakeTickerChanCounter <- observer.timeNow()
			runtime.Gosched()
		}
	}
	// Make the publisher drain it's inbox
	runtime.Gosched()
	p := DefaultHTTPPublisher
	fmt.Println("series=", len(p.series))
	for name, s := range p.series {
		fmt.Printf("  %v: %v\n", name, len(s.store[s.oldestAt:]))
		for _, o := range s.store[s.oldestAt:] {
			fmt.Println("  ", o.Timestamp, o.Value)
		}
	}
}
