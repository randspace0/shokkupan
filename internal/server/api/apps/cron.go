package apps

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/randspace0/shokkupan/internal/env"
	"github.com/randspace0/shokkupan/internal/server/commands"
	"github.com/randspace0/shokkupan/internal/server/dto"
	"github.com/texm/dokku-go"
	"net/http"
	"strings"
)

// parseCronList parses `dokku cron:list` output. Dokku pads columns, but
// the command column holds spaces, so split by field count instead: a
// schedule is 5 fields, or 1 for macros like @daily.
func parseCronList(out string) []dto.CronTask {
	tasks := []dto.CronTask{}
	for _, line := range strings.Split(out, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 3 || fields[0] == "ID" || strings.HasPrefix(fields[0], "!") {
			continue
		}
		n := 5
		if strings.HasPrefix(fields[1], "@") {
			n = 1
		}
		if len(fields) < 1+n+1 {
			continue
		}
		tasks = append(tasks, dto.CronTask{
			ID:       fields[0],
			Schedule: strings.Join(fields[1:1+n], " "),
			Command:  strings.Join(fields[1+n:], " "),
		})
	}
	return tasks
}

func GetAppCron(e *env.Env, c echo.Context) error {
	var req dto.GetAppOverviewRequest
	if err := dto.BindRequest(c, &req); err != nil {
		return err.ToHTTP()
	}

	out, err := e.Dokku.Exec("cron:list " + req.Name)
	if err != nil {
		var xe *dokku.ExitCodeError
		// apps without cron entries make dokku exit non-zero
		if errors.As(err, &xe) && strings.Contains(xe.Output(), "No cron") {
			return c.JSON(http.StatusOK, dto.GetAppCronResponse{Tasks: []dto.CronTask{}})
		}
		return fmt.Errorf("listing cron tasks: %w", err)
	}

	return c.JSON(http.StatusOK, dto.GetAppCronResponse{Tasks: parseCronList(out)})
}

func RunAppCron(e *env.Env, c echo.Context) error {
	var req dto.RunAppCronRequest
	if err := dto.BindRequest(c, &req); err != nil {
		return err.ToHTTP()
	}

	cmd := func() (*dokku.CommandOutputStream, error) {
		return e.Dokku.ExecStreaming(fmt.Sprintf("cron:run %s %s", req.Name, req.ID))
	}

	return c.JSON(http.StatusOK, dto.CommandExecutionResponse{
		ExecutionID: commands.RequestExecution(cmd, nil),
	})
}
