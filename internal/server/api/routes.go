package api

import (
	"github.com/labstack/echo/v4"
	"github.com/randspace0/shokkupan/internal/env"
	"github.com/randspace0/shokkupan/internal/server/api/apps"
	"github.com/randspace0/shokkupan/internal/server/api/auth"
	"github.com/randspace0/shokkupan/internal/server/api/monitoring"
	"github.com/randspace0/shokkupan/internal/server/api/notifications"
	"github.com/randspace0/shokkupan/internal/server/api/services"
	"github.com/randspace0/shokkupan/internal/server/api/settings"
	"github.com/randspace0/shokkupan/internal/server/api/setup"
)

func RegisterRoutes(e *env.Env, authMiddleware []echo.MiddlewareFunc) {
	g := e.Router.Group("/api")

	setup.RegisterRoutes(e, g.Group("/setup"))
	auth.RegisterRoutes(e, g.Group("/auth"), g.Group("/auth", authMiddleware...))

	protectedApi := g.Group("", authMiddleware...)
	apps.RegisterRoutes(e, protectedApi.Group("/apps"))
	monitoring.RegisterRoutes(e, protectedApi.Group("/monitoring"))
	notifications.RegisterRoutes(e, protectedApi.Group("/notifications"))
	services.RegisterRoutes(e, protectedApi.Group("/services"))
	settings.RegisterRoutes(e, protectedApi.Group("/settings"))

	protectedApi.GET("/exec/status", e.H(GetCommandExecutionStatus))
	g.POST("/github/events", e.H(ReceiveGithubWebhookEvent))
}
