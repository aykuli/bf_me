package models

import "gorm.io/gorm"

type Exercise struct {
	gorm.Model
	TitleEn  string `gorm:"unique;not null"`
	TitleRu  string `gorm:"unique;not null"`
	Filename string `gorm:"unique;not null"`
	Tags     []Tag  `gorm:"many2many:exercises_tags;"`
}
