package monitor

import "testing"

func TestTransitions(t *testing.T) {
	prev := map[string]bool{"a": true, "b": false, "c": true}
	cur := map[string]bool{"a": false, "b": true, "c": true, "new": false}
	got := map[string]bool{}
	for _, c := range transitions(prev, cur) {
		got[c.App] = c.Running
	}
	if len(got) != 2 || got["a"] != false || got["b"] != true {
		t.Fatalf("got %v", got)
	}
}
