package main

import (
	"github.com/dmitriygoldberg/image-previewer/internal"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

var shaCommit = "local"

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := log.With().Str("sha_commit", shaCommit).Logger()

	if err := godotenv.Load(); err != nil {
		logger.Warn().Err(err).Msg("Loading .env error")
	}

	handler := internal.Handler{Logger: logger}
	server := &http.Server{
		Addr:         ":8080",
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Second * 60,
		Handler:      handler.GetRouter(),
	}

	if err := server.ListenAndServe(); err != nil {
		logger.Fatal().Err(err).Msg("Error starting http server")
	}
}
