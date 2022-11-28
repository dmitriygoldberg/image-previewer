package previewer

import (
	"bytes"
	"errors"
	"fmt"
	"image/jpeg"
	"io"
	"os"
	"time"

	"github.com/disintegration/imaging"
	"github.com/rs/zerolog"
)

type cache interface {
	Set(key string, value interface{}) bool
	Get(key string) (interface{}, bool)
	Clear()
}

type Image struct {
	Img     []byte
	Headers map[string][]string
}

type Previewer struct {
	logger          zerolog.Logger
	downloader      Downloader
	resizedCache    cache
	downloadedCache cache
}

func (p Previewer) Fill(params *FillParams) (*FillResponse, error) {
	cacheKeyResized := p.makeCacheKeyResizes(params.width, params.height, params.url)
	cachedResizedImg, ok := p.resizedCache.Get(cacheKeyResized)

	if ok {
		fillResponse := NewFillResponse(cachedResizedImg.(*Image).Img, cachedResizedImg.(*Image).Headers)
		return fillResponse, nil
	}

	cacheKeyDownloaded := p.makeCacheKeyDownloaded(params.url)
	cachedDownloadedImage, ok := p.downloadedCache.Get(cacheKeyDownloaded)

	var downloaded *Image
	var err error

	if ok {
		downloaded = &Image{
			Img:     cachedDownloadedImage.(*Image).Img,
			Headers: cachedDownloadedImage.(*Image).Headers,
		}
	} else {
		downloaded, err = p.downloader.DownloadByURL(params.ctx, params.url, params.headers)
		if err != nil {
			p.logger.Err(err).Msg("Image can't be downloaded")
			return nil, err
		}

		p.downloadedCache.Set(cacheKeyDownloaded, downloaded)
	}

	if downloaded == nil {
		return nil, errors.New("image download error")
	}

	resizedImg, err := p.resize(downloaded.Img, params.width, params.height)
	if err != nil {
		p.logger.Err(err).Msg("Image can not be resized")
		return nil, err
	}

	p.resizedCache.Set(cacheKeyResized, &Image{
		Img:     resizedImg,
		Headers: downloaded.Headers,
	})

	return NewFillResponse(resizedImg, downloaded.Headers), nil
}

func NewPreviewer(
	logger zerolog.Logger,
	downloader Downloader,
	resizedCache cache,
	downloadedCache cache,
) *Previewer {
	return &Previewer{
		logger:          logger,
		downloader:      downloader,
		resizedCache:    resizedCache,
		downloadedCache: downloadedCache,
	}
}

func (p Previewer) resize(img []byte, width, height int) ([]byte, error) {
	tmpImgName := fmt.Sprintf("./%d.jpg", time.Now().Unix())
	file, err := os.Create(tmpImgName)
	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			p.logger.Err(err).Msg("File handling error")
		}
	}(file)

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			p.logger.Err(err).Msg("Tmp file removing error")
		}
	}(tmpImgName)

	if _, err := io.Copy(file, bytes.NewReader(img)); err != nil {
		p.logger.Err(err).Msg(err.Error())
	}

	src, err := imaging.Open(tmpImgName)
	if err != nil {
		return nil, err
	}

	resized := imaging.Resize(src, width, height, imaging.Lanczos)
	imgBuffer := new(bytes.Buffer)
	err = jpeg.Encode(imgBuffer, resized, nil)

	return imgBuffer.Bytes(), err
}

func (p Previewer) makeCacheKeyResizes(width, height int, url string) string {
	return fmt.Sprintf("%d_%d_%s", width, height, url)
}

func (p Previewer) makeCacheKeyDownloaded(url string) string {
	return url
}
