// Package monitor watches whether deployed apps are running, records
// changes and alerts through the dokku-discord webhook.
package monitor

import (
	"errors"
	"github.com/randspace0/shokkupan/internal/discord"
	"github.com/randspace0/shokkupan/internal/env"
	"github.com/randspace0/shokkupan/internal/models"
	"github.com/randspace0/shokkupan/internal/server/dokku"
	"github.com/randspace0/shokkupan/internal/stats"
	"github.com/rs/zerolog/log"
	dokkugo "github.com/texm/dokku-go"
	"time"
)

const (
	pollInterval   = time.Minute
	retention      = 30 * 24 * time.Hour
	statsRetention = 7 * 24 * time.Hour

	colorDown = 0xE74C3C
	colorUp   = 0x2ECC71
)

type change struct {
	App     string
	Running bool
}

// transitions returns apps whose running state changed. Apps absent from
// prev (first sight, or the first poll) are not reported.
func transitions(prev, cur map[string]bool) []change {
	var out []change
	for app, running := range cur {
		if was, seen := prev[app]; seen && was != running {
			out = append(out, change{App: app, Running: running})
		}
	}
	return out
}

func snapshot(e *env.Env) (map[string]bool, error) {
	reports, err := dokku.GetProcessReports(e)
	if err != nil {
		return nil, err
	}
	cur := map[string]bool{}
	for name, r := range reports {
		if name != dokku.FilteredApp && r.Deployed {
			cur[name] = r.Running
		}
	}
	return cur, nil
}

func alert(e *env.Env, c change) {
	url, err := discord.Get(e.Dokku.Exec, c.App)
	if err != nil {
		if !errors.Is(err, discord.ErrNotConfigured) && !errors.Is(err, discord.ErrNotInstalled) {
			log.Error().Err(err).Str("app", c.App).Msg("monitor: failed to get discord webhook")
		}
		return
	}

	title, desc, color := c.App+" is down", "Not running (stopped or crashed).", colorDown
	if c.Running {
		title, desc, color = c.App+" is running", "Back up.", colorUp
	}
	if err := discord.Notify(url, title, desc, color); err != nil {
		log.Error().Err(err).Str("app", c.App).Msg("monitor: failed to send alert")
	}
}

// collectStats stores one summed sample per app; it is a no-op until the
// dokku-stats plugin is installed.
func collectStats(e *env.Env) {
	samples, err := stats.Report(e.Dokku.Exec, "")
	if err != nil {
		if !errors.Is(err, stats.ErrNotInstalled) {
			log.Error().Err(err).Msg("monitor: failed to get stats")
		}
		return
	}

	perApp := map[string]*models.AppStatsSample{}
	for _, s := range samples {
		row, ok := perApp[s.App]
		if !ok {
			row = &models.AppStatsSample{AppName: s.App}
			perApp[s.App] = row
		}
		row.CPU += s.CPU
		row.MemBytes += s.MemBytes
		row.MemLimit += s.MemLimit
	}
	for _, row := range perApp {
		if err := e.DB.Create(row).Error; err != nil {
			log.Error().Err(err).Msg("monitor: failed to record stats")
		}
	}
	e.DB.Unscoped().Where("created_at < ?", time.Now().Add(-statsRetention)).Delete(&models.AppStatsSample{})
}

func poll(e *env.Env, prev map[string]bool) map[string]bool {
	cur, err := snapshot(e)
	if err != nil {
		if !errors.Is(err, dokkugo.NoDeployedAppsError) {
			log.Error().Err(err).Msg("monitor: failed to get process report")
			return prev
		}
		cur = map[string]bool{}
	}

	if prev != nil {
		for _, c := range transitions(prev, cur) {
			if err := e.DB.Create(&models.AppStatusEvent{AppName: c.App, Running: c.Running}).Error; err != nil {
				log.Error().Err(err).Msg("monitor: failed to record event")
			}
			alert(e, c)
		}
	}

	e.DB.Unscoped().Where("created_at < ?", time.Now().Add(-retention)).Delete(&models.AppStatusEvent{})
	collectStats(e)
	return cur
}

// Run polls forever; the first poll only seeds state so a restart of
// shokkupan does not alert on every app.
func Run(e *env.Env) {
	var state map[string]bool
	for {
		state = poll(e, state)
		time.Sleep(pollInterval)
	}
}
