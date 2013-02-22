package gotelem

import (
	"fmt"
	"os"
)

type receiver interface {
	receiverChannel() chan<- *Observation
}

type broadcaster []chan<- *Observation

func (b *broadcaster) AddReceiver(r receiver) {
	if c := r.receiverChannel(); c != nil {
		*b = append(*b, c)
	} else {
		fmt.Fprintln(os.Stderr, "WARN: nil value from receiverChannel() from", r)
	}
}

func (b broadcaster) broadcast(o *Observation) {
	for _, c := range b {
		c <- o
	}
}
