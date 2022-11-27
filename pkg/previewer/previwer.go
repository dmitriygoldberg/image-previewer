package previewer

import "github.com/rs/zerolog"

type cache interface {
	Set(key string, value interface{}) bool
	Get(key string) (interface{}, bool)
	Clear()
}

type Previewer struct {
	logger          zerolog.Logger
	downloader      Downloader
	resizedCache    cache
	downloadedCache cache
}

func (p Previewer) Fill(params *FillParams) (*FillResponse, error) {
	// TODO implement me
	panic("implement me")
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
