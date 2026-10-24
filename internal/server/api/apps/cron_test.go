package apps

import "testing"

func TestParseCronList(t *testing.T) {
	out := `ID                                   Schedule      Command
a1b2c3   */5 * * * *   bundle exec rake sync --verbose
d4e5f6   @daily        node jobs/cleanup.js`
	got := parseCronList(out)
	if len(got) != 2 {
		t.Fatalf("len = %d", len(got))
	}
	if got[0].ID != "a1b2c3" || got[0].Schedule != "*/5 * * * *" || got[0].Command != "bundle exec rake sync --verbose" {
		t.Errorf("task 0: %+v", got[0])
	}
	if got[1].Schedule != "@daily" || got[1].Command != "node jobs/cleanup.js" {
		t.Errorf("task 1: %+v", got[1])
	}
}

func TestParseCronListEmpty(t *testing.T) {
	if got := parseCronList(""); len(got) != 0 {
		t.Fatalf("got %v", got)
	}
}
