package sysmon

import (
	"testing"
	"time"

	"github.com/randspace0/shokkupan/internal/models"
)

func TestParsers(t *testing.T) {
	a, _ := parseCPU("cpu  100 0 100 700 100 0 0 0\ncpu0 1 1 1 1 1")
	b, _ := parseCPU("cpu  200 0 200 750 100 0 0 0")
	// total +250, idle +50 -> 80% busy
	if got := cpuPercent(a, b); got != 80 {
		t.Fatalf("cpu = %v, want 80", got)
	}

	mem, _ := parseMemPercent("MemTotal: 1000 kB\nMemAvailable: 250 kB\n")
	if mem != 75 {
		t.Fatalf("mem = %v, want 75", mem)
	}

	rx, tx := parseNet("Inter-|   Receive\n face |bytes\n    lo: 999 0 0 0 0 0 0 0 999 0\n  eth0: 10 0 0 0 0 0 0 0 20 0\n")
	if rx != 10 || tx != 20 {
		t.Fatalf("net = %d/%d, want 10/20", rx, tx)
	}
}

func TestDetectSpikes(t *testing.T) {
	if got := detectSpikes(models.ServerSample{CPUPercent: 50, MemPercent: 50}); len(got) != 0 {
		t.Fatalf("quiet sample produced %v", got)
	}
	got := detectSpikes(models.ServerSample{CPUPercent: 95, MemPercent: 91, NetRxBps: 60e6})
	kinds := map[string]bool{}
	for _, sp := range got {
		kinds[sp.Kind] = true
	}
	if !kinds["cpu"] || !kinds["memory"] || !kinds["network"] {
		t.Fatalf("want cpu+memory+network, got %v", got)
	}
}

func TestCooldown(t *testing.T) {
	var cd cooldown
	now := time.Now()
	if !cd.due("cpu", now) {
		t.Fatal("first alert must fire")
	}
	if cd.due("cpu", now.Add(time.Minute)) {
		t.Fatal("second alert inside cooldown must be suppressed")
	}
	if !cd.due("memory", now.Add(time.Minute)) {
		t.Fatal("other kinds have their own cooldown")
	}
	if !cd.due("cpu", now.Add(alertCooldown+time.Second)) {
		t.Fatal("alert must fire again after cooldown")
	}
}
