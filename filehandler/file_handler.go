package filehandler

import (
	"io"
	"strings"
)

type FileHandler interface {
	Open(uri string) (io.ReadCloser, error)
	Save(uri string, r io.Reader) (string, error)
	Delete(uri string) error
	Exists(uri string) (bool, error)
}

type FileHandlerFactory struct {
	localFileHandler *LocalFileHandler
}

func (f *FileHandlerFactory) GetHandler(uri string) FileHandler {
	uri = strings.TrimSpace(uri)
	if uri == "" {
		return f.localFileHandler
	}

	if strings.HasPrefix(uri, LOCAL_FILE_PREFIX) {
		return f.localFileHandler
	}

	return nil
}
