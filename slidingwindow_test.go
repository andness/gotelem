package gotelem

import (
	"testing"
	"time"
)

func TestSmokeTestSlidingWindow(t *testing.T) {
	w := NewSlidingWindow(1 * time.Millisecond)
	if len(w.items) != 0 {
		t.Errorf("Window should be empty when nothing has been inserted")
	}
	expired := w.Add(&Observation{time.Now().UTC(), "Test", float64(42)})
	if len(expired) != 0 {
		t.Errorf("No expired items should be returned from empty window")
	}

	// Now make sure the first items expires
	time.Sleep(2 * time.Millisecond)

	expired = w.Add(&Observation{time.Now().UTC(), "Test", float64(84)})
	if len(expired) != 1 {
		t.Errorf("Add should return one expired item")
	} else {
		if expired[0] == nil {
			t.Errorf("Expired item should not be nil")
		} else {
			if expired[0].Value != 42.0 {
				t.Errorf("Expected expired item with Value 84 but got %v", expired[0].Value)
			}
		}
	}
	items := w.items[w.oldestAt:]
	if len(items) != 1 {
		t.Fatalf("Window should contain 1 item, found %v", len(items))
	}
	if items[0] == nil {
		t.Fatalf("Remaining item should not be nil")
	}
}
