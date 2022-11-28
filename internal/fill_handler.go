package internal

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/dmitriygoldberg/image-previewer/pkg/previewer"
	"github.com/gorilla/mux"

	"github.com/rs/zerolog"
)

type fillHandler struct {
	logger    zerolog.Logger
	previewer *previewer.Previewer
}

func (fh *fillHandler) Process(response http.ResponseWriter, request *http.Request) {
	fillParams, err := fh.parseParams(request.Context(), mux.Vars(request), request.Header)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("Request validation error"))
		fh.logger.Err(err).Msg(err.Error())
		return
	}

	fillResponse, err := fh.previewer.Fill(fillParams)
	if err != nil {
		response.WriteHeader(http.StatusBadGateway)
		response.Write([]byte("Unable to process image"))
		fh.logger.Err(err).Msg(err.Error())
		return
	}

	for key, values := range fillResponse.Headers {
		for _, value := range values {
			response.Header().Set(key, value)
		}
	}

	response.Header().Set("Content-Length", strconv.Itoa(len(fillResponse.Img)))
	if _, err := response.Write(fillResponse.Img); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		fh.logger.Err(err).Msg(err.Error())
		return
	}
}

func (fh *fillHandler) parseParams(
	ctx context.Context,
	vars map[string]string,
	headers map[string][]string,
) (*previewer.FillParams, error) {
	width, err := strconv.Atoi(vars["width"])
	if err != nil {
		return nil, errors.New("width is not correct")
	}

	height, err := strconv.Atoi(vars["height"])
	if err != nil {
		return nil, errors.New("height is not correct")
	}

	imageURL, err := url.Parse(vars["imageURL"])
	if err != nil {
		return nil, errors.New("imageURL is not correct")
	}

	imageURL.Scheme = "http"

	return previewer.NewFillParams(ctx, width, height, imageURL.String(), headers), nil
}
