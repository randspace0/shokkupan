package dokku

import (
	"sync"
	"time"

	"github.com/randspace0/shokkupan/internal/env"
	"github.com/rs/zerolog/log"
	dokkugo "github.com/texm/dokku-go"
)

// `dokku ps:report` inspects every container and takes 10-20s on a busy
// host, so the result is shared: callers get the last report immediately and
// a stale one triggers a single background refresh.
const (
	psFresh = 20 * time.Second
)

type psCache struct {
	mu         sync.Mutex
	report     dokkugo.ProcessReport
	err        error
	at         time.Time
	refreshing bool
	done       chan struct{} // closed when the in-flight refresh ends
}

var processReports psCache

func (c *psCache) refresh(e *env.Env) {
	report, err := e.Dokku.GetAllProcessReport()
	c.mu.Lock()
	if err == nil || c.report == nil {
		c.report, c.err = report, err
	} else {
		log.Error().Err(err).Msg("ps:report refresh failed; serving the previous report")
	}
	c.at = time.Now()
	c.refreshing = false
	close(c.done)
	c.mu.Unlock()
}

// GetProcessReports returns the cached `ps:report` for all apps. The first
// call (or one after a failure with no data) waits for a fresh report.
func GetProcessReports(e *env.Env) (dokkugo.ProcessReport, error) {
	c := &processReports
	c.mu.Lock()
	if !c.refreshing && time.Since(c.at) > psFresh {
		c.refreshing = true
		c.done = make(chan struct{})
		go c.refresh(e)
	}
	if c.report == nil {
		done := c.done
		c.mu.Unlock()
		<-done
		c.mu.Lock()
	}
	report, err := c.report, c.err
	c.mu.Unlock()
	return report, err
}
