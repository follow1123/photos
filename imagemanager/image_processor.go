package imagemanager

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"image"
	"image/jpeg"
	_ "image/png"
)

type ImageInfo struct {
	Name   string
	Size   int64
	Format string
	Width  int64
	Height int64
}

type ImageProcessor struct {
	data      []byte
	imageInfo *ImageInfo
}

func (ip *ImageProcessor) GetData() []byte {
	return ip.data
}

func (ip *ImageProcessor) GetHexSum() string {
	sum := md5.Sum(ip.data)
	return hex.EncodeToString(sum[:])
}

func (ip *ImageProcessor) GetImageInfo() (*ImageInfo, error) {
	if ip.imageInfo == nil {
		// 获取图片其他信息
		imgConfig, format, err := image.DecodeConfig(bytes.NewReader(ip.data))
		if err != nil {
			return nil, err
		}
		ip.imageInfo = &ImageInfo{
			Size:   int64(len(ip.data)),
			Format: format,
			Width:  int64(imgConfig.Width),
			Height: int64(imgConfig.Height),
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
			return nil, err
		}
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 30})
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	case "png":
		// TODO: compressed png image
		return ip.data, nil
	}
	return nil, errors.New("error image format")
}
