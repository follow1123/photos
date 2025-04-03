package model

import (
	"database/sql"

	"gorm.io/gorm"
)

type Photo struct {
	gorm.Model
	Desc       string
	Type       sql.NullString
	Path       string
	Size       sql.NullInt64
	Resolution sql.NullString
	PhotoDate  sql.NullTime
}
