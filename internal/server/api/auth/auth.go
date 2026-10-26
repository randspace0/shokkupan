package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/randspace0/shokkupan/internal/env"
	"github.com/randspace0/shokkupan/internal/models"
	"github.com/randspace0/shokkupan/internal/server/auth"
	"github.com/rs/zerolog/log"
	"net/http"
)

func HandleLogout(e *env.Env, c echo.Context) error {
	e.Auth.ClearTokenCookies(c)
	return c.NoContent(http.StatusOK)
}

func HandleRefreshAuth(e *env.Env, c echo.Context) error {
	user, err := e.Auth.GetUserFromContext(c)
	if err != nil {
		log.Error().Msg("failed to parse user from context")
		return echo.ErrInternalServerError
	}

	claims := auth.UserClaims{
		Name: user.Name,
	}
	token, err := e.Auth.NewToken(claims)
	if err != nil {
		log.Error().Err(err).Msg("failed to create jwt")
		return echo.ErrInternalServerError
	}

	e.Auth.SetTokenCookies(c, token)

	return c.NoContent(http.StatusOK)
}

func HandleGetDetails(e *env.Env, c echo.Context) error {
	user, err := e.Auth.GetUserFromContext(c)
	if err != nil {
		log.Error().Msg("failed to parse user from context")
		return echo.ErrInternalServerError
	}

	var dbUser models.User
	if err := e.DB.Where("name = ?", user.Name).Take(&dbUser).Error; err != nil {
		return echo.ErrUnauthorized
	}
	apps := []string{}
	var rows []models.AppAccess
	e.DB.Where("user_name = ?", user.Name).Find(&rows)
	for _, r := range rows {
		apps = append(apps, r.AppName)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"username": user.Name,
		"role":     dbUser.Role,
		"apps":     apps,
	})
}
