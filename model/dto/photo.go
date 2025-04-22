package dto

import (
	"time"

	"github.com/follow1123/photos/model"
)

const (
	PHOTO_OPE_ID_REQUIRED uint = iota
	PHOTO_OPE_CREATE
)

type CreatePhotoFailedResult struct {
	UploadID uint   `json:"uploadId"`
	Message  string `json:"message"`
}

type PhotoDto struct {
	Operate   uint      `json:"-"`
	UploadID  uint      `json:"uploadId"`
	ID        uint      `json:"id" uri:"id" binding:"required_if=Operate 0,omitempty,min=1"`
	Desc      string    `json:"desc" form:"desc" binding:"required_if=Operate 1"`
	Format    string    `json:"format"`
	Uri       string    `json:"uri"`
	Size      int64     `json:"size"`
	Sum       string    `json:"-"`
	Width     int64     `json:"width"`
	Height    int64     `json:"height"`
	PhotoDate time.Time `json:"photoDate" time_format:"2006-01-02 15:04:05"`
}

func (p *PhotoDto) Update(photo *model.Photo) {
	p.ID = photo.ID
	p.Desc = photo.Desc
	p.Uri = photo.Uri
	p.Sum = photo.Sum
	p.Format = photo.Format
	p.Size = photo.Size
	p.Width = photo.Width
	p.Height = photo.Height
	p.PhotoDate = photo.PhotoDate
}

func (p *PhotoDto) ToModel() *model.Photo {
	photo := model.Photo{
		Desc:   p.Desc,
		Uri:    p.Uri,
		Sum:    p.Sum,
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
