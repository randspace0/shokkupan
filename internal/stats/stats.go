// Package stats reads container usage from the dokku-stats plugin
// (https://github.com/randspace0/dokku-stats).
package stats

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrNotInstalled = errors.New("dokku-stats plugin is not installed")

type Sample struct {
	App       string
	Process   string
	Container string
	CPU       float64 // percent of one core, can exceed 100
	MemBytes  int64
	MemLimit  int64
}

type outputError interface{ Output() string }

// Report runs stats:report for app, or all apps when app is empty.
func Report(exec func(string) (string, error), app string) ([]Sample, error) {
	cmd := "stats:report"
	if app != "" {
		cmd += " " + app
	}
	out, err := exec(cmd)
	if err != nil {
		var oe outputError
		if errors.As(err, &oe) && strings.Contains(oe.Output(), "is not a dokku command") {
			return nil, ErrNotInstalled
		}
		return nil, err
	}
	return Parse(out)
}

// Parse reads lines of: app, process, container, cpu%, "used / limit" (tab separated).
func Parse(out string) ([]Sample, error) {
	var samples []Sample
	for _, line := range strings.Split(out, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		cols := strings.Split(line, "\t")
		if len(cols) != 5 {
			return nil, fmt.Errorf("unexpected stats line: %q", line)
		}
		cpu, err := strconv.ParseFloat(strings.TrimSuffix(strings.TrimSpace(cols[3]), "%"), 64)
		if err != nil {
			return nil, fmt.Errorf("parsing cpu %q: %w", cols[3], err)
		}
		used, limit, ok := strings.Cut(cols[4], "/")
		if !ok {
			return nil, fmt.Errorf("parsing memory %q", cols[4])
		}
		usedB, err := parseBytes(used)
		if err != nil {
			return nil, err
		}
		limitB, err := parseBytes(limit)
		if err != nil {
			return nil, err
		}
		samples = append(samples, Sample{
			App: cols[0], Process: cols[1], Container: cols[2],
			CPU: cpu, MemBytes: usedB, MemLimit: limitB,
		})
	}
	return samples, nil
}

var units = []struct {
	suffix string
	mult   float64
}{
	{"TiB", 1 << 40}, {"GiB", 1 << 30}, {"MiB", 1 << 20}, {"KiB", 1 << 10},
	{"TB", 1e12}, {"GB", 1e9}, {"MB", 1e6}, {"kB", 1e3}, {"B", 1},
}

func parseBytes(s string) (int64, error) {
	s = strings.TrimSpace(s)
	for _, u := range units {
		if num, ok := strings.CutSuffix(s, u.suffix); ok {
			v, err := strconv.ParseFloat(strings.TrimSpace(num), 64)
			if err != nil {
				return 0, fmt.Errorf("parsing bytes %q: %w", s, err)
			}
			return int64(v * u.mult), nil
		}
	}
	return 0, fmt.Errorf("unknown size %q", s)
}
