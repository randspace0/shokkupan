package setup

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/randspace0/shokkupan/internal/access"
	"github.com/randspace0/shokkupan/internal/models"
	"github.com/randspace0/shokkupan/internal/randstr"
	"github.com/randspace0/shokkupan/internal/server/auth"
	"github.com/rs/zerolog/log"
	"image/png"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pquerna/otp/totp"
	"github.com/randspace0/shokkupan/internal/env"
	"github.com/randspace0/shokkupan/internal/server/dto"
)

func GenerateTotp(e *env.Env, c echo.Context) error {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "shokkupan",
		AccountName: "admin account",
	})
	if err != nil {
		return fmt.Errorf("failed to generate totp key: %w", err)
	}

	var imgBuf bytes.Buffer
	img, imgErr := key.Image(160, 160)
	if imgErr != nil {
		return fmt.Errorf("failed to generate totp image: %w", imgErr)
	}
	if encErr := png.Encode(&imgBuf, img); encErr != nil {
		return fmt.Errorf("failed to encode totp image to png: %w", encErr)
	}
	return c.JSON(http.StatusOK, dto.GenerateTotpResponse{
		Secret:       key.Secret(),
		Image:        base64.StdEncoding.EncodeToString(imgBuf.Bytes()),
		RecoveryCode: randstr.String(12),
	})
}

func ConfirmTotp(e *env.Env, c echo.Context) error {
	var req dto.ConfirmTotpRequest
	if err := dto.BindRequest(c, &req); err != nil {
		return err.ToHTTP()
	}

	return c.JSON(http.StatusOK, dto.ConfirmTotpResponse{
		Valid: totp.Validate(req.Code, req.Secret),
	})
}

func CompletePasswordSetup(e *env.Env, c echo.Context) error {
	var req dto.CompletePasswordSetupRequest
	if err := dto.BindRequest(c, &req); err != nil {
		return err.ToHTTP()
	}

	if err := setupServerWithAuthMethod(e, auth.MethodPassword); err != nil {
		log.Error().Err(err).Msg("failed to setup github auth")
		return echo.ErrInternalServerError
	}

	pw, ok := e.Auth.(*auth.PasswordAuthenticator)
	if !ok {
		log.Error().Msg("failed to cast e.Auth to pw auth")
		return echo.ErrInternalServerError
	}
	passwordHash, pwErr := pw.HashPassword([]byte(req.Password))
	if pwErr != nil {
		log.Error().Err(pwErr).Msg("failed to hash password")
		return echo.ErrInternalServerError
	}

	user := models.User{
		Name:         req.Username,
		Source:       "manual",
		Role:         access.RoleAdmin,
		PasswordHash: passwordHash,
		TotpEnabled:  req.Enable2FA,
		RecoveryCode: req.RecoveryCode,
		TotpSecret:   req.TotpSecret,
	}
	if err := e.DB.Save(&user).Error; err != nil {
		log.Error().Err(err).Msg("failed to save initial password auth user")
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}
