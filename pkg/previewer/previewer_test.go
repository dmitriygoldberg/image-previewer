package previewer

import (
	"context"
	"errors"
	lruCache "github.com/dmitriygoldberg/image-previewer/pkg/cache"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestDefaultService_Fill_positive(t *testing.T) {
	imageOrigin := loadImage("_gopher_original_1024x504.jpg")
	image1Resized := loadImage("gopher_100x100.jpg")
	downloadedImage := &Image{Img: imageOrigin}
	resizedCache := lruCache.NewCache(2)
	downloadedCache := lruCache.NewCache(2)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDownloader := NewMockImageDownloader(ctrl)
	logger := log.With().Logger()

	type fields struct {
		l               zerolog.Logger
		downloader      ImageDownloader
		resizedCache    lruCache.Cache
		downloadedCache lruCache.Cache
	}
	tests := []struct {
		name    string
		fields  fields
		params  *FillParams
		want    *FillResponse
		wantErr bool
	}{
		{
			name: "good_resized",
			fields: fields{
				l:               logger,
				downloader:      mockDownloader,
				resizedCache:    resizedCache,
				downloadedCache: downloadedCache,
			},
			params: NewFillParams(
				context.Background(),
				1000,
				500,
				"http://raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg",
				make(map[string][]string),
			),
			want:    NewFillResponse(image1Resized, make(map[string][]string)),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Previewer{
				logger:          tt.fields.l,
				downloader:      tt.fields.downloader,
				resizedCache:    tt.fields.resizedCache,
				downloadedCache: tt.fields.downloadedCache,
			}

			mockDownloader.EXPECT().DownloadByURL(
				context.Background(),
				tt.params.url,
				tt.params.headers,
			).Return(downloadedImage, nil)

			_, err := svc.Fill(tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fill() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDefaultService_Fill_negative(t *testing.T) {
	image1Resized := loadImage("gopher_100x100.jpg")
	resizedCache := lruCache.NewCache(2)
	downloadedCache := lruCache.NewCache(2)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDownloader := NewMockImageDownloader(ctrl)
	logger := log.With().Logger()

	type fields struct {
		l               zerolog.Logger
		downloader      ImageDownloader
		resizedCache    lruCache.Cache
		downloadedCache lruCache.Cache
	}
	tests := []struct {
		name       string
		fields     fields
		params     *FillParams
		want       *FillResponse
		downloaded *Image
		err        error
		wantErr    bool
	}{
		{
			name: "false_resized",
			fields: fields{
				l:               logger,
				downloader:      mockDownloader,
				resizedCache:    resizedCache,
				downloadedCache: downloadedCache,
			},
			params: NewFillParams(
				context.Background(),
				1000,
				500,
				"http://raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg",
				make(map[string][]string),
			),
			want:       NewFillResponse(image1Resized, make(map[string][]string)),
			downloaded: nil,
			err:        errors.New("image can't be downloaded"),
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Previewer{
				logger:          tt.fields.l,
				downloader:      tt.fields.downloader,
				resizedCache:    tt.fields.resizedCache,
				downloadedCache: tt.fields.downloadedCache,
			}

			mockDownloader.EXPECT().DownloadByURL(
				context.Background(),
				tt.params.url,
				tt.params.headers,
			).Return(tt.downloaded, tt.err)

			_, err := svc.Fill(tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fill() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDefaultService_Fill_from_cache(t *testing.T) {
	imageOrigin := loadImage("_gopher_original_1024x504.jpg")
	image1Resized := loadImage("gopher_100x100.jpg")
	downloadedImage := &Image{Img: imageOrigin}
	resizedCache := lruCache.NewCache(2)
	downloadedCache := lruCache.NewCache(2)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDownloader := NewMockImageDownloader(ctrl)
	logger := log.With().Logger()

	type fields struct {
		l               zerolog.Logger
		downloader      ImageDownloader
		resizedCache    lruCache.Cache
		downloadedCache lruCache.Cache
	}
	tests := []struct {
		name    string
		fields  fields
		params  *FillParams
		want    *FillResponse
		wantErr bool
	}{
		{
			name: "good_resized",
			fields: fields{
				l:               logger,
				downloader:      mockDownloader,
				resizedCache:    resizedCache,
				downloadedCache: downloadedCache,
			},
			params: NewFillParams(
				context.Background(),
				1000,
				500,
				"http://raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg",
				make(map[string][]string),
			),
			want:    NewFillResponse(image1Resized, make(map[string][]string)),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Previewer{
				logger:          tt.fields.l,
				downloader:      tt.fields.downloader,
				resizedCache:    tt.fields.resizedCache,
				downloadedCache: tt.fields.downloadedCache,
			}

			mockDownloader.EXPECT().DownloadByURL(
				context.Background(),
				tt.params.url,
				tt.params.headers,
			).Return(downloadedImage, nil)

			_, err := svc.Fill(tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fill() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			cachedResp, err := svc.Fill(tt.params)
			cacheKey := svc.makeCacheKeyResizes(tt.params.width, tt.params.height, tt.params.url)
			cacheResized, ok := resizedCache.Get(cacheKey)
			require.True(t, ok)
			require.Equal(t, cachedResp.Img, cacheResized.(*Image).Img)

			tt.params.width++
			_, err = svc.Fill(tt.params)
			cachedKey := svc.makeCacheKeyDownloaded(tt.params.url)
			cacheDownloaded, ok := downloadedCache.Get(cachedKey)
			require.True(t, ok)
			require.Equal(t, imageOrigin, cacheDownloaded.(*Image).Img)
		})
	}
}

func TestNewDefaultService(t *testing.T) {
	type args struct {
		l               zerolog.Logger
		downloader      ImageDownloader
		resizedCache    lruCache.Cache
		downloadedCache lruCache.Cache
	}
	tests := []struct {
		name string
		args args
		want *Previewer
	}{
		{
			name: "good",
			want: &Previewer{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPreviewer(tt.args.l, tt.args.downloader, tt.args.resizedCache, tt.args.downloadedCache); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDefaultService() = %v, want %v", got, tt.want)
			}
		})
	}
}
