package imagemanager

import (
	"github.com/follow1123/photos/logger"
)

type ImageManager struct {
	logger    *logger.AppLogger
	filesRoot string
	cache     *ImageCache
}

func (im *ImageManager) Deinit() {
	im.cache.Close()
}

func NewImageManager(filesRoot string, cache *ImageCache, logger *logger.AppLogger) *ImageManager {
	return &ImageManager{
		filesRoot: filesRoot,
		logger:    logger,
		cache:     cache,
	}
}

func (im *ImageManager) NewUploadManager(source ImageSource) *UploadImageManager {
	return NewUploadImageManager(im.filesRoot, source, im.logger, im.cache)
}

func (im *ImageManager) NewDownloadManager(uri string) *DownloadImageManager {
	return NewDownloadImageManager(im.filesRoot, uri, im.logger, im.cache)
}

func (im *ImageManager) NewDeleteManager(uri string) *DeleteImageManager {
	return NewDeleteImageManager(im.filesRoot, uri)
}
