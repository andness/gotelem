package gotelem

import (
	"fmt"
)

func newLogger(logFunc func(v ...interface{})) (l *logger) {
	l = &logger{inbox: make(chan *Observation, 128), logFunc: logFunc}
	go l.process()
	return
}

type logger struct {
	inbox   chan *Observation
	logFunc func(v ...interface{})
}

func (l *logger) receiverChannel() chan<- *Observation {
	return l.inbox
}

func (l *logger) process() {
	for {
		o := <-l.inbox
		l.logFunc(fmt.Sprintf("%v,%v,%v", o.Timestamp.UnixNano(), o.Name, o.Value))
	}
}
