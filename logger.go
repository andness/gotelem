package gotelem

import (
	"fmt"
	"io"
)

func newLogger(sink io.Writer) (l *logger) {
	l = &logger{inbox: make(chan *Observation, 128), sink: sink}
	go l.process()
	return
}

type logger struct {
	inbox chan *Observation
	sink  io.Writer
}

func (l *logger) receiverChannel() chan<- *Observation {
	return l.inbox
}

func (l *logger) process() {
	for {
		o := <-l.inbox
		fmt.Fprintf(l.sink, "%v,%v,%v\n", o.Timestamp.UnixNano(), o.Name, o.Value)
	}
}
