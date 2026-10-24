package discord

import (
	"errors"
	"testing"
)

type outErr struct{ out string }

func (e outErr) Error() string  { return e.out }
func (e outErr) Output() string { return e.out }

const goodURL = "https://discord.com/api/webhooks/123456/abc-DEF_9"

func TestSetRejectsInjection(t *testing.T) {
	called := false
	exec := func(string) (string, error) { called = true; return "", nil }
	for _, u := range []string{"", "http://discord.com/api/webhooks/1/a", goodURL + "; rm -rf /", "https://evil.com/api/webhooks/1/a"} {
		if err := Set(exec, u, ""); !errors.Is(err, ErrInvalidWebhook) {
			t.Errorf("%q: got %v", u, err)
		}
	}
	if called {
		t.Fatal("exec must not run for invalid urls")
	}
}

func TestSetCommand(t *testing.T) {
	var got string
	exec := func(c string) (string, error) { got = c; return "", nil }
	if err := Set(exec, goodURL, "myapp"); err != nil {
		t.Fatal(err)
	}
	if want := "discord:set " + goodURL + " myapp"; got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestGetClassifiesErrors(t *testing.T) {
	cases := map[string]error{
		"!     No Discord webhook URL configured globally": ErrNotConfigured,
		"!     `discord:get` is not a dokku command.":      ErrNotInstalled,
	}
	for out, want := range cases {
		_, err := Get(func(string) (string, error) { return "", outErr{out} }, "")
		if !errors.Is(err, want) {
			t.Errorf("%q: got %v", out, err)
		}
	}
}

func TestMask(t *testing.T) {
	if got := Mask(goodURL); got != "https://discord.com/api/webhooks/123456/****" {
		t.Fatal(got)
	}
}
