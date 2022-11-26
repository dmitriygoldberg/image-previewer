package internal

import (
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

type Handler struct {
	Logger zerolog.Logger
}

func (handler *Handler) GetRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/fill/{width:[0-9]+}/{height:[0-9]+}/{imageURL:.*}", handler.newFillHandler(handler.Logger).Process)

	return router
}

func (handler *Handler) newFillHandler(logger zerolog.Logger) *fillHandler {
	return &fillHandler{logger: logger}
}
