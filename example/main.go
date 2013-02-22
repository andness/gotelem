package main

import (
	telem "gotelem"
	_ "gotelem/goruntime"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var factory *telem.Factory = &telem.Factory{
	Logger:            log.Println,
	SamplingInterval:  1 * time.Second,
	SummarizerWindows: []time.Duration{time.Minute, 5 * time.Minute},
	HTTPPublisher:     telem.DefaultHTTPPublisher}

func main() {
	go makeSomeObservations()
	telem.DefaultHTTPPublisher.SetBaseURL("http://localhost:8888")
	http.Handle("/", telem.DefaultHTTPPublisher)
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func makeSomeObservations() {
	observer := factory.NewObserver("BAPI_Schedule_ExecTime")
	counter := factory.NewCounter("BAPI_Schedule_CallCount")
	for {
		observer.Observe(rand.Float64() * 300)
		counter.Inc()
		time.Sleep(300 * time.Millisecond)
	}
}
