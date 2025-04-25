package imagemanager

import (
	"crypto/md5"
	"encoding/hex"
	"os"

	"github.com/follow1123/photos/logger"
)

type UploadImageManager struct {
	logger    *logger.AppLogger
	filesRoot string
	source    ImageSource
	processor *ImageProcessor
	cache     *ImageCache
}

func NewUploadImageManager(
	filesRoot string,
	imageSource ImageSource,
	logger *logger.AppLogger,
	cache *ImageCache,
) *UploadImageManager {
	return &UploadImageManager{
		filesRoot: filesRoot,
		source:    imageSource,
		cache:     cache,
		logger:    logger,
	}
}
func (uim *UploadImageManager) initImageProcessor() error {
	if uim.processor == nil {
		rc, err := uim.source.GetReader()
		if err != nil {
			return err
		}
		uim.processor = NewImageProcessor(rc, uim.logger)
	}
	return nil
}

func (uim *UploadImageManager) Save() (string, error) {
	err := uim.initImageProcessor()
	if err != nil {
		return "", err
	}
	// TODO check other file type
	fileUri := CreateLocalFileUri(uim.filesRoot)

	originalFileName := fileUri.GetOriginalFilePath()
	compressedFileName := fileUri.GetCompressedFilePath()

	err = fileUri.CreateFilePath()
	if err != nil {
		return "", err
	}

	if fileUri.Is(LOCAL_FILE) {
		data, err := uim.processor.GetData()
		if err != nil {
			return "", err
		}
		os.WriteFile(originalFileName, data, 0666)
	}

	data, err := uim.processor.GetCompressedData()
	uim.cache.Set(fileUri.String(), data, 1)
	if err != nil {
		return "", err
	}
	os.WriteFile(compressedFileName, data, 0666)

	return fileUri.String(), nil
}

func (uim *UploadImageManager) GetImageName() string {
	return uim.source.GetName()
}

func (uim *UploadImageManager) GetHexSum() (string, error) {
	err := uim.initImageProcessor()
	if err != nil {
		return "", err
	}

	data, err := uim.processor.GetData()
	if err != nil {
		return "", err
	}

	sum := md5.Sum(data)
	return hex.EncodeToString(sum[:]), nil
}

func (uim *UploadImageManager) GetImageInfo() (*ImageInfo, error) {
	err := uim.initImageProcessor()
	if err != nil {
		return nil, err
	}

	return uim.processor.GetImageInfo()
}
