package apps

import "testing"

func TestParseBuilds(t *testing.T) {
	out := `[{"id":"a1","kind":"build","status":"succeeded","started_at":"2026-07-24T07:43:30Z","exit_code":0},
	{"id":"b2","kind":"deploy","status":"failed","started_at":"2026-10-26T02:27:49Z","exit_code":1}]`
	b, err := parseBuilds(out)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 2 || b[0].ID != "b2" {
		t.Fatalf("want newest first, got %+v", b)
	}
	if empty, err := parseBuilds("  \n"); err != nil || len(empty) != 0 {
		t.Fatalf("empty output: %+v %v", empty, err)
	}
	if _, err := parseBuilds("not json"); err == nil {
		t.Fatal("expected error for bad json")
	}
}
