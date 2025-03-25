package models

import "gorm.io/gorm"

type Training struct {
	gorm.Model
	TitleEn        string          `gorm:"not null"`
	TitleRu        string          `gorm:"not null"`
	Draft          bool            `gorm:"default:true"`
	Blocks         []Block         `gorm:"many2many:trainings_blocks;"`
	TrainingBlocks []TrainingBlock `gorm:"foreignKey:TrainingID;references:ID"`
}
