package imagemanager_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/follow1123/photos/generator/appgen"
	"github.com/follow1123/photos/generator/imagegen"
	"github.com/follow1123/photos/imagemanager"
	"github.com/follow1123/photos/logger"
	"github.com/stretchr/testify/suite"
)

type ImageProcessorTestSuite struct {
	suite.Suite
	logger *logger.AppLogger
}

func TestImageProcessorTestSuite(t *testing.T) {
	suite.Run(t, &ImageProcessorTestSuite{})
}

func (s *ImageProcessorTestSuite) SetupSuite() {
	appComponents := &appgen.AppComponents{}
	appLogger, err := appgen.GenAppLogger(appComponents)
	s.Nil(err)
	s.logger = appLogger
}

func (s *ImageProcessorTestSuite) TestGetData() {
	buf := new(bytes.Buffer)
	_, err := imagegen.GenImage(buf)
	s.Nil(err)
	expectedData := buf.Bytes()

	rc := io.NopCloser(bytes.NewReader(expectedData))

	ip := imagemanager.NewImageProcessor(rc, s.logger)

	data, err := ip.GetData()
	s.Nil(err)
	s.Equal(expectedData, data)
}

func (s *ImageProcessorTestSuite) TestGetImageInfo() {
	buf := new(bytes.Buffer)
	expectedImgInfo, err := imagegen.GenImage(buf)
	s.Nil(err)

	rc := io.NopCloser(bytes.NewReader(buf.Bytes()))
	ip := imagemanager.NewImageProcessor(rc, s.logger)
	imgInfo, err := ip.GetImageInfo()
	s.Nil(err)
	s.Equal(expectedImgInfo.Format, imgInfo.Format)
	s.Equal(int64(expectedImgInfo.Width), imgInfo.Width)
	s.Equal(int64(expectedImgInfo.Height), imgInfo.Height)
}

func (s *ImageProcessorTestSuite) TestGetCompressedData() {
	buf := new(bytes.Buffer)
	// TODO：test other image format
	_, err := imagegen.GenImage(buf, imagegen.WithFormat(imagegen.FORMAT_JPEG))
	s.Nil(err)

	rc := io.NopCloser(bytes.NewReader(buf.Bytes()))
	ip := imagemanager.NewImageProcessor(rc, s.logger)
	data, err := ip.GetCompressedData()
	s.Nil(err)
	s.NotEqual(buf.Bytes(), data)
}
