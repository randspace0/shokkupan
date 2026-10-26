package env

import (
	"github.com/labstack/echo/v4"
	"github.com/randspace0/shokkupan/internal/server/auth"
	"github.com/texm/dokku-go"
	"gorm.io/gorm"
)

// Version is set at build time: -ldflags "-X github.com/randspace0/shokkupan/internal/env.Version=v1.2.3"
var Version = "dev"

type Env struct {
	Version string
	Router  *echo.Echo
	Dokku   *dokku.SSHClient
	DB      *gorm.DB
	Auth    auth.Authenticator

	DebugMode      bool
	SetupCompleted bool
}

func New(debugMode bool) *Env {
	return &Env{
		DebugMode: debugMode,
		Version:   Version,
	}
}

type handlerFunc func(env *Env, c echo.Context) error

func (e *Env) H(f handlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error { return f(e, c) }
}

func (e *Env) Shutdown() error {
	if err := e.Dokku.Close(); err != nil {
		return err
	}
	return nil
}
