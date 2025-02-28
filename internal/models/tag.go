package models

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	TitleEn   string     `gorm:"unique;not null"`
	TitleRu   string     `gorm:"unique;not null"`
	Exercises []Exercise `gorm:"many2many:exercises_tags;"`
}
