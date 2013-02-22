package gotelem

import (
	"time"
)

type Observation struct {
	Timestamp time.Time
	Name      string
	Value     float64
}
