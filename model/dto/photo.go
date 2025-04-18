package dto

import (
	"database/sql"
	"mime/multipart"
	"time"

	"github.com/follow1123/photos/model"
)

const (
	PHOTO_OPE_ID_REQUIRED uint = iota
	PHOTO_OPE_CREATE
)

type PhotoDto struct {
	Operate       uint                  `json:"-"`
	UploadID      uint                  `json:"uploadId"`
	ID            uint                  `json:"id" uri:"id" binding:"required_if=Operate 0,omitempty,min=1"`
	MultipartFile *multipart.FileHeader `json:"-"`
	Desc          string                `json:"desc" form:"desc" binding:"required_if=Operate 1"`
	Type          string                `json:"type"`
	Uri           string                `json:"uri" form:"uri"`
	Size          int64                 `json:"size"`
	Sum           string                `json:"-"`
	Resolution    string                `json:"resolution"`
	PhotoDate     time.Time             `json:"photoDate" time_format:"2006-01-02 15:04:05"`
}

func (p *PhotoDto) Update(photo *model.Photo) {
	p.ID = photo.ID
	p.Desc = photo.Desc
	p.Uri = photo.Uri
	p.Sum = photo.Sum
	if photo.Type.Valid {
		p.Type = photo.Type.String
	}
	if photo.Size.Valid {
		p.Size = photo.Size.Int64
	}
	if photo.Resolution.Valid {
		p.Resolution = photo.Resolution.String
	}
	if photo.PhotoDate.Valid {
		p.PhotoDate = photo.PhotoDate.Time
	}
}

func (p *PhotoDto) UpdateToModel(photo *model.Photo) {
	photo.ID = p.ID
	if p.Type != "" {
		photo.Type.String = p.Type
		photo.Type.Valid = true
	} else {
		photo.Type.Valid = false
	}
	if p.Size != 0 {
		photo.Size.Int64 = p.Size
		photo.Size.Valid = true
	} else {
		photo.Size.Valid = false
	}
	if p.Resolution != "" {
		photo.Resolution.String = p.Resolution
		photo.Resolution.Valid = true
	} else {
		photo.Resolution.Valid = false
	}
	if !p.PhotoDate.IsZero() {
		photo.PhotoDate.Time = p.PhotoDate
		photo.PhotoDate.Valid = true
	} else {
		photo.PhotoDate.Valid = false
	}
}

func (p *PhotoDto) ToModel() *model.Photo {
	photo := model.Photo{Desc: p.Desc, Uri: p.Uri, Sum: p.Sum}
	photo.Type = sql.NullString{}
	if p.Type != "" {
		photo.Type.String = p.Type
		photo.Type.Valid = true
	}
	photo.Size = sql.NullInt64{}
	if p.Size != 0 {
		photo.Size.Int64 = p.Size
		photo.Size.Valid = true
	}
	if p.Resolution != "" {
		photo.Resolution.String = p.Resolution
		photo.Resolution.Valid = true
	}
	if p.PhotoDate.IsZero() {
		photo.PhotoDate.Time = time.Now()
		photo.PhotoDate.Valid = true
	} else {
		photo.PhotoDate.Time = p.PhotoDate
		photo.PhotoDate.Valid = true
	}
	return &photo
}
