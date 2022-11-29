package previewer

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type ImageDownloader interface {
	DownloadByURL(ctx context.Context, url string, headers map[string][]string) (*Image, error)
}

type Downloader struct{}

func NewDownloader() *Downloader {
	return &Downloader{}
}

func (d *Downloader) DownloadByURL(
	ctx context.Context,
	url string,
	headers map[string][]string,
) (*Image, error) {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Image download error")
	}

	req.Header = headers
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	img, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := d.validate(img); err != nil {
		return nil, err
	}

	downloadedImage := &Image{Img: img, Headers: resp.Header}

	if err := resp.Body.Close(); err != nil {
		return nil, err
	}

	return downloadedImage, nil
}

func (d *Downloader) validate(img []byte) error {
	if len(img) == 0 {
		return errors.New("Image download error")
	}

	allowedFormats := map[string]string{
		"\xff\xd8\xff":      "image/jpeg",
		"\x89PNG\r\n\x1a\n": "image/png",
		"GIF87a":            "image/gif",
		"GIF89a":            "image/gif",
	}

	imgStr := string(img)
	for format := range allowedFormats {
		if strings.HasPrefix(imgStr, format) {
			return nil
		}
	}

	return errors.New("Downloaded file is not image")
}
