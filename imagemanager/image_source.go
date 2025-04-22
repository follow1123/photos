package imagemanager

import (
	"io"
	"mime/multipart"
)

type ImageSource interface {
	GetReader() (io.ReadCloser, error)
	GetName() string
}

type MultipartSource struct {
	FileHeader *multipart.FileHeader
}

func (ms *MultipartSource) GetReader() (io.ReadCloser, error) {
	file, err := ms.FileHeader.Open()
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (ms *MultipartSource) GetName() string {
	return ms.FileHeader.Filename
}
