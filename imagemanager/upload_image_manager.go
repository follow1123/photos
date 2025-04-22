package imagemanager

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/follow1123/photos/logger"
)

type UploadImageManager struct {
	logger.AppLogger
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
) (*UploadImageManager, error) {
	r, err := imageSource.GetReader()
	if err != nil {
		return nil, err
	}
	defer r.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return &UploadImageManager{
		filesRoot: filesRoot,
		source:    imageSource,
		processor: &ImageProcessor{data: data},
		cache:     cache,
		AppLogger: *logger,
	}, nil
}

func (uim *UploadImageManager) saveFile(filePath string, src io.Reader) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, src)
	if err != nil {
		return err
	}
	return nil
}

func (uim *UploadImageManager) Save() (string, error) {
	// TODO check other file type
	fileUri := createLocalFileUri(uim.filesRoot)

	originalFileName := fileUri.GetOriginalFilePath()
	compressedFileName := fileUri.GetCompressedFilePath()

	err := os.MkdirAll(filepath.Dir(compressedFileName), 0755)
	if err != nil {
		return "", err
	}

	if fileUri.Is(LOCAL_FILE) {
		uim.saveFile(originalFileName, bytes.NewReader(uim.processor.GetData()))
	}

	data, err := uim.processor.GetCompressedData()
	uim.cache.Set(fileUri.String(), data, 1)
	if err != nil {
		return "", err
	}
	uim.saveFile(compressedFileName, bytes.NewReader(data))

	return fileUri.String(), nil
}

func (uim *UploadImageManager) GetHexSum() string {
	return uim.processor.GetHexSum()
}

func (uim *UploadImageManager) GetImageInfo() (ImageInfo, error) {
	imgInfo, err := uim.processor.GetImageInfo()
	if err != nil {
		return ImageInfo{}, err
	}
	imgInfo.Name = uim.source.GetName()
	return *imgInfo, nil
}
