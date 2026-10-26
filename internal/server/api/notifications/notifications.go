package notifications

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/randspace0/shokkupan/internal/discord"
	"github.com/randspace0/shokkupan/internal/env"
	"github.com/randspace0/shokkupan/internal/server/dto"
	"net/http"
)

func RegisterRoutes(e *env.Env, g *echo.Group) {
	g.GET("", e.H(GetNotifications))
	g.POST("", e.H(SetNotifications))
	g.DELETE("", e.H(ClearNotifications))
	g.POST("/test", e.H(TestNotifications))
}

// mapErr turns plugin sentinel errors into HTTP errors.
func mapErr(err error, action string) error {
	switch {
	case errors.Is(err, discord.ErrNotInstalled):
		return echo.NewHTTPError(http.StatusFailedDependency, err.Error())
	case errors.Is(err, discord.ErrNotConfigured):
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	case errors.Is(err, discord.ErrInvalidWebhook):
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return fmt.Errorf("%s: %w", action, err)
}

func GetNotifications(e *env.Env, c echo.Context) error {
	var req dto.NotificationsRequest
	if err := dto.BindRequest(c, &req); err != nil {
		return err.ToHTTP()
	}

	url, err := discord.Get(e.Dokku.Exec, req.App)
	switch {
	case errors.Is(err, discord.ErrNotInstalled):
		return c.JSON(http.StatusOK, dto.GetNotificationsResponse{})
	case errors.Is(err, discord.ErrNotConfigured):
		return c.JSON(http.StatusOK, dto.GetNotificationsResponse{Installed: true})
	case err != nil:
		return fmt.Errorf("getting discord webhook: %w", err)
	}

	return c.JSON(http.StatusOK, dto.GetNotificationsResponse{
		Installed:  true,
		Configured: true,
		WebhookURL: discord.Mask(url),
	})
}

func SetNotifications(e *env.Env, c echo.Context) error {
	var req dto.SetNotificationsRequest
	if err := dto.BindRequest(c, &req); err != nil {
		return err.ToHTTP()
	}
	if err := discord.Set(e.Dokku.Exec, req.WebhookURL, req.App); err != nil {
		return mapErr(err, "setting discord webhook")
	}
	return c.NoContent(http.StatusOK)
}

func ClearNotifications(e *env.Env, c echo.Context) error {
	var req dto.NotificationsRequest
	if err := dto.BindRequest(c, &req); err != nil {
		return err.ToHTTP()
	}
	if err := discord.Clear(e.Dokku.Exec, req.App); err != nil {
		return mapErr(err, "clearing discord webhook")
	}
	return c.NoContent(http.StatusOK)
}

func TestNotifications(e *env.Env, c echo.Context) error {
	var req dto.NotificationsRequest
	if err := dto.BindRequest(c, &req); err != nil {
		return err.ToHTTP()
	}
	if err := discord.Test(e.Dokku.Exec, req.App); err != nil {
		return mapErr(err, "testing discord webhook")
	}
	return c.NoContent(http.StatusOK)
}
