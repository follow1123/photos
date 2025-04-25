package imagemanager

import (
	"bytes"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"slices"

	"github.com/follow1123/photos/logger"
)

var SupportedFormats = [...]string{"jpeg", "png"}

type ImageInfo struct {
	Size   int64
	Format string
	Width  int64
	Height int64
}

type ImageProcessor struct {
	logger    *logger.AppLogger
	data      []byte
	reader    io.ReadCloser
	imageInfo *ImageInfo
}

func NewImageProcessor(rc io.ReadCloser, logger *logger.AppLogger) *ImageProcessor {
	return &ImageProcessor{reader: rc, logger: logger}
}

func (_ *ImageProcessor) checkFormat(format string) error {
	if !slices.Contains(SupportedFormats[:], format) {
		return ErrUnsupportedImageFormat
	}
	return nil
}

func (ip *ImageProcessor) GetData() ([]byte, error) {
	if ip.data == nil {
		data, err := io.ReadAll(ip.reader)
		if err != nil {
			ip.logger.Error("read image error: %v", err)
			return nil, err
		}

		if err = ip.reader.Close(); err != nil {
			ip.logger.Error("close image error: %v", err)
			return nil, err
		}

		// 获取图片其他信息
		imgConfig, format, err := image.DecodeConfig(bytes.NewReader(data))
		if err != nil {
			ip.logger.Error("decode image config error: %v", err)
			return nil, err
		}
		if err := ip.checkFormat(format); err != nil {
			ip.logger.Error("check image format error: %v", err)
			return nil, err
		}

		ip.imageInfo = &ImageInfo{
			Size:   int64(len(data)),
			Format: format,
			Width:  int64(imgConfig.Width),
			Height: int64(imgConfig.Height),
		}
		ip.data = data
	}
	return ip.data, nil
}

func (ip *ImageProcessor) GetImageInfo() (*ImageInfo, error) {
	if ip.imageInfo == nil {
		_, err := ip.GetData()
		if err != nil {
			return nil, err
		}
	}
	return ip.imageInfo, nil
}

func (ip *ImageProcessor) GetCompressedData() ([]byte, error) {
	imgInfo, err := ip.GetImageInfo()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	switch imgInfo.Format {
	case "jpeg":
		// 压缩并保存
		img, err := jpeg.Decode(bytes.NewReader(ip.data))
		if err != nil {
			ip.logger.Error("decode image error: %v", err)
			return nil, err
		}
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 30})
		if err != nil {
			ip.logger.Error("compresse image error: %v", err)
			return nil, err
		}
		return buf.Bytes(), nil
	case "png":
		ip.logger.Warn("compressed png format image to be implemented")
		// TODO: compressed png image
		return ip.data, nil
	}
	return nil, ErrUnsupportedImageFormat
}
