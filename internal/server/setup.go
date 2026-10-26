package server

import (
	"fmt"
	"github.com/randspace0/shokkupan/internal/env"
	"github.com/randspace0/shokkupan/internal/models"
	"github.com/randspace0/shokkupan/internal/randstr"
	"github.com/randspace0/shokkupan/internal/server/auth"
	"github.com/randspace0/shokkupan/internal/server/github"
	"github.com/rs/zerolog/log"
	"time"
)

func maybeSyncState(e *env.Env, state *models.Server) error {
	if state.AuthMethod == auth.MethodGithub {
		timeSinceLast := time.Since(state.LastSync)
		needsSync := timeSinceLast > 6*time.Hour
		if !needsSync {
			return nil
		}
		if err := github.SyncUsersToDB(e); err != nil {
			return fmt.Errorf("failed to sync github state: %w", err)
		}
	}

	state.LastSync = time.Now()
	if err := e.DB.Save(state).Error; err != nil {
		return fmt.Errorf("failed to update last sync time")
	}

	return nil
}

func initServerState(e *env.Env) (*models.Server, error) {
	var state models.Server
	if err := e.DB.FirstOrCreate(&state).Error; err != nil {
		log.Error().
			Err(err).
			Msg("failed to get state")
		return nil, err
	}

	log.Debug().
		Bool("is_setup", state.IsSetup).
		Str("auth_method", string(state.AuthMethod)).
		Time("last_sync", state.LastSync).
		Msg("server setup state")

	e.SetupCompleted = state.IsSetup

	if state.IsSetup {
		if syncErr := maybeSyncState(e, &state); syncErr != nil {
			return nil, syncErr
		}
		return &state, nil
	}

	if state.SetupKey == "" {
		state.SetupKey = randstr.String(16)
		if err := e.DB.Save(&state).Error; err != nil {
			log.Error().
				Err(err).
				Msg("failed to update setup key")
			return nil, err
		}
	}

	log.Info().
		Str("setup_key", state.SetupKey).
		Msg("running in setup mode")

	return &state, nil
}
