package internal

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
)

type fillHandler struct {
	logger zerolog.Logger
}

func (fh *fillHandler) Process(response http.ResponseWriter, request *http.Request) {
	var respContent struct {
		Res string
	}

	respContent.Res = fmt.Sprintf("%s %s", "Hello", "World")
	rawResp, err := json.Marshal(respContent)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = response.Write(rawResp); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		fh.logger.Err(err).Msg("Handler processing error")
		return
	}

	_ = request
}
