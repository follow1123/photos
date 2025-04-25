package imagegen_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/follow1123/photos/generator/imagegen"
	"github.com/stretchr/testify/suite"
)

type ImageGeneratorTestSuite struct {
	suite.Suite
}

func TestImageGeneratorTestSuite(t *testing.T) {
	suite.Run(t, &ImageGeneratorTestSuite{})
}

func (s *ImageGeneratorTestSuite) TestGenImage() {
	for i := range 10 {
		_ = i
		buf := new(bytes.Buffer)
		imgInfo, err := imagegen.GenImage(buf)
		s.Nil(err)
		s.False(imgInfo.Format == "")
		s.False(imgInfo.Height == 0)
		s.False(imgInfo.Width == 0)
		s.False(len(buf.Bytes()) == 0)
	}
}

// 查看生成的随机图片的效果，方法名第一个改成大写执行
// 注释掉 RemoveAll 查看
func (s *ImageGeneratorTestSuite) testGenImageFile() {
	cwd, err := os.Getwd()
	s.Nil(err)
	imageDir := filepath.Join(cwd, fmt.Sprintf("image_test"))
	err = os.MkdirAll(imageDir, 0755)
	s.Nil(err)

	for i := range 10 {
		buf := new(bytes.Buffer)
		imgInfo, err := imagegen.GenImage(buf)
		s.Nil(err)

		imageFile, err := os.Create(filepath.Join(imageDir, fmt.Sprintf("image_%d.%s", i, imgInfo.Format)))
		s.Nil(err)
		defer imageFile.Close()
		io.Copy(imageFile, bytes.NewReader(buf.Bytes()))
	}

	err = os.RemoveAll(imageDir)
	s.Nil(err)
	// 防止测试缓存
	s.False(true)
}
