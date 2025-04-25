package imagemanager

import "errors"

var ErrUnsupportedImageFormat = errors.New("Unsupported image format")
var ErrInvalidFileType = errors.New("invalid file type")
var ErrUnsupportedRemoteFiles = errors.New("unsupported remote file")
