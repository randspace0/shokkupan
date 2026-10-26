// Package sysmon samples host CPU, memory and network while the
// "server monitoring" switch is on, and logs a warning on spikes.
//
// Shokku runs in a container: /proc/stat and /proc/meminfo report the host,
// but /proc/net/dev only covers the container's network namespace.
package sysmon

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/randspace0/shokkupan/internal/discord"
	"github.com/randspace0/shokkupan/internal/env"
	"github.com/randspace0/shokkupan/internal/models"
	"github.com/rs/zerolog/log"
)

const (
	interval  = 15 * time.Second
	retention = 7 * 24 * time.Hour

	cpuSpike = 90.0 // percent
	memSpike = 90.0 // percent
	netSpike = 50e6 // bytes/second, rx or tx

	colorSpike = 0xE67E22
)

type cpuTimes struct{ idle, total uint64 }

// parseCPU reads the aggregate "cpu" line of /proc/stat.
func parseCPU(stat string) (cpuTimes, error) {
	for _, line := range strings.Split(stat, "\n") {
		f := strings.Fields(line)
		if len(f) < 5 || f[0] != "cpu" {
			continue
		}
		var t cpuTimes
		for i, v := range f[1:] {
			n, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return cpuTimes{}, err
			}
			t.total += n
			if i == 3 || i == 4 { // idle, iowait
				t.idle += n
			}
		}
		return t, nil
	}
	return cpuTimes{}, fmt.Errorf("no cpu line in /proc/stat")
}

func cpuPercent(prev, cur cpuTimes) float64 {
	total := cur.total - prev.total
	if cur.total < prev.total || total == 0 {
		return 0
	}
	return 100 * float64(total-(cur.idle-prev.idle)) / float64(total)
}

// parseMemPercent returns used memory (total - available) as a percentage.
func parseMemPercent(meminfo string) (float64, error) {
	var total, avail float64
	for _, line := range strings.Split(meminfo, "\n") {
		f := strings.Fields(line)
		if len(f) < 2 {
			continue
		}
		v, _ := strconv.ParseFloat(f[1], 64)
		switch f[0] {
		case "MemTotal:":
			total = v
		case "MemAvailable:":
			avail = v
		}
	}
	if total == 0 {
		return 0, fmt.Errorf("no MemTotal in /proc/meminfo")
	}
	return 100 * (total - avail) / total, nil
}

// parseNet sums rx/tx bytes over all interfaces except loopback.
func parseNet(dev string) (rx, tx uint64) {
	for _, line := range strings.Split(dev, "\n") {
		name, rest, ok := strings.Cut(line, ":")
		if !ok || strings.TrimSpace(name) == "lo" {
			continue
		}
		f := strings.Fields(rest)
		if len(f) < 9 {
			continue
		}
		r, _ := strconv.ParseUint(f[0], 10, 64)
		t, _ := strconv.ParseUint(f[8], 10, 64)
		rx, tx = rx+r, tx+t
	}
	return rx, tx
}

func rate(prev, cur uint64, d time.Duration) float64 {
	if cur < prev || d <= 0 {
		return 0
	}
	return float64(cur-prev) / d.Seconds()
}

func read(path string) string {
	b, _ := os.ReadFile(path)
	return string(b)
}

// Enabled reports whether the switch is on.
func Enabled(e *env.Env) bool {
	var s models.Server
	if err := e.DB.First(&s).Error; err != nil {
		return false
	}
	return s.ServerMonitoring
}

// spike is one threshold crossing in a sample.
type spike struct {
	Kind    string // "cpu", "memory" or "network"
	Message string
}

func detectSpikes(s models.ServerSample) []spike {
	var out []spike
	if s.CPUPercent >= cpuSpike {
		out = append(out, spike{"cpu", fmt.Sprintf("CPU at %.0f%% (threshold %.0f%%)", s.CPUPercent, cpuSpike)})
	}
	if s.MemPercent >= memSpike {
		out = append(out, spike{"memory", fmt.Sprintf("Memory at %.0f%% (threshold %.0f%%)", s.MemPercent, memSpike)})
	}
	if s.NetRxBps >= netSpike || s.NetTxBps >= netSpike {
		out = append(out, spike{"network", fmt.Sprintf("Network in %.1f MB/s, out %.1f MB/s (threshold %.0f MB/s)",
			s.NetRxBps/1e6, s.NetTxBps/1e6, netSpike/1e6)})
	}
	return out
}

// alertCooldown stops a sustained spike from posting every 15 seconds.
const alertCooldown = 15 * time.Minute

type cooldown struct{ last map[string]time.Time }

// due reports whether kind may alert now, and records it if so.
func (c *cooldown) due(kind string, now time.Time) bool {
	if c.last == nil {
		c.last = map[string]time.Time{}
	}
	if t, ok := c.last[kind]; ok && now.Sub(t) < alertCooldown {
		return false
	}
	c.last[kind] = now
	return true
}

// alertSpikes logs every spike and posts the ones past their cooldown to the
// global Discord webhook (a no-op when the plugin or webhook is missing).
func alertSpikes(e *env.Env, cd *cooldown, spikes []spike, now time.Time) {
	for _, sp := range spikes {
		log.Warn().Str("kind", sp.Kind).Msg("sysmon: " + sp.Message)
		if !cd.due(sp.Kind, now) {
			continue
		}
		url, err := discord.Get(e.Dokku.Exec, "")
		if err != nil {
			if !errors.Is(err, discord.ErrNotConfigured) && !errors.Is(err, discord.ErrNotInstalled) {
				log.Error().Err(err).Msg("sysmon: failed to get discord webhook")
			}
			return
		}
		if err := discord.Notify(url, "Server "+sp.Kind+" spike", sp.Message, colorSpike); err != nil {
			log.Error().Err(err).Msg("sysmon: failed to send alert")
		}
	}
}

// Run samples forever; it does nothing while the switch is off.
func Run(e *env.Env) {
	var (
		prevCPU  cpuTimes
		prevRx   uint64
		prevTx   uint64
		prevTime time.Time
		primed   bool
		cd       cooldown
	)
	for ; ; time.Sleep(interval) {
		if !Enabled(e) {
			primed = false
			continue
		}

		now := time.Now()
		cpu, err := parseCPU(read("/proc/stat"))
		mem, memErr := parseMemPercent(read("/proc/meminfo"))
		if err != nil || memErr != nil {
			log.Error().Err(err).AnErr("mem", memErr).Msg("sysmon: failed to read /proc")
			continue
		}
		rx, tx := parseNet(read("/proc/net/dev"))

		if primed {
			s := models.ServerSample{
				CPUPercent: cpuPercent(prevCPU, cpu),
				MemPercent: mem,
				NetRxBps:   rate(prevRx, rx, now.Sub(prevTime)),
				NetTxBps:   rate(prevTx, tx, now.Sub(prevTime)),
			}
			if err := e.DB.Create(&s).Error; err != nil {
				log.Error().Err(err).Msg("sysmon: failed to record sample")
			}
			alertSpikes(e, &cd, detectSpikes(s), now)
			e.DB.Unscoped().Where("created_at < ?", now.Add(-retention)).Delete(&models.ServerSample{})
		}
		prevCPU, prevRx, prevTx, prevTime, primed = cpu, rx, tx, now, true
	}
}
