package generator

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"math/rand/v2"
)

const (
	FORMAT_PNG  = "png"
	FORMAT_JPEG = "jpeg"
)

var (
	COLOR_TRANSPARENT = color.RGBA{0, 0, 0, 0}
	COLOR_BLACK       = color.RGBA{0, 0, 0, 255}
	COLOR_WHITE       = color.RGBA{255, 255, 255, 255}
	COLOR_RED         = color.RGBA{255, 0, 0, 255}
	COLOR_GREEN       = color.RGBA{0, 255, 0, 255}
	COLOR_BLUE        = color.RGBA{0, 0, 255, 255}
	COLOR_YELLOW      = color.RGBA{255, 255, 0, 255}
	COLOR_WHEAT       = color.RGBA{255, 231, 186, 255}
	COLOR_PURPLE      = color.RGBA{160, 32, 240, 255}
	COLOR_GREY        = color.RGBA{181, 181, 181, 255}
	COLOR_CORAL       = color.RGBA{255, 114, 86, 255}
	COLOR_DARKCYAN    = color.RGBA{0, 139, 139, 255}
	COLOR_LIGHTCYAN   = color.RGBA{209, 238, 238, 255}
	COLOR_MAROON      = color.RGBA{176, 48, 96, 255}
	COLOR_ORANGE      = color.RGBA{255, 165, 0, 255}
)

var (
	STYLE_PURE_COLOR_RED   = &StylePureColor{color: COLOR_RED}
	STYLE_PURE_COLOR_GREEN = &StylePureColor{color: COLOR_GREEN}
	STYLE_PURE_COLOR_BLUE  = &StylePureColor{color: COLOR_BLUE}

	STYLE_CHESS_BOARD_1 = &StyleChessboard{color1: COLOR_BLACK, color2: COLOR_WHITE}
	STYLE_CHESS_BOARD_2 = &StyleChessboard{color1: COLOR_ORANGE, color2: COLOR_YELLOW}
	STYLE_CHESS_BOARD_3 = &StyleChessboard{color1: COLOR_CORAL, color2: COLOR_RED}

	STYLE_GRADIENT = &StyleGradient{}

	STYLE_RANDOM_NOISE = &StyleRandomNoise{}

	STYLE_LINEAR_GRADIENT = &StyleLinearGradient{}

	STYLE_CIRCULAR_1 = &StyleCircular{innerColor: COLOR_RED, outerColor: COLOR_TRANSPARENT}
	STYLE_CIRCULAR_2 = &StyleCircular{innerColor: COLOR_GREY, outerColor: COLOR_MAROON}
	STYLE_CIRCULAR_3 = &StyleCircular{innerColor: COLOR_DARKCYAN, outerColor: COLOR_BLUE}
)

type ImageStyle interface {
	Render(img *image.RGBA, width int, height int)
}

type StylePureColor struct {
	color color.RGBA
}

func (spc *StylePureColor) Render(img *image.RGBA, width int, height int) {
	for x := range width {
		for y := range height {
			img.Set(x, y, spc.color)
		}
	}
}

type StyleChessboard struct {
	color1 color.RGBA
	color2 color.RGBA
}

func (sc *StyleChessboard) Render(img *image.RGBA, width int, height int) {
	for x := range width {
		for y := range height {
			if (x/10+y/10)%2 == 0 {
				img.Set(x, y, sc.color1)
			} else {
				img.Set(x, y, sc.color2)
			}
		}
	}
}

type StyleGradient struct {
	// color color.RGBA
}

func (sg *StyleGradient) Render(img *image.RGBA, width int, height int) {
	for x := range width {
		for y := range height {
			r := uint8((x * 255) / (width - 1))
			b := 255 - r
			img.Set(x, y, color.RGBA{r, 0, b, 255})
		}
	}
}

type StyleCircular struct {
	innerColor color.RGBA
	outerColor color.RGBA
}

func (sc *StyleCircular) Render(img *image.RGBA, width int, height int) {
	for x := range width {
		for y := range height {
			cx, cy := width/2, height/2
			r := float64(min(width, height)) * 0.3
			dx := float64(x - cx)
			dy := float64(y - cy)
			if dx*dx+dy*dy < r*r {
				img.Set(x, y, sc.innerColor)
			} else {
				img.Set(x, y, sc.outerColor)
			}
		}
	}
}

type StyleRandomNoise struct {
}

func (_ *StyleRandomNoise) Render(img *image.RGBA, width int, height int) {
	for x := range width {
		for y := range height {
			img.Set(x, y, color.RGBA{
				R: uint8(rand.IntN(256)),
				G: uint8(rand.IntN(256)),
				B: uint8(rand.IntN(256)),
				A: 255,
			})
		}
	}
}

type StyleLinearGradient struct {
	// color color.RGBA
}

func (slg *StyleLinearGradient) Render(img *image.RGBA, width int, height int) {
	for x := range width {
		for y := range height {
			c := uint8((x + y) * 255 / 198)
			img.Set(x, y, color.RGBA{c, c / 2, 255 - c, 255})
		}
	}
}

func GenImage(format string, style ImageStyle, writer io.Writer) error {
	switch format {
	case FORMAT_PNG, FORMAT_JPEG:
	default:
		return fmt.Errorf("invalid image format: %s", format)
	}
	var width = 300
	var height = 300
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	style.Render(img, width, height)

	var err error
	switch format {
	case FORMAT_PNG:
		err = png.Encode(writer, img)
	case FORMAT_JPEG:
		err = jpeg.Encode(writer, img, nil)
	}
	if err != nil {
		return err
	}
	return nil
}
