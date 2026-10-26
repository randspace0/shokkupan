package resources

import (
	"fmt"
	"strconv"
	"strings"

	dokkugo "github.com/texm/dokku-go"
)

// The dokku-go resource report parser predates dokku 0.3x: it expects
// "<proc> limit cpu" keys (dokku now prints "resource <proc> limit cpu") and
// integer amounts (cpu limits can be fractional, e.g. 0.5). These types keep
// the same JSON shape with float amounts.

type Resource struct {
	Type   dokkugo.ResourceSpec `json:"type"`
	Amount float64              `json:"amount"`
}

type ResourceUnits struct {
	CPU            *Resource `json:"cpu"`
	Memory         *Resource `json:"memory"`
	MemorySwap     *Resource `json:"memory_swap"`
	Network        *Resource `json:"network"`
	NetworkIngress *Resource `json:"network_ingress"`
	NetworkEgress  *Resource `json:"network_egress"`
	NvidiaGPU      *Resource `json:"nvidia_gpu"`
}

type ResourceSettings struct {
	Limits       ResourceUnits `json:"limits"`
	Reservations ResourceUnits `json:"reservations"`
}

type AppResourceReport struct {
	Defaults  ResourceSettings            `json:"defaults"`
	Processes map[string]ResourceSettings `json:"processes"`
}

const defaultProcess = "_default_"

func parseResource(kind, raw string) (*Resource, error) {
	spec := dokkugo.ResourceSpec{Name: kind}
	if kind == "memory" || kind == "memory-swap" {
		spec.Suffix = "m" // dokku's default memory unit is MB
		if last := raw[len(raw)-1:]; strings.Contains("bkmg", last) {
			spec.Suffix, raw = last, raw[:len(raw)-1]
		}
	}
	amount, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return nil, fmt.Errorf("resource amount %q: %w", raw, err)
	}
	return &Resource{Type: spec, Amount: amount}, nil
}

func (u *ResourceUnits) set(kind string, r *Resource) {
	switch kind {
	case "cpu":
		u.CPU = r
	case "memory":
		u.Memory = r
	case "memory-swap":
		u.MemorySwap = r
	case "network":
		u.Network = r
	case "network-ingress":
		u.NetworkIngress = r
	case "network-egress":
		u.NetworkEgress = r
	case "nvidia-gpu":
		u.NvidiaGPU = r
	}
	// unknown resource kinds from newer dokku are ignored, not fatal
}

// ParseResourceReport parses `dokku resource:report <app>`; lines look like
// "resource worker limit cpu:     0.5" (older dokku omits the "resource" word).
func ParseResourceReport(output string) (*AppResourceReport, error) {
	report := &AppResourceReport{Processes: map[string]ResourceSettings{}}
	for _, line := range strings.Split(output, "\n") {
		key, val, ok := strings.Cut(line, ":")
		val = strings.TrimSpace(val)
		if !ok || val == "" || strings.HasPrefix(strings.TrimSpace(line), "=====>") {
			continue
		}
		fields := strings.Fields(key)
		if len(fields) > 0 && fields[0] == "resource" {
			fields = fields[1:]
		}
		if len(fields) != 3 {
			continue
		}
		proc, mode, kind := fields[0], fields[1], fields[2]

		res, err := parseResource(kind, val)
		if err != nil {
			return nil, err
		}

		settings := report.Defaults
		if proc != defaultProcess {
			settings = report.Processes[proc]
		}
		switch mode {
		case "limit":
			settings.Limits.set(kind, res)
		case "reserve":
			settings.Reservations.set(kind, res)
		default:
			return nil, fmt.Errorf("unknown resource management %q", mode)
		}
		if proc == defaultProcess {
			report.Defaults = settings
		} else {
			report.Processes[proc] = settings
		}
	}
	return report, nil
}
