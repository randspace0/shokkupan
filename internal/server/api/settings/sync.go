package settings

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/randspace0/shokkupan/internal/env"
	"github.com/randspace0/shokkupan/internal/server/dokku"
)

// SyncDokkuState re-reads apps and services from dokku into shokkupan's DB.
func SyncDokkuState(e *env.Env, c echo.Context) error {
	dokku.SyncState(e)
	return c.NoContent(http.StatusOK)
}
