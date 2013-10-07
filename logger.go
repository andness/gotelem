package gotelem

import (
	"fmt"
	"io"
	"os"
)

func NewLogger(sink io.Writer) (l *Logger) {
	l = &Logger{inbox: make(chan *Observation, 128), sink: sink}
	go l.process()
	return
}

type Logger struct {
	inbox chan *Observation
	sink  io.Writer
}

func (l *Logger) receiverChannel() chan<- *Observation {
	return l.inbox
}

func (l *Logger) process() {
	for {
		o := <-l.inbox
		if _, err := fmt.Fprintf(l.sink, "%v,%v,%f", o.Timestamp.UnixNano(), o.Name, o.Value); err != nil {
			fmt.Fprint(os.Stderr, "logger: Error from Fprintf: %s", err)
		}
		if _, err := fmt.Fprintln(l.sink); err != nil {
			fmt.Fprint(os.Stderr, "logger: Error from Fprintln: %s", err)
		}
	}
}
