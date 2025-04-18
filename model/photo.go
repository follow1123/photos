package model

import (
	"database/sql"

	"gorm.io/gorm"
)

type Photo struct {
	gorm.Model
	Desc       string
	Type       sql.NullString
	Uri        string
	Size       sql.NullInt64
	Sum        string
	Resolution sql.NullString
	PhotoDate  sql.NullTime
}
