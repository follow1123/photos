package imagemanager

import (
	"io"
	"os"

	"github.com/follow1123/photos/logger"
)

type DownloadImageManager struct {
	logger.AppLogger
	uri   FileUri
	cache *ImageCache
}

func NewDownloadImageManager(
	filesRoot string,
	uri string,
	logger *logger.AppLogger,
	cache *ImageCache,
) *DownloadImageManager {
	return &DownloadImageManager{
		uri:       *NewFileUri(filesRoot, uri),
		cache:     cache,
		AppLogger: *logger,
	}
}

func (dim *DownloadImageManager) OpenOriginal() (io.ReadCloser, error) {
	originalFilePath := dim.uri.GetOriginalFilePath()
	dim.Debug("download manager open original file: %s", originalFilePath)
	if dim.uri.Is(LOCAL_FILE) {
		file, err := os.Open(originalFilePath)
		if err != nil {
			dim.Debug("download manager open original file %s error: %v", originalFilePath, err)
			return nil, err
		}
		return file, nil
	} else {
		panic("TODO: implement read remote file path")
	}
}

func (dim *DownloadImageManager) GetCompressed() ([]byte, error) {
	data, ok := dim.cache.Get(dim.uri.String())
	if ok {
		dim.Debug("get compressed image from cache")
		return data, nil
	}
	dim.Debug("get compressed image from local file")
	compressedFilePath := dim.uri.GetCompressedFilePath()
	imageData, err := os.ReadFile(compressedFilePath)
	if err != nil {
		return nil, err
	}
	dim.cache.Set(dim.uri.String(), imageData, 1)
	return imageData, nil
}
