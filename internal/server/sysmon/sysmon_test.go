package sysmon

import "testing"

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
