package resources

import "testing"

func TestParseResourceReport(t *testing.T) {
	out := `=====> opentransum-backend resource information
       resource _default_ limit cpu:  1
       resource worker limit cpu:     0.5
       resource worker limit memory:  512m
       resource worker reserve memory: 1g
       resource web limit nvidia-gpu: 1
       resource web limit future-thing: 7
`
	r, err := ParseResourceReport(out)
	if err != nil {
		t.Fatal(err)
	}
	if r.Defaults.Limits.CPU == nil || r.Defaults.Limits.CPU.Amount != 1 {
		t.Fatalf("default cpu = %+v", r.Defaults.Limits.CPU)
	}
	w := r.Processes["worker"]
	if w.Limits.CPU.Amount != 0.5 {
		t.Fatalf("worker cpu = %v, want 0.5", w.Limits.CPU.Amount)
	}
	if m := w.Limits.Memory; m.Amount != 512 || m.Type.Suffix != "m" {
		t.Fatalf("worker memory limit = %+v", m)
	}
	if m := w.Reservations.Memory; m.Amount != 1 || m.Type.Suffix != "g" {
		t.Fatalf("worker memory reserve = %+v", m)
	}
	if r.Processes["web"].Limits.NvidiaGPU == nil {
		t.Fatal("web gpu missing")
	}
}

func TestParseResourceReportLegacyKeysAndEmpty(t *testing.T) {
	r, err := ParseResourceReport("       _default_ limit cpu: 2\n")
	if err != nil || r.Defaults.Limits.CPU.Amount != 2 {
		t.Fatalf("legacy: %+v %v", r, err)
	}
	if r, err = ParseResourceReport("=====> app resource information\n"); err != nil || len(r.Processes) != 0 {
		t.Fatalf("empty: %+v %v", r, err)
	}
}
