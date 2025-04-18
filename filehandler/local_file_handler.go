package filehandler

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/follow1123/photos/logger"
	"github.com/google/uuid"
)

const LOCAL_FILE_PREFIX = "local://"

type LocalFileHandler struct {
	filesRoot string
	logger logger.AppLogger
}

func NewLocalFileHandler(filesPath string, logger logger.AppLogger) *LocalFileHandler {
	return &LocalFileHandler{filesRoot: filesPath, logger: logger}
}

func (l *LocalFileHandler) checkUri(uri string) (string, error) {
	if strings.HasPrefix(uri, LOCAL_FILE_PREFIX) {
		msg := fmt.Sprintf("invalid local file uri: %s", uri)
		l.logger.Error(msg)
		return "", errors.New(msg)
	}
	return strings.TrimPrefix(uri, LOCAL_FILE_PREFIX), nil
}

// 获取下一个文件保存的路径
// 如果路径不存在，需要手动创建
func (l *LocalFileHandler) nextFilePath() string {
	t := time.Now()
	id := strings.ReplaceAll(uuid.New().String(), "-", "")
	dir := filepath.Join(l.filesRoot, t.Format("200601/02/15"))
	timestamp := t.Format("20060102150405")
	fileName := fmt.Sprintf("%s_%s", id, timestamp)
	return filepath.Join(dir, fileName)
}

func (l *LocalFileHandler) Save(_ string, src io.Reader) (string, error) {
	filePath := l.nextFilePath()
	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		l.logger.Error("cannot create next files directory: %s", filepath.Dir(filePath))
		return "", nil
	}

	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		return "", err
	}

	return strings.Replace(filePath, l.filesRoot, LOCAL_FILE_PREFIX, 1), nil
}

func (l *LocalFileHandler) Open(uri string) (io.ReadCloser, error) {
	filePath, err := l.checkUri(uri)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (l *LocalFileHandler) Delete(uri string) error {
	filePath, err := l.checkUri(uri)
	if err != nil {
		return err
	}

	err = os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}

func (l *LocalFileHandler) Exists(uri string) (bool, error) {
	filePath, err := l.checkUri(uri)
	if err != nil {
		return false, err
	}

	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}
	return info.Mode().IsRegular(), nil
}
