package monitoring

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	dokkugo "github.com/texm/dokku-go"
	"gitlab.com/texm/shokku/internal/env"
	"gitlab.com/texm/shokku/internal/models"
	"gitlab.com/texm/shokku/internal/server/dokku"
	"gitlab.com/texm/shokku/internal/server/dto"
	"net/http"
	"sort"
	"time"
)

func RegisterRoutes(e *env.Env, g *echo.Group) {
	g.GET("", e.H(GetMonitoring))
	g.GET("/stats", e.H(GetMonitoringStats))
}

func GetMonitoring(e *env.Env, c echo.Context) error {
	apps := []dto.MonitoringApp{}
	reports, err := e.Dokku.GetAllProcessReport()
	if err != nil && !errors.Is(err, dokkugo.NoDeployedAppsError) {
		return fmt.Errorf("getting process report: %w", err)
	}
	for name, r := range reports {
		if name == dokku.FilteredApp {
			continue
		}
		apps = append(apps, dto.MonitoringApp{
			Name:         name,
			IsDeployed:   r.Deployed,
			IsRunning:    r.Running,
			NumProcesses: r.Processes,
			Restore:      r.Restore,
		})
	}
	sort.Slice(apps, func(i, j int) bool { return apps[i].Name < apps[j].Name })

	var rows []models.AppStatusEvent
	if err := e.DB.Order("created_at desc").Limit(50).Find(&rows).Error; err != nil {
		return fmt.Errorf("getting status events: %w", err)
	}
	events := make([]dto.MonitoringEvent, len(rows))
	for i, r := range rows {
		events[i] = dto.MonitoringEvent{
			AppName: r.AppName,
			Running: r.Running,
			At:      r.CreatedAt.UTC().Format(time.RFC3339),
		}
	}

	// a recent sample means the dokku-stats plugin is installed and working
	var recent int64
	e.DB.Model(&models.AppStatsSample{}).
		Where("created_at > ?", time.Now().Add(-10*time.Minute)).Count(&recent)

	return c.JSON(http.StatusOK, dto.GetMonitoringResponse{
		StatsAvailable: recent > 0,
		Apps:           apps,
		Events:         events,
	})
}

func GetMonitoringStats(e *env.Env, c echo.Context) error {
	var req dto.GetMonitoringStatsRequest
	if err := dto.BindRequest(c, &req); err != nil {
		return err.ToHTTP()
	}
	if req.Hours == 0 {
		req.Hours = 6
	}

	var rows []models.AppStatsSample
	err := e.DB.Where("app_name = ? AND created_at > ?", req.App,
		time.Now().Add(-time.Duration(req.Hours)*time.Hour)).
		Order("created_at asc").Find(&rows).Error
	if err != nil {
		return fmt.Errorf("getting stats: %w", err)
	}

	points := make([]dto.MonitoringStatsPoint, len(rows))
	for i, r := range rows {
		points[i] = dto.MonitoringStatsPoint{
			At:       r.CreatedAt.UTC().Format(time.RFC3339),
			CPU:      r.CPU,
			MemBytes: r.MemBytes,
			MemLimit: r.MemLimit,
		}
	}
	return c.JSON(http.StatusOK, dto.GetMonitoringStatsResponse{Points: points})
}
