package dto

import (
	"mime/multipart"
	"time"

	"github.com/follow1123/photos/model"
)

type CreatePhotoParam struct {
	UploadID   uint                  `json:"uploadId"`
	Desc       string                `json:"desc"`
	Uri        string                `json:"uri"`
	PhotoDate  time.Time             `json:"photoDate" time_format:"2006-01-02 15:04:05"`
	FileHeader *multipart.FileHeader `json:"-"`
}

func (cpp *CreatePhotoParam) ToModel() *model.Photo {
	photo := model.Photo{
		Desc: cpp.Desc,
		Uri:  cpp.Uri,
	}
	if cpp.PhotoDate.IsZero() {
		photo.PhotoDate = time.Now()
	} else {
		photo.PhotoDate = cpp.PhotoDate
	}
	return &photo
}

type CreatePhotoFailedResult struct {
	UploadID uint   `json:"uploadId"`
	Message  string `json:"message"`
}

type PhotoPageParam struct {
	Desc string `json:"desc"`
}

func (ppp *PhotoPageParam) ToModel() *model.Photo {
	return &model.Photo{Desc: ppp.Desc}
}

type PhotoParam struct {
	ID        uint      `json:"id" uri:"id" binding:"required"`
	Desc      string    `json:"desc"`
	PhotoDate time.Time `json:"photoDate" time_format:"2006-01-02 15:04:05"`
}

type PhotoDto struct {
	ID        uint      `json:"id"`
	Desc      string    `json:"desc"`
	Format    string    `json:"format"`
	Size      int64     `json:"size"`
	Width     int64     `json:"width"`
	Height    int64     `json:"height"`
	PhotoDate time.Time `json:"photoDate" time_format:"2006-01-02 15:04:05"`
}

func (p *PhotoDto) Update(photo *model.Photo) {
	p.ID = photo.ID
	p.Desc = photo.Desc
	p.Format = photo.Format
	p.Size = photo.Size
	p.Width = photo.Width
	p.Height = photo.Height
	p.PhotoDate = photo.PhotoDate
}

func (p *PhotoDto) ToModel() *model.Photo {
	photo := model.Photo{
		Desc:   p.Desc,
		Size:   p.Size,
		Format: p.Format,
		Width:  p.Width,
		Height: p.Height,
	}
	if p.PhotoDate.IsZero() {
		photo.PhotoDate = time.Now()
	} else {
		photo.PhotoDate = p.PhotoDate
	}
	return &photo
}
