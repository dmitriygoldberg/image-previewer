package internal

import (
	"github.com/dmitriygoldberg/image-previewer/config"
	"github.com/dmitriygoldberg/image-previewer/pkg/cache"
	"github.com/dmitriygoldberg/image-previewer/pkg/previewer"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

type Handler struct {
	Logger    zerolog.Logger
	AppConfig config.Config
}

func (handler *Handler) GetRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/fill/{width:[0-9]+}/{height:[0-9]+}/{imageURL:.*}", handler.fillHandler().Process)

	return router
}

func (handler *Handler) fillHandler() *fillHandler {
	resizedCache := cache.NewCache(handler.AppConfig.Cache.Capacity)
	downloadedCache := cache.NewCache(handler.AppConfig.Cache.Capacity)

	downloader := previewer.NewDownloader()
	previewerService := previewer.NewPreviewer(handler.Logger, downloader, resizedCache, downloadedCache)

	return &fillHandler{logger: handler.Logger, previewer: previewerService}
}
