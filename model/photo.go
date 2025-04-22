package model

import (
	"time"

	"gorm.io/gorm"
)

type Photo struct {
	gorm.Model
	Desc      string
	Format    string
	Uri       string
	Size      int64
	Sum       string
	Width     int64
	Height    int64
	PhotoDate time.Time
}
