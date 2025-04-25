package imagemanager_test

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"image"
	"os"
	"testing"

	"github.com/follow1123/photos/config"
	"github.com/follow1123/photos/generator/appgen"
	"github.com/follow1123/photos/generator/imagegen"
	"github.com/follow1123/photos/imagemanager"
	"github.com/follow1123/photos/logger"
	"github.com/stretchr/testify/suite"
)

type UploadImageManagerTestSuite struct {
	suite.Suite
	conf   *config.Config
	logger *logger.AppLogger
	cache  *imagemanager.ImageCache
}

func TestUploadImageManagerTestSuite(t *testing.T) {
	suite.Run(t, &UploadImageManagerTestSuite{})
}

func (s *UploadImageManagerTestSuite) SetupSuite() {
	appComponents := &appgen.AppComponents{}
	conf, err := appgen.GenConfig(appComponents)
	s.Nil(err)
	conf.CreatePath()
	appLogger, err := appgen.GenAppLogger(appComponents)
	s.Nil(err)
	imageCache, err := appgen.GenImageCache(appComponents)
	s.Nil(err)

	s.conf = conf
	s.logger = appLogger
	s.cache = imageCache
}

func (s *UploadImageManagerTestSuite) TearDownSuite() {
	s.cache.Close()
	s.conf.DeletePath()
}

func (s *UploadImageManagerTestSuite) SetupTest() {
	s.cache.Clear()
}

func (s *UploadImageManagerTestSuite) TestGetImageName() {
	data := []byte("hello world")
	expectedName := "aaa"

	uploadMgr := imagemanager.NewUploadImageManager(
		s.conf.GetFilesPath(),
		imagemanager.NewReaderSource(bytes.NewReader(data), expectedName),
		s.logger,
		s.cache,
	)
	s.Equal(expectedName, uploadMgr.GetImageName())
}

func (s *UploadImageManagerTestSuite) TestGetHexSumSuccess() {
	name := "aaa"

	type Scenario = struct {
		source         imagemanager.ImageSource
		expectedHexSum string
	}

	newScenario := func() Scenario {
		buf := new(bytes.Buffer)
		_, err := imagegen.GenImage(buf)
		s.Nil(err)

		source := imagemanager.NewReaderSource(bytes.NewReader(buf.Bytes()), name)
		sum := md5.Sum(buf.Bytes())
		return Scenario{
			source:         source,
			expectedHexSum: hex.EncodeToString(sum[:]),
		}
	}

	scenarios := []Scenario{
		newScenario(),
		newScenario(),
		newScenario(),
	}

	for _, scenario := range scenarios {
		uploadMgr := imagemanager.NewUploadImageManager(
			s.conf.GetFilesPath(),
			scenario.source,
			s.logger,
			s.cache,
		)

		actualHexSum, err := uploadMgr.GetHexSum()
		s.Nil(err)
		s.Equal(scenario.expectedHexSum, actualHexSum)
	}
}

func (s *UploadImageManagerTestSuite) TestGetHexSumFailure() {
	name := "aaa"
	data := []byte("24123423")

	uploadMgr := imagemanager.NewUploadImageManager(
		s.conf.GetFilesPath(),
		imagemanager.NewReaderSource(bytes.NewReader(data), name),
		s.logger,
		s.cache,
	)

	actualHexSum, err := uploadMgr.GetHexSum()
	s.Equal("", actualHexSum)
	s.Equal(image.ErrFormat, err)
}

func (s *UploadImageManagerTestSuite) TestGetImageInfoSuccess() {
	name := "aaa"

	type Scenario = struct {
		source         imagemanager.ImageSource
		expectedSize   int64
		expectedFormat string
		expectedWidth  int64
		expectedHeight int64
	}

	newScenario := func() Scenario {
		buf := new(bytes.Buffer)
		imgInfo, err := imagegen.GenImage(buf)
		s.Nil(err)

		source := imagemanager.NewReaderSource(bytes.NewReader(buf.Bytes()), name)
		return Scenario{
			source:         source,
			expectedSize:   int64(len(buf.Bytes())),
			expectedFormat: imgInfo.Format,
			expectedWidth:  int64(imgInfo.Width),
			expectedHeight: int64(imgInfo.Height),
		}
	}

	scenarios := []Scenario{
		newScenario(),
		newScenario(),
		newScenario(),
	}

	for _, scenario := range scenarios {
		uploadMgr := imagemanager.NewUploadImageManager(
			s.conf.GetFilesPath(),
			scenario.source,
			s.logger,
			s.cache,
		)

		imgInfo, err := uploadMgr.GetImageInfo()
		s.Nil(err)
		s.Equal(scenario.expectedSize, imgInfo.Size)
		s.Equal(scenario.expectedFormat, imgInfo.Format)
		s.Equal(scenario.expectedWidth, imgInfo.Width)
		s.Equal(scenario.expectedHeight, imgInfo.Height)
	}
}

func (s *UploadImageManagerTestSuite) TestGetImageInfoFailure() {
	name := "aaa"
	data := []byte("24123423")

	uploadMgr := imagemanager.NewUploadImageManager(
		s.conf.GetFilesPath(),
		imagemanager.NewReaderSource(bytes.NewReader(data), name),
		s.logger,
		s.cache,
	)

	actualImageInfo, err := uploadMgr.GetImageInfo()
	s.Nil(actualImageInfo)
	s.Equal(image.ErrFormat, err)
}

func (s *UploadImageManagerTestSuite) TestSaveSuccess() {
	name := "aaa"
	filesRoot := s.conf.GetFilesPath()

	type Scenario = struct {
		source imagemanager.ImageSource
	}

	newScenario := func() Scenario {
		buf := new(bytes.Buffer)
		_, err := imagegen.GenImage(buf)
		s.Nil(err)
		source := imagemanager.NewReaderSource(bytes.NewReader(buf.Bytes()), name)
		return Scenario{
			source: source,
		}
	}

	scenarios := []Scenario{
		newScenario(),
		newScenario(),
		newScenario(),
	}

	for _, scenario := range scenarios {
		uploadMgr := imagemanager.NewUploadImageManager(
			filesRoot,
			scenario.source,
			s.logger,
			s.cache,
		)

		uri, err := uploadMgr.Save()
		s.Nil(err)
		s.NotPanics(func() {
			fileUri := imagemanager.NewFileUri(filesRoot, uri)
			compressedFilePath := fileUri.GetCompressedFilePath()
			info, err := os.Stat(compressedFilePath)
			s.False(os.IsNotExist(err))
			s.False(info.IsDir())
		})
	}
}
