package imagemanager

import "os"

type DeleteImageManager struct {
	uri FileUri
}

func NewDeleteImageManager(filesRoot string, uri string) *DeleteImageManager {
	return &DeleteImageManager{
		uri: *NewFileUri(filesRoot, uri),
	}
}

func (dim *DeleteImageManager) Delete() error {
	originalFilePath := dim.uri.GetOriginalFilePath()
	compressedFilePath := dim.uri.GetCompressedFilePath()

	if dim.uri.Is(LOCAL_FILE) {
		if err := os.Remove(originalFilePath); err != nil {
			return err
		}
	}
	if err := os.Remove(compressedFilePath); err != nil {
		return err
	}
	return nil
}
