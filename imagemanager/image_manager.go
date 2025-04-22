package imagemanager

import (
	"github.com/dgraph-io/ristretto/v2"
	"github.com/follow1123/photos/logger"
)

type ImageCache struct {
	*ristretto.Cache[string, []byte]
}

type ImageManager struct {
	logger    *logger.AppLogger
	filesRoot string
	cache     *ImageCache
}

func (im *ImageManager) Deinit() {
	im.cache.Close()
}

func NewImageManager(filesRoot string, logger *logger.AppLogger) *ImageManager {
	cache, err := ristretto.NewCache(&ristretto.Config[string, []byte]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		logger.Fatal("init cache error: %v", err)
	}

	return &ImageManager{filesRoot: filesRoot, logger: logger, cache: &ImageCache{Cache: cache}}
}

func (im *ImageManager) NewUploadManager(source ImageSource) (*UploadImageManager, error) {
	return NewUploadImageManager(im.filesRoot, source, im.logger, im.cache)
}

func (im *ImageManager) NewDownloadManager(uri string) *DownloadImageManager {
	return NewDownloadImageManager(im.filesRoot, uri, im.logger, im.cache)
}

func (im *ImageManager) NewDeleteManager(uri string) *DeleteImageManager {
	return NewDeleteImageManager(im.filesRoot, uri)
}
