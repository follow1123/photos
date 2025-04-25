package imagegen

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"math/rand/v2"
)

const (
	FORMAT_PNG  = "png"
	FORMAT_JPEG = "jpeg"

	IMG_MIN_WIDTH  = 1280
	IMG_MIN_HEIGHT = 800
	IMG_MAX_WIDTH  = 1920
	IMG_MAX_HEIGHT = 1080
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

var FixedColors = []color.RGBA{
	COLOR_TRANSPARENT, COLOR_BLACK, COLOR_WHITE, COLOR_RED, COLOR_GREEN, COLOR_BLUE,
	COLOR_YELLOW, COLOR_WHEAT, COLOR_PURPLE, COLOR_GREY, COLOR_CORAL, COLOR_DARKCYAN,
	COLOR_LIGHTCYAN, COLOR_MAROON, COLOR_ORANGE,
}

var FixedFormats = []string{FORMAT_PNG, FORMAT_JPEG}

var FixedStyles = []Randomize{
	&StylePureColor{},
	&StyleCircular{},
	&StyleGradient{},
	&StyleChessboard{},
	&StyleRandomNoise{},
	&StyleLinearGradient{},
}

type ImageStyle interface {
	Render(img *image.RGBA, width int, height int)
}

type Randomize interface {
	Randomize() ImageStyle
}

type StylePureColor struct {
	color color.RGBA
}

func (spc *StylePureColor) Randomize() ImageStyle {
	spc.color = RandomRGBAColor()
	return spc
}

func (spc *StylePureColor) Render(img *image.RGBA, width int, height int) {
	for x := range width {
		for y := range height {
			img.Set(x, y, spc.color)
		}
	}
}

type StyleChessboard struct {
	color1   color.RGBA
	color2   color.RGBA
	cellSize int
}

func (sc *StyleChessboard) Randomize() ImageStyle {
	sc.color1 = RandomRGBAColor()
	sc.color2 = RandomRGBAColor()
	sc.cellSize = rand.IntN(70) + 15
	return sc
}

func (sc *StyleChessboard) Render(img *image.RGBA, width int, height int) {
	for x := range width {
		for y := range height {
			if (x/sc.cellSize+y/sc.cellSize)%2 == 0 {
				img.Set(x, y, sc.color1)
			} else {
				img.Set(x, y, sc.color2)
			}
		}
	}
}

type StyleGradient struct {
	color1 color.RGBA
	color2 color.RGBA
	angle  float64 // 渐变的角度
}

func (sg *StyleGradient) Randomize() ImageStyle {
	sg.angle = rand.Float64() * 360
	sg.color1 = RandomRGBAColor()
	sg.color2 = RandomRGBAColor()
	return sg
}

func (sg *StyleGradient) Render(img *image.RGBA, width int, height int) {
	// 使用角度来决定渐变方向
	angleRad := sg.angle * math.Pi / 180.0 // 转换为弧度
	cosAngle := math.Cos(angleRad)
	sinAngle := math.Sin(angleRad)

	// 渐变的颜色过渡
	for x := range width {
		for y := range height {
			// 计算每个点到渐变起始点的距离
			t := float64(x)*cosAngle + float64(y)*sinAngle
			t = (t + float64(width+height)) / float64(width+height) // 将 t 标准化为 [0, 1]

			// 渐变颜色计算：基于 t 计算颜色的过渡
			r := uint8(float64(sg.color1.R)*(1-t) + float64(sg.color2.R)*t)
			g := uint8(float64(sg.color1.G)*(1-t) + float64(sg.color2.G)*t)
			b := uint8(float64(sg.color1.B)*(1-t) + float64(sg.color2.B)*t)
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}
}

type StyleCircular struct {
	innerColor color.RGBA
	outerColor color.RGBA
}

func (sc *StyleCircular) Randomize() ImageStyle {
	sc.innerColor = RandomRGBAColor()
	sc.outerColor = RandomRGBAColor()
	return sc
}

func (sc *StyleCircular) Render(img *image.RGBA, width int, height int) {
	// 随机生成圆心的位置，保证圆形不超出图像边界
	maxRadius := float64(min(width, height)) * 0.3               // 最大半径
	radius := float64(rand.IntN(int(maxRadius))) + maxRadius*0.3 // 随机半径

	// 随机生成圆心位置，确保圆形不会超出图像
	cx := rand.IntN(width-int(2*radius)) + int(radius)
	cy := rand.IntN(height-int(2*radius)) + int(radius)

	for x := range width {
		for y := range height {
			// 计算 (x, y) 是否在圆形范围内
			dx := float64(x - cx)
			dy := float64(y - cy)
			if dx*dx+dy*dy <= radius*radius {
				img.Set(x, y, sc.innerColor)
			} else {
				img.Set(x, y, sc.outerColor)
			}
		}
	}
}

type StyleRandomNoise struct {
}

func (sc *StyleRandomNoise) Randomize() ImageStyle {
	// do nothing
	return sc
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
	color1 color.RGBA
	color2 color.RGBA
	angle  float64 // 渐变的角度
	steps  int     // 渐变的步长
}

func (slg *StyleLinearGradient) Randomize() ImageStyle {
	// 随机选择一个角度和步长
	slg.angle = rand.Float64() * 360
	slg.steps = rand.IntN(10) + 2 // 随机步长，2到12之间
	slg.color1 = RandomRGBAColor()
	slg.color2 = RandomRGBAColor()
	return slg
}

func (slg *StyleLinearGradient) Render(img *image.RGBA, width int, height int) {
	// 使用角度来决定渐变方向
	angleRad := slg.angle * math.Pi / 180.0
	cosAngle := math.Cos(angleRad)
	sinAngle := math.Sin(angleRad)

	// 渐变的颜色过渡
	for x := range width {
		for y := range height {
			// 计算每个点到渐变起始点的距离
			t := float64(x)*cosAngle + float64(y)*sinAngle
			t = (t + float64(width+height)) / float64(width+height)   // 标准化为 [0, 1]
			t = math.Floor(t*float64(slg.steps)) / float64(slg.steps) // 步长控制

			// 渐变颜色计算，基于 t 计算颜色的过渡
			r := uint8(float64(slg.color1.R)*(1-t) + float64(slg.color2.R)*t)
			g := uint8(float64(slg.color1.G)*(1-t) + float64(slg.color2.G)*t)
			b := uint8(float64(slg.color1.B)*(1-t) + float64(slg.color2.B)*t)
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}
}

type Option interface {
	apply(*ImageGenOption)
}

type optionFunc func(*ImageGenOption)

func (f optionFunc) apply(opts *ImageGenOption) {
	f(opts)
}

func WithSize(width int, height int) Option {
	if width < IMG_MIN_WIDTH || width > IMG_MAX_WIDTH {
		width = IMG_MIN_WIDTH // 设置为默认最小值
	}
	if height < IMG_MIN_HEIGHT || height > IMG_MAX_HEIGHT {
		height = IMG_MIN_HEIGHT // 设置为默认最小值
	}
	return optionFunc(func(igo *ImageGenOption) {
		igo.Width = width
		igo.Height = height
	})
}

func WithFormat(format string) Option {
	if format != FORMAT_PNG && format != FORMAT_JPEG {
		panic("invalid image format")
	}
	return optionFunc(func(igo *ImageGenOption) {
		igo.Format = format
	})
}

func WithStyle(style ImageStyle) Option {
	return optionFunc(func(igo *ImageGenOption) {
		igo.Style = style
	})
}

type ImageGenOption struct {
	Format string
	Style  ImageStyle
	Width  int
	Height int
}

func RandomWidth() int {
	return rand.IntN(IMG_MAX_WIDTH-IMG_MIN_WIDTH+1) + IMG_MIN_WIDTH
}

func RandomHeight() int {
	return rand.IntN(IMG_MAX_HEIGHT-IMG_MIN_HEIGHT+1) + IMG_MIN_HEIGHT
}

func RandomFormat() string {
	idx := rand.IntN(len(FixedFormats))
	return FixedFormats[idx]
}

func RandomStyle() ImageStyle {
	idx := rand.IntN(len(FixedStyles))
	return FixedStyles[idx].Randomize()
}

func RandomRGBAColor() color.RGBA {
	return color.RGBA{
		R: uint8(rand.IntN(256)),
		G: uint8(rand.IntN(256)),
		B: uint8(rand.IntN(256)),
		A: uint8(rand.IntN(256)),
	}
}

func RandomFixedRGBAColor() color.RGBA {
	idx := rand.IntN(len(FixedColors))
	return FixedColors[idx]
}

type ImageInfo struct {
	Format string
	Width  int
	Height int
}

func GenImage(writer io.Writer, opts ...Option) (*ImageInfo, error) {
	genOption := ImageGenOption{
		Width:  RandomWidth(),
		Height: RandomHeight(),
		Format: RandomFormat(),
		Style:  RandomStyle(),
	}

	for _, opt := range opts {
		opt.apply(&genOption)
	}

	format := genOption.Format
	style := genOption.Style
	width := genOption.Width
	height := genOption.Height

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
		return nil, err
	}
	return &ImageInfo{
		Format: format,
		Width:  width,
		Height: height,
	}, nil
}
