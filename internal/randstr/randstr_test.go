package randstr

import "testing"

func TestStringLengthAndAlphabet(t *testing.T) {
	s := String(64)
	if len(s) != 64 {
		t.Fatalf("len = %d", len(s))
	}
	for _, c := range s {
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z') {
			t.Fatalf("bad char %q", c)
		}
	}
}
