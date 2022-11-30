package main

import (
	"net/http"

	"github.com/dmitriygoldberg/image-previewer/config"
	"github.com/dmitriygoldberg/image-previewer/internal"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var shaCommit = "local"

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := log.With().Str("sha_commit", shaCommit).Logger()

	initEnv(logger)
	appConfig := config.New()

	handler := internal.Handler{Logger: logger, AppConfig: *appConfig}
	server := &http.Server{
		Addr:         appConfig.Server.Address,
		WriteTimeout: appConfig.Server.WriteTimeout,
		ReadTimeout:  appConfig.Server.ReadTimeout,
		IdleTimeout:  appConfig.Server.IdleTimeout,
		Handler:      handler.GetRouter(),
	}

	if err := server.ListenAndServe(); err != nil {
		logger.Fatal().Err(err).Msg("Error starting http server")
	}
}

func initEnv(logger zerolog.Logger) {
	if err := godotenv.Load(); err != nil {
		logger.Warn().Err(err).Msg("Loading .env error")
	}
}
