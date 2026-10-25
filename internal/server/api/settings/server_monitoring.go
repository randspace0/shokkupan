package settings

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/texm/shokku/internal/env"
	"gitlab.com/texm/shokku/internal/models"
)

type serverMonitoringBody struct {
	Enabled bool `json:"enabled"`
}

func GetServerMonitoring(e *env.Env, c echo.Context) error {
	var s models.Server
	if err := e.DB.First(&s).Error; err != nil {
		return fmt.Errorf("getting server state: %w", err)
	}
	return c.JSON(http.StatusOK, serverMonitoringBody{Enabled: s.ServerMonitoring})
}

func SetServerMonitoring(e *env.Env, c echo.Context) error {
	var body serverMonitoringBody
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	err := e.DB.Model(&models.Server{}).Where("1 = 1").
		Update("server_monitoring", body.Enabled).Error
	if err != nil {
		return fmt.Errorf("saving server monitoring: %w", err)
	}
	return c.JSON(http.StatusOK, body)
}
