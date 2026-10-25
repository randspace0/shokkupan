// Package discord drives the dokku-discord plugin
// (plugins/discord) and posts alerts to its webhook.
package discord

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// Exec runs a dokku command and returns its output, e.g. (*dokku.SSHClient).Exec.
type Exec func(cmd string) (string, error)

var (
	ErrNotInstalled   = errors.New("dokku-discord plugin is not installed")
	ErrNotConfigured  = errors.New("no discord webhook configured")
	ErrInvalidWebhook = errors.New("invalid discord webhook url")

	// strict on purpose: the url is interpolated into a dokku command line
	webhookRe = regexp.MustCompile(`^https://(?:(?:ptb|canary)\.)?discord(?:app)?\.com/api/webhooks/\d+/[A-Za-z0-9_-]+$`)

	httpClient = &http.Client{Timeout: 10 * time.Second}
)

type outputError interface{ Output() string }

// classify maps dokku's textual failures onto sentinel errors.
func classify(err error) error {
	var oe outputError
	if !errors.As(err, &oe) {
		return err
	}
	out := oe.Output()
	switch {
	case strings.Contains(out, "is not a dokku command"):
		return ErrNotInstalled
	case strings.Contains(out, "No Discord webhook URL configured"):
		return ErrNotConfigured
	}
	return err
}

func appArg(app string) string {
	if app == "" {
		return ""
	}
	return " " + app
}

// Get returns the webhook for app, falling back to the global one like the
// plugin does. An empty app returns the global webhook.
func Get(exec Exec, app string) (string, error) {
	out, err := exec("discord:get" + appArg(app))
	if err != nil {
		return "", classify(err)
	}
	return strings.TrimSpace(out), nil
}

func Set(exec Exec, webhookURL, app string) error {
	if !webhookRe.MatchString(webhookURL) {
		return ErrInvalidWebhook
	}
	_, err := exec("discord:set " + webhookURL + appArg(app))
	return classify(err)
}

func Clear(exec Exec, app string) error {
	_, err := exec("discord:clear" + appArg(app))
	return classify(err)
}

func Test(exec Exec, app string) error {
	_, err := exec("discord:test" + appArg(app))
	return classify(err)
}

// Mask hides the webhook token, which is a credential.
func Mask(webhookURL string) string {
	i := strings.LastIndex(webhookURL, "/")
	if i < 0 {
		return "****"
	}
	return webhookURL[:i+1] + "****"
}

// Notify posts a single embed to the webhook.
func Notify(webhookURL, title, description string, color int) error {
	body, err := json.Marshal(map[string]any{
		"embeds": []map[string]any{{
			"title":       title,
			"description": description,
			"color":       color,
			"timestamp":   time.Now().UTC().Format(time.RFC3339),
		}},
	})
	if err != nil {
		return err
	}
	resp, err := httpClient.Post(webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("posting discord webhook: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("discord webhook returned %s", resp.Status)
	}
	return nil
}
