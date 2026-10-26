package apps

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/randspace0/shokkupan/internal/env"
	"github.com/randspace0/shokkupan/internal/server/dto"
)

// build is one row of `dokku builds:list --format json`.
type build struct {
	ID         string `json:"id"`
	Kind       string `json:"kind"`
	Status     string `json:"status"`
	Source     string `json:"source"`
	StartedAt  string `json:"started_at"`
	FinishedAt string `json:"finished_at"`
	ExitCode   *int   `json:"exit_code"`
}

func parseBuilds(out string) ([]build, error) {
	out = strings.TrimSpace(out)
	if out == "" {
		return []build{}, nil
	}
	var builds []build
	if err := json.Unmarshal([]byte(out), &builds); err != nil {
		return nil, fmt.Errorf("parsing builds:list output: %w", err)
	}
	// dokku lists oldest first; show the newest first
	sort.SliceStable(builds, func(i, j int) bool { return builds[i].StartedAt > builds[j].StartedAt })
	return builds, nil
}

func ListAppDeploys(e *env.Env, c echo.Context) error {
	var req dto.GetAppDeploysRequest
	if err := dto.BindRequest(c, &req); err != nil {
		return err.ToHTTP()
	}

	out, err := e.Dokku.Exec("builds:list " + req.Name + " --format json")
	if err != nil {
		if strings.Contains(err.Error(), "is not a dokku command") {
			return echo.NewHTTPError(http.StatusNotImplemented, "this dokku version has no builds:list; upgrade dokku")
		}
		return fmt.Errorf("listing builds: %w", err)
	}
	builds, err := parseBuilds(out)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]any{"deploys": builds})
}

func GetAppDeployOutput(e *env.Env, c echo.Context) error {
	var req dto.GetAppDeployOutputRequest
	if err := dto.BindRequest(c, &req); err != nil {
		return err.ToHTTP()
	}

	// both values are validated (app name charset, alphanumeric id) before
	// being placed in the dokku command line
	out, err := e.Dokku.Exec("builds:output " + req.Name + " " + req.ID)
	if err != nil {
		return fmt.Errorf("getting build output: %w", err)
	}
	return c.JSON(http.StatusOK, map[string]string{"output": out})
}
