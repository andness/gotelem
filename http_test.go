package gotelem

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestHTTPPublisher(t *testing.T) {
	// Typical case of a publisher keeping 1 hour of per second data.
	p := NewHTTPPublisher(3600)
	ts := time.Date(1978, 2, 12, 16, 0, 0, 0, time.UTC)
	inbox := p.receiverChannel()
	for i := 0; i < 4; i++ {
		name := fmt.Sprintf("Test%v", i)
		for j := 0; j < 5000; j++ {
			inbox <- &Observation{ts.Add(time.Duration(j) * time.Second), name, float64(j)}
		}
	}
	// Make the publisher drain it's inbox
	runtime.Gosched()

	fmt.Println("num series=", len(p.series))
	for name, s := range p.series {
		fmt.Printf("  %v: %v\n observations", name, len(s.store[s.oldestAt:]))
		fmt.Printf("    First=%v\n", s.store[s.oldestAt])
		fmt.Printf("    Last=%v\n", s.store[len(s.store)-1])
		// for _, o := range s.store[s.oldestAt:] {
		// 	fmt.Println("  ", o.Timestamp, o.Value)
		// }
	}
}

func _BenchmarkRingBuffer(b *testing.B, keep int, capacity int) {
	rb := newObservationRingBufferWithCapacity(keep, float64(capacity)/float64(keep))
	obs := &Observation{time.Now(), "Test", 666.6}
	//ms := &runtime.MemStats{}
	//runtime.ReadMemStats(ms)
	//before := ms.Alloc
	for i := 0; i < b.N; i++ {
		rb.update(obs)
	}
	//runtime.ReadMemStats(ms)
	//after := ms.Alloc
	//fmt.Printf("Mem: %v -> %v Δ=%v\n", before, after, after-before)
}

func MakeBenchmarkFunc(keep, capacity int) func(b *testing.B) {
	return func(b *testing.B) { _BenchmarkRingBuffer(b, keep, capacity) }
}

func TestBenchIt(t *testing.T) {
	for keep := 100; keep <= 1000000; keep *= 10 {
		for c := 10; c > 0; c-- {
			capacity := keep + (keep/10)*c
			fn := MakeBenchmarkFunc(keep, capacity)
			res := testing.Benchmark(fn)
			fmt.Printf("keep=%v cap=%v overhead=%v: %s\n", keep, capacity, float64(capacity)/float64(keep), res)
		}
	}
}

/*

Results from TestBenchIt. Seems like 20% or 30% overhead is sufficient
to give good perf in most cases.

=== RUN TestBenchIt
keep=100 cap=200 overhead=2: 200000000	         8.55 ns/op
keep=100 cap=190 overhead=1.9: 200000000	         8.67 ns/op
keep=100 cap=180 overhead=1.8: 200000000	         8.79 ns/op
keep=100 cap=170 overhead=1.7: 200000000	         8.93 ns/op
keep=100 cap=160 overhead=1.6: 200000000	         9.23 ns/op
keep=100 cap=150 overhead=1.5: 200000000	         9.39 ns/op
keep=100 cap=140 overhead=1.4: 200000000	         9.88 ns/op
keep=100 cap=130 overhead=1.3: 100000000	        10.7 ns/op
keep=100 cap=120 overhead=1.2: 100000000	        12.2 ns/op
keep=100 cap=110 overhead=1.1: 100000000	        15.1 ns/op
keep=1000 cap=2000 overhead=2: 200000000	         8.13 ns/op
keep=1000 cap=1900 overhead=1.9: 500000000	         7.89 ns/op
keep=1000 cap=1800 overhead=1.8: 200000000	         8.02 ns/op
keep=1000 cap=1700 overhead=1.7: 200000000	         8.48 ns/op
keep=1000 cap=1600 overhead=1.6: 200000000	         8.57 ns/op
keep=1000 cap=1500 overhead=1.5: 200000000	         8.71 ns/op
keep=1000 cap=1400 overhead=1.4: 200000000	         8.44 ns/op
keep=1000 cap=1300 overhead=1.3: 200000000	         8.56 ns/op
keep=1000 cap=1200 overhead=1.2: 200000000	         9.09 ns/op
keep=1000 cap=1100 overhead=1.1: 100000000	        10.6 ns/op
keep=10000 cap=20000 overhead=2: 200000000	         8.16 ns/op
keep=10000 cap=19000 overhead=1.9: 200000000	         8.29 ns/op
keep=10000 cap=18000 overhead=1.8: 200000000	         8.30 ns/op
keep=10000 cap=17000 overhead=1.7: 200000000	         8.37 ns/op
keep=10000 cap=16000 overhead=1.6: 200000000	         8.49 ns/op
keep=10000 cap=15000 overhead=1.5: 200000000	         8.69 ns/op
keep=10000 cap=14000 overhead=1.4: 200000000	         8.94 ns/op
keep=10000 cap=13000 overhead=1.3: 200000000	         9.34 ns/op
keep=10000 cap=12000 overhead=1.2: 200000000	         9.96 ns/op
keep=10000 cap=11000 overhead=1.1: 100000000	        12.0 ns/op
keep=100000 cap=200000 overhead=2: 200000000	         8.31 ns/op
keep=100000 cap=190000 overhead=1.9: 200000000	         8.41 ns/op
keep=100000 cap=180000 overhead=1.8: 200000000	         8.51 ns/op
keep=100000 cap=170000 overhead=1.7: 200000000	         8.60 ns/op
keep=100000 cap=160000 overhead=1.6: 200000000	         8.77 ns/op
keep=100000 cap=150000 overhead=1.5: 200000000	         8.90 ns/op
keep=100000 cap=140000 overhead=1.4: 200000000	         9.29 ns/op
keep=100000 cap=130000 overhead=1.3: 200000000	         9.78 ns/op
keep=100000 cap=120000 overhead=1.2: 100000000	        10.8 ns/op
keep=100000 cap=110000 overhead=1.1: 100000000	        13.9 ns/op
keep=1000000 cap=2000000 overhead=2: 200000000	         9.10 ns/op
keep=1000000 cap=1900000 overhead=1.9: 200000000	         9.05 ns/op
keep=1000000 cap=1800000 overhead=1.8: 200000000	         9.26 ns/op
keep=1000000 cap=1700000 overhead=1.7: 200000000	         9.50 ns/op
keep=1000000 cap=1600000 overhead=1.6: 200000000	         9.80 ns/op
keep=1000000 cap=1500000 overhead=1.5: 100000000	        10.2 ns/op
keep=1000000 cap=1400000 overhead=1.4: 100000000	        10.9 ns/op
keep=1000000 cap=1300000 overhead=1.3: 100000000	        11.4 ns/op
keep=1000000 cap=1200000 overhead=1.2: 100000000	        13.8 ns/op
keep=1000000 cap=1100000 overhead=1.1: 100000000	        18.7 ns/op
*/
