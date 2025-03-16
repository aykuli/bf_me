package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Exercise struct {
	gorm.Model
	TitleEn  string         `gorm:"not null"`
	TitleRu  string         `gorm:"not null"`
	Filename string         `gorm:"unique;not null"`
	Tips     pq.StringArray `gorm:"type:text[]"`
}
