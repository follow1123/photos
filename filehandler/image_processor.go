package filehandler

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
)

type ImageMetaInfo struct {
	Size       int64
	ImgFormat  string
	Resolution string
}

type ImageProcessor struct {
	Data           []byte
	size           int64
	hexSum         *string
	compressedData []byte
	metaInfo       *ImageMetaInfo
}

func NewImageProcessor(reader io.Reader) (*ImageProcessor, error) {
	imgData, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return &ImageProcessor{
		Data: imgData,
		size: int64(len(imgData)),
	}, nil
}

func (ip *ImageProcessor) GetHexSum() string {
	if ip.hexSum == nil {
		sum := md5.Sum(ip.Data)
		hexSum := hex.EncodeToString(sum[:])
		ip.hexSum = &hexSum
	}
	return *ip.hexSum
}

func (ip *ImageProcessor) GetMetaInfo() (ImageMetaInfo, error) {
	if ip.metaInfo == nil {
		// 获取图片其他信息
		imgConfig, format, err := image.DecodeConfig(bytes.NewReader(ip.Data))
		if err != nil {
			return *ip.metaInfo, err
		}
		ip.metaInfo = &ImageMetaInfo{
			Size:       ip.size,
			ImgFormat:  format,
			Resolution: fmt.Sprintf("%dx%d", imgConfig.Width, imgConfig.Height),
		}
	}
	return *ip.metaInfo, nil
}

func (ip *ImageProcessor) GetCompressedData() ([]byte, error) {
	if ip.size <= 204800 {
		return ip.Data, nil
	}

	metaInfo, err := ip.GetMetaInfo()
	if err != nil {
		return nil, err
	}

	if ip.compressedData == nil {
		var buf bytes.Buffer
		switch metaInfo.ImgFormat {
		case "jpeg":
			// 压缩并保存
			img, err := jpeg.Decode(bytes.NewReader(ip.Data))
			if err != nil {
				return nil, err
			}
			err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 30})
			if err != nil {
				return nil, err
			}
			ip.compressedData = buf.Bytes()
		case "png":
			return nil, errors.New("TODO: compressed png image")
		}
	}
	return ip.compressedData, nil
}
