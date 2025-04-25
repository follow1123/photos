package imagemanager

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type FileUriTestSuite struct {
	suite.Suite
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, &FileUriTestSuite{})
}

func (s *FileUriTestSuite) SetupSuite() {
}

func (s *FileUriTestSuite) TestNewFileUriSuccess() {
	var filesRoot = "/a/b/c"

	scenarios := []struct {
		uri string
	}{
		{"local://fdskfaj"},
		{"ftp://fdskfaj"},
		{"scp://fdskfaj"},
	}

	for _, scenario := range scenarios {
		s.NotPanics(func() { NewFileUri(filesRoot, scenario.uri) })
	}
}

func (s *FileUriTestSuite) TestNewFileUriFailure() {
	var filesRoot = "/a/b/c"

	scenarios := []struct {
		uri string
	}{
		{"http://fdskfaj"},
		{"abc://fdskfaj"},
		{"fdskfaj"},
	}

	for _, scenario := range scenarios {
		s.Panics(func() { NewFileUri(filesRoot, scenario.uri) })
	}
}

func (s *FileUriTestSuite) TestCreateRemoteFileUriSuccess() {
	var filesRoot = "/a/b/c"

	scenarios := []struct {
		uri string
	}{
		{"ftp://fdskfaj"},
		{"scp://fdskfaj"},
	}

	for _, scenario := range scenarios {
		fileUri, err := CreateRemoteFileUri(filesRoot, scenario.uri)
		s.Nil(err)
		s.NotNil(fileUri)
	}
}

func (s *FileUriTestSuite) TestCreateRemoteFileUriFailure() {
	var filesRoot = "/a/b/c"

	scenarios := []struct {
		uri string
	}{
		{"local://fdskfaj"},
		{"124234"},
	}

	for _, scenario := range scenarios {
		fileUri, err := CreateRemoteFileUri(filesRoot, scenario.uri)
		s.Nil(fileUri)
		s.Equal(ErrUnsupportedRemoteFiles, err)
	}
}

func (s *FileUriTestSuite) TestCreateLocalFileUriStructure() {
	var filesRoot = "/a/b/c"
	var originalSuffix = "_original"
	var compressedSuffix = "_compressed"
	fileUri := CreateLocalFileUri(filesRoot)

	s.Equal(LOCAL_FILE, fileUri.fileType)
	s.Equal(filesRoot, fileUri.filesRoot)
	s.True(strings.HasPrefix(fileUri.filePath, filesRoot))
	s.True(strings.HasPrefix(fileUri.uri, fileUri.fileType))
	s.True(strings.HasSuffix(fileUri.GetOriginalFilePath(), originalSuffix))
	s.True(strings.HasSuffix(fileUri.GetCompressedFilePath(), compressedSuffix))
}

func (s *FileUriTestSuite) TestCreateRemoteFileUriStructureSuccess() {
	var filesRoot = "/a/b/c"
	var compressedSuffix = "_compressed"

	scenarios := []struct {
		remoteUri        string
		expectedFileType string
	}{
		{"ftp://localhost:1234/a/b/c", FTP_FILE},
		{"scp://za@localhost:5678/a/b/c", SCP_FILE},
	}

	for _, scenario := range scenarios {
		fileUri, err := CreateRemoteFileUri(filesRoot, scenario.remoteUri)
		s.Nil(err)
		s.NotNil(fileUri)

		s.Equal(scenario.expectedFileType, fileUri.fileType)
		s.Equal(filesRoot, fileUri.filesRoot)
		s.True(strings.HasPrefix(fileUri.filePath, filesRoot))
		s.True(strings.HasPrefix(fileUri.uri, fileUri.fileType))
		s.Equal(scenario.remoteUri, fileUri.GetOriginalFilePath())
		s.True(strings.HasSuffix(fileUri.GetCompressedFilePath(), compressedSuffix))
	}
}

func (s *FileUriTestSuite) TestNewFileUriStructure() {
	var filesRoot = "/a/b/c"

	localFileUri := CreateLocalFileUri(filesRoot)
	remoteFileUri, err := CreateRemoteFileUri(filesRoot, "scp://za@localhost:5678/a/b/c")
	s.NotNil(remoteFileUri)
	s.Nil(err)

	scenarios := []struct {
		uri                        string
		expectedFileType           string
		expectedFilePath           string
		expectedOriginalFilePath   string
		expectedCompressedFilePath string
	}{
		{localFileUri.String(), localFileUri.fileType, localFileUri.filePath, localFileUri.GetOriginalFilePath(), localFileUri.GetCompressedFilePath()},
		{remoteFileUri.String(), remoteFileUri.fileType, remoteFileUri.filePath, remoteFileUri.GetOriginalFilePath(), remoteFileUri.GetCompressedFilePath()},
	}

	for _, scenario := range scenarios {
		fileUri := NewFileUri(filesRoot, scenario.uri)
		s.Equal(scenario.expectedFileType, fileUri.fileType)
		s.Equal(filesRoot, fileUri.filesRoot)
		s.Equal(scenario.expectedFilePath, fileUri.filePath)
		s.Equal(scenario.expectedOriginalFilePath, fileUri.GetOriginalFilePath())
		s.Equal(scenario.expectedCompressedFilePath, fileUri.GetCompressedFilePath())
	}
}
