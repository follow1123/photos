package generator_test

import (
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"testing"

	"github.com/follow1123/photos/generator"
	"github.com/stretchr/testify/suite"
)

type ImageGeneratorTestSuite struct {
	suite.Suite
}

func TestImageGeneratorTestSuite(t *testing.T) {
	suite.Run(t, &ImageGeneratorTestSuite{})
}

func (s *ImageGeneratorTestSuite) TestGenImage() {
	var dir = "/home/yf/space/tmp/photos/"

	styleList := [...]generator.ImageStyle{
		generator.STYLE_CHESS_BOARD_2,
		generator.STYLE_PURE_COLOR_BLUE,
		generator.STYLE_CIRCULAR_2,
		generator.STYLE_RANDOM_NOISE,
		generator.STYLE_CIRCULAR_3,
		generator.STYLE_LINEAR_GRADIENT,
	}

	formatList := [...]string{generator.FORMAT_PNG, generator.FORMAT_JPEG}

	styleCount := len(styleList)
	formatCount := len(formatList)

	for i := range 10 {
		styleIdx := rand.IntN(styleCount)
		formatIdx := rand.IntN(formatCount)

		// idx := rand.IntN(styleCount)
		format := formatList[formatIdx]
		imageFile, err := os.Create(filepath.Join(dir, fmt.Sprintf("image_%d.%s", i, format)))
		s.Nil(err)
		defer imageFile.Close()
		// fmt.Printf("idx: %d\n", idx)
		err = generator.GenImage(format, styleList[styleIdx], imageFile)
		s.Nil(err)
	}
	s.True(false)
}
