package imagemanager

import (
	"encoding/base64"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	LOCAL_FILE string = "local://"
	FTP_FILE          = "ftp://"
	SCP_FILE          = "scp://"
)

var ErrInvalidFileType = errors.New("invalid file type")
var ErrUnsupportedRemoteFiles = errors.New("unsupported remote files")

type FileUri struct {
	uri       string
	fileType  string
	filePath  string
	filesRoot string
}

func newFileUri(filesRoot string, uri string) *FileUri {
	var fileType string
	if strings.HasPrefix(uri, LOCAL_FILE) {
		fileType = LOCAL_FILE
	} else if strings.HasPrefix(uri, FTP_FILE) {
		fileType = FTP_FILE
	} else if strings.HasPrefix(uri, SCP_FILE) {
		fileType = SCP_FILE
	} else {
		panic(ErrInvalidFileType)
	}
	return &FileUri{
		uri:       uri,
		fileType:  fileType,
		filePath:  strings.Replace(uri, fileType, filesRoot, 1),
		filesRoot: filesRoot,
	}
}
func createLocalFileUri(filesRoot string) *FileUri {
	var fileType = LOCAL_FILE
	t := time.Now()
	id := strings.ReplaceAll(uuid.New().String(), "-", "")
	dir := filepath.Join(filesRoot, t.Format("200601/02/15"))
	timestamp := t.Format("20060102150405")
	fileName := fmt.Sprintf("%s_%s", timestamp, id)
	filePath := filepath.Join(dir, fileName)

	return &FileUri{
		uri:       strings.Replace(filePath, filesRoot, fileType, 1),
		fileType:  fileType,
		filePath:  filePath,
		filesRoot: filesRoot,
	}
}

func createRemoteFileUri(filesRoot string, remoteUri string) (*FileUri, error) {
	var fileType string
	if strings.HasPrefix(remoteUri, FTP_FILE) {
		fileType = FTP_FILE
	} else if strings.HasPrefix(remoteUri, SCP_FILE) {
		fileType = SCP_FILE
	} else {
		return nil, ErrUnsupportedRemoteFiles
	}

	t := time.Now()
	encodedUri := base64.StdEncoding.EncodeToString([]byte(remoteUri))
	dir := filepath.Join(filesRoot, t.Format("200601/02/15"))
	timestamp := t.Format("20060102150405")
	fileName := fmt.Sprintf("%s_%s", timestamp, encodedUri)

	filePath := filepath.Join(dir, fileName)
	return &FileUri{
		uri:       strings.Replace(filePath, filesRoot, fileType, 1),
		fileType:  fileType,
		filePath:  filePath,
		filesRoot: filesRoot,
	}, nil
}

func (fu *FileUri) String() string {
	return fu.uri
}

func (fu *FileUri) Is(fileType string) bool {
	return fu.fileType == fileType
}

func (fu *FileUri) GetOriginalFilePath() string {
	if fu.fileType == LOCAL_FILE {
		return fmt.Sprintf("%s_original", fu.filePath)
	} else {
		baseName := filepath.Base(fu.filePath)
		idx := strings.Index(baseName, "_") + 1
		encodedUri := baseName[idx:]
		uri, err := base64.StdEncoding.DecodeString(encodedUri)
		if err != nil {
			panic(err)
		}
		return string(uri)
	}
}

func (fu *FileUri) GetCompressedFilePath() string {
	return fmt.Sprintf("%s_compressed", fu.filePath)
}
