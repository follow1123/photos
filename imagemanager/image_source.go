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

func NewMultipartSource(fileHeader *multipart.FileHeader) ImageSource {
	return &MultipartSource{FileHeader: fileHeader}
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

func NewRemoteUriSource(uri string) ImageSource {
	// TODO
	return nil
}

type ReaderSource struct {
	reader io.Reader
	name   string
}

func NewReaderSource(r io.Reader, name string) ImageSource {
	return &ReaderSource{reader: r, name: name}
}

func (rs *ReaderSource) GetReader() (io.ReadCloser, error) {
	return io.NopCloser(rs.reader), nil
}

func (rs *ReaderSource) GetName() string {
	return rs.name
}
