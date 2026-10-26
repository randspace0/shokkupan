package dokku

import (
	"testing"
	"time"
)

func TestFreshWindowIsShorterThanMonitorPoll(t *testing.T) {
	// the monitor polls every minute; a report older than that must refresh
	if psFresh >= time.Minute {
		t.Fatalf("psFresh = %v, want < 1m", psFresh)
	}
}
