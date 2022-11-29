package internal

import (
	"errors"
	"github.com/dmitriygoldberg/image-previewer/pkg/previewer"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func TestHandlers_FillHandler_Positive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPreviewer := previewer.NewMockPreviewInterface(ctrl)
	l := log.With().Logger()

	image1 := loadImage("gopher_200x700.jpg")
	image2 := loadImage("gopher_1024x252.jpg")

	tests := []struct {
		name         string
		width        int64
		height       int64
		url          string
		response     string
		fillResponse *previewer.FillResponse
		err          error
		httpStatus   int64
	}{
		{
			name:         "good response",
			width:        200,
			height:       300,
			url:          "http://raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/gopher_200x700.jpg",
			response:     string(image1),
			fillResponse: &previewer.FillResponse{Img: image1},
			httpStatus:   http.StatusOK,
		},
		{
			name:         "good response 2",
			width:        400,
			height:       500,
			url:          "http://raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/gopher_1024x252.jpg",
			response:     string(image2),
			fillResponse: &previewer.FillResponse{Img: image2},
			httpStatus:   http.StatusOK,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
			req = mux.SetURLVars(req, map[string]string{
				"width":    strconv.Itoa(int(tt.width)),
				"height":   strconv.Itoa(int(tt.height)),
				"imageURL": tt.url,
			})

			fillParams := previewer.NewFillParams(req.Context(), int(tt.width), int(tt.height), tt.url, req.Header)
			mockPreviewer.EXPECT().Fill(fillParams).Return(tt.fillResponse, tt.err)
			fh := &fillHandler{
				logger:    l,
				previewer: mockPreviewer,
			}

			w := httptest.NewRecorder()
			fh.Process(w, req)
			require.Equal(t, int(tt.httpStatus), w.Result().StatusCode)
			require.Equal(t, strings.TrimSpace(w.Body.String()), tt.response)
		})
	}
}

func TestHandlers_FillHandler_Negative(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPreviewer := previewer.NewMockPreviewInterface(ctrl)
	l := log.With().Logger()

	tests := []struct {
		name         string
		width        int64
		height       int64
		url          string
		response     string
		fillResponse *previewer.FillResponse
		err          error
		httpStatus   int64
	}{
		{
			name:       "validation error",
			width:      300,
			height:     400,
			url:        "http://user^:passwo^rd@foo.com/",
			response:   "Request validation error",
			httpStatus: http.StatusBadRequest,
		},
		{
			name:         "fill error",
			width:        300,
			height:       400,
			url:          "http://raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/gopher_200x700.jpg",
			response:     "Unable to process image",
			fillResponse: nil,
			httpStatus:   http.StatusBadGateway,
			err:          errors.New("ошибка"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
			req = mux.SetURLVars(req, map[string]string{
				"width":    strconv.Itoa(int(tt.width)),
				"height":   strconv.Itoa(int(tt.height)),
				"imageURL": tt.url,
			})

			if tt.fillResponse != nil || tt.err != nil {
				fillParams := previewer.NewFillParams(req.Context(), int(tt.width), int(tt.height), tt.url, req.Header)
				mockPreviewer.EXPECT().Fill(fillParams).Return(tt.fillResponse, tt.err)
			}

			fh := &fillHandler{
				logger:    l,
				previewer: mockPreviewer,
			}

			w := httptest.NewRecorder()

			fh.Process(w, req)
			require.Equal(t, int(tt.httpStatus), w.Result().StatusCode)
			require.Equal(t, strings.TrimSpace(w.Body.String()), tt.response)
		})
	}
}

func TestHandlers_FillHandler_Headers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPreviewer := previewer.NewMockPreviewInterface(ctrl)
	l := log.With().Logger()

	image1 := loadImage("gopher_200x700.jpg")

	headers := map[string][]string{
		"Source-Age":                  {0: "3"},
		"Access-Control-Allow-Origin": {0: "*"},
		"Content-Length":              {0: "30146"},
		"Content-Type":                {0: "image/jpeg"},
	}

	tests := []struct {
		name         string
		width        int64
		height       int64
		url          string
		response     string
		fillResponse *previewer.FillResponse
		err          error
		httpStatus   int64
	}{
		{
			name:         "good headers",
			width:        200,
			height:       300,
			url:          "http://raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/gopher_200x700.jpg",
			response:     string(image1),
			fillResponse: &previewer.FillResponse{Img: image1, Headers: headers},
			httpStatus:   http.StatusOK,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
			req = mux.SetURLVars(req, map[string]string{
				"width":    strconv.Itoa(int(tt.width)),
				"height":   strconv.Itoa(int(tt.height)),
				"imageURL": tt.url,
			})

			fillParams := previewer.NewFillParams(req.Context(), int(tt.width), int(tt.height), tt.url, req.Header)
			mockPreviewer.EXPECT().Fill(fillParams).Return(tt.fillResponse, tt.err)
			fh := &fillHandler{
				logger:    l,
				previewer: mockPreviewer,
			}

			w := httptest.NewRecorder()
			fh.Process(w, req)

			for name, values := range tt.fillResponse.Headers {
				for _, value := range values {
					require.Equal(t, value, w.Header().Get(name))
				}
			}
		})
	}
}
