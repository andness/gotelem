package gotelem

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
)

var DefaultHTTPPublisher *HTTPPublisher = NewHTTPPublisher(300)

// The HTTP Publisher receives and stores the N latest
// observations. It implements ServeHTTP and will expose the available
// time series as JSON data.
type HTTPPublisher struct {
	baseURL string
	inbox   chan *Observation
	keep    int
	series  map[string]*observationFIFOQueue
}

func (h *HTTPPublisher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	localPath := r.URL.Path
	//println("localPath", localPath)
	w.Header().Set("Content-Type", "application/json")
	if localPath == "/" {
		h.RespondAvailableSeries(w, r)
	} else if localPath == "/series" {
		h.RespondSelectedSeries(w, r)
	}
}

func (h *HTTPPublisher) receiverChannel() chan<- *Observation {
	return h.inbox
}

type timeseries struct {
	Name string
	URL  string
}

func (h *HTTPPublisher) RespondAvailableSeries(w http.ResponseWriter, r *http.Request) {
	series := make([]*timeseries, len(h.series))
	i := 0
	fmt.Println("url=", r.URL.String())
	baseUrl := h.baseURL + "/series?q="
	for k, _ := range h.series {
		series[i] = &timeseries{k, baseUrl + k}
		i++
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(series)
}

type TimeSeries map[string][]TimeSeriesPoint

type TimeSeriesPoint struct {
	Timestamp int64
	Value     float64
}

func (h *HTTPPublisher) RespondSelectedSeries(w http.ResponseWriter, r *http.Request) {
	selected := r.URL.Query()["q"]
	fmt.Println("selected=", selected)
	result := make(map[string][]TimeSeriesPoint)
	for _, name := range selected {
		seriesData := h.series[name]
		if seriesData == nil {
			continue
		}
		observations := seriesData.values()
		points := make([]TimeSeriesPoint, len(observations))
		for i, o := range observations {
			points[i] = TimeSeriesPoint{o.Timestamp.UnixNano(), o.Value}
		}
		result[name] = points
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(result)
}

func NewHTTPPublisher(keep int) *HTTPPublisher {
	publisher := &HTTPPublisher{
		keep:   keep,
		inbox:  make(chan *Observation, 256),
		series: make(map[string]*observationFIFOQueue)}
	go publisher.processInbox()
	return publisher
}

func (h *HTTPPublisher) SetBaseURL(url string) {
	h.baseURL = url
}

func (h *HTTPPublisher) processInbox() {
	for {
		o := <-h.inbox
		rb := h.series[o.Name]
		if rb == nil {
			rb = newObservationRingBuffer(h.keep)
			h.series[o.Name] = rb
		}
		rb.update(o)
	}
}

func newObservationRingBuffer(keep int) *observationFIFOQueue {
	return newObservationRingBufferWithCapacity(keep, 1.3)
}

func newObservationRingBufferWithCapacity(keep int, capacityMultiplier float64) *observationFIFOQueue {
	if capacityMultiplier < 1.0 {
		panic("newObservationRingBuffer: capacityMultiplier cannot be < 1.0")
	}
	capacity := int(math.Ceil(float64(keep) * capacityMultiplier))
	store := make([]*Observation, 0, capacity)
	return &observationFIFOQueue{store: store, keep: keep}
}

// FIFO Queue with max size. Reuses the same slice for the lifetime of
// the queue to avoid generating garbage.
type observationFIFOQueue struct {
	store    []*Observation
	keep     int
	oldestAt int
}

func (rb *observationFIFOQueue) values() (values []*Observation) {
	values = make([]*Observation, len(rb.store)-rb.oldestAt)
	copy(values, rb.store[rb.oldestAt:])
	return
}

func (rb *observationFIFOQueue) update(o *Observation) {
	//println(int(o.Value))
	if len(rb.store) == cap(rb.store) {
		// Buffer is full, move all data to start of buffer and reset
		// oldestAt. Note that we discard the oldest value by slicing
		// from oldestAt + 1.
		//println("enter copy", int(o.Value), rb.oldestAt, len(rb.store), cap(rb.store))
		copy(rb.store, rb.store[rb.oldestAt+1:])
		rb.store = rb.store[0 : rb.keep-1]
		//println("copied", rb.oldestAt, len(rb.store), cap(rb.store))
		rb.oldestAt = 0
	}
	// Discard the oldest value. Note that right after a copy oldestAt
	// will always be 1 less than keep and we already discarded the
	// oldest value in the copy operation.
	if len(rb.store) >= rb.keep {
		rb.oldestAt++
	}
	rb.store = append(rb.store, o)
}
