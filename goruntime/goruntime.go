// Telemetry for the Go runtime package. Exposes NumCPU, NumGoroutine,
// NumCgoCall as well as the output from ReadMemStats. Default
// sampling interval is 5 seconds. Data is pushed to the
// DefaultHTTPPublisher and no summarizers are used.

// To use, you can simply import the package:
//
//   import _ gotelem/goruntime
//

// You can reconfigure the defaults by creating a factory and passing
// it to InitFromFactory method. DefaultFactory is the factory used
// when you import the goruntime package.

package goruntime

import (
	telem "gotelem"
	"runtime"
	"time"
)

func init() {
	InitFromFactory(DefaultFactory)
}

var DefaultFactory *telem.Factory = &telem.Factory{SamplingInterval: 5 * time.Second, HTTPPublisher: telem.DefaultHTTPPublisher}
var currentObserver *telem.CallbackObserver

// Stops the running sampler and starts a new by calling
// NewCallbackObserver on the factory.
func InitFromFactory(factory *telem.Factory) {
	if currentObserver != nil {
		currentObserver.Stop()
	}
	currentObserver = factory.NewCallbackObserver(sample)
}

func sample(t time.Time) []*telem.Observation {
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)
	return []*telem.Observation{
		&telem.Observation{t, "Goruntime_NumGoroutine", float64(runtime.NumGoroutine())},
		&telem.Observation{t, "Goruntime_NumCgoCall", float64(runtime.NumCgoCall())},
		&telem.Observation{t, "Goruntime_MemAlloc", float64(m.Alloc)},
		&telem.Observation{t, "Goruntime_MemTotalAlloc", float64(m.TotalAlloc)},
		&telem.Observation{t, "Goruntime_MemSys", float64(m.Sys)},
		&telem.Observation{t, "Goruntime_MemLookups", float64(m.Lookups)},
		&telem.Observation{t, "Goruntime_MemMallocs", float64(m.Mallocs)},
		&telem.Observation{t, "Goruntime_MemHeapAlloc", float64(m.HeapAlloc)},
		&telem.Observation{t, "Goruntime_MemHeapSys", float64(m.HeapSys)},
		&telem.Observation{t, "Goruntime_MemHeapIdle", float64(m.HeapIdle)},
		&telem.Observation{t, "Goruntime_MemHeapInuse", float64(m.HeapInuse)},
		&telem.Observation{t, "Goruntime_MemHeapReleased", float64(m.HeapReleased)},
		&telem.Observation{t, "Goruntime_MemHeapObjects", float64(m.HeapObjects)},
		&telem.Observation{t, "Goruntime_MemNextGC", float64(m.NextGC)},
		&telem.Observation{t, "Goruntime_MemLastGC", float64(m.LastGC)},
		&telem.Observation{t, "Goruntime_MemPauseTotalNs", float64(m.PauseTotalNs)},
		&telem.Observation{t, "Goruntime_MemNumGC", float64(m.NumGC)},
	}
}
