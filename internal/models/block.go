package models

import "gorm.io/gorm"

type Block struct {
	gorm.Model
	TitleEn        string          `gorm:"unique;not null"`
	TitleRu        string          `gorm:"unique;not null"`
	TotalDuration  uint8           // minutes
	OnTime         uint8           // seconds
	RelaxTime      uint8           // seconds
	Draft          bool            `gorm:"default:true"`
	Exercises      []Exercise      `gorm:"many2many:exercises_blocks;"`
	ExerciseBlocks []ExerciseBlock `gorm:"foreignKey:BlockID;references:ID"`
}
