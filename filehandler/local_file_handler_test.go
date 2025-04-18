package filehandler

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/follow1123/photos/logger"
	"github.com/stretchr/testify/suite"
)

type LocalFileHandlerTestSuite struct {
	suite.Suite
	handler *LocalFileHandler
	dataDir string
	logger logger.AppLogger
	testFilePath string
}

func TestLocalFileHandlerTestSuite(t *testing.T) {
	suite.Run(t, &LocalFileHandlerTestSuite{})
}

func (s *LocalFileHandlerTestSuite) SetupSuite() {
	dir, err := os.Getwd()
	s.Nil(err, "test directory init failed")
	s.dataDir = filepath.Join(dir, "files")

	err = os.Mkdir(s.dataDir, 0755)
	s.Nil(err, "test directory init failed")
	s.testFilePath = filepath.Join(s.dataDir, "test_file.txt")
	f, err := os.Create(s.testFilePath)
	s.Nil(err, "test file create failed")
	_, err = f.WriteString("1234321uqweripwejifqwe")
	s.Nil(err, "test file write failed")

	baseLogger := logger.InitBaseLogger()
	s.logger = logger.NewAppLogger(baseLogger)
	s.handler = NewLocalFileHandler(s.dataDir, s.logger)
}

func (s *LocalFileHandlerTestSuite) TearDownSuite() {
	err := os.RemoveAll(s.dataDir)
	s.Nil(err, "test directory clean failed")
}

func (s *LocalFileHandlerTestSuite) TestNextFileName() {
	filePath := s.handler.nextFilePath()
	s.True(strings.HasPrefix(filePath, s.dataDir))
}

func (s *LocalFileHandlerTestSuite) TestSaveSuccess() {
	src, err := os.Open(s.testFilePath)
	s.Nil(err)

	path, err := s.handler.Save("", src)
	s.Nil(err)

	s.True(strings.HasPrefix(path, LOCAL_FILE_PREFIX))
}
