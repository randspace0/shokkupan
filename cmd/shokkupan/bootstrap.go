package main

import (
	"github.com/randspace0/shokkupan/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func bootstrapServer() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if err := server.Bootstrap(); err != nil {
		log.Fatal().Err(err).Msg("failed to create key")
	}
}
