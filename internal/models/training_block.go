package models

import "gorm.io/gorm"

type TrainingBlock struct {
	gorm.Model
	TrainingID uint `gorm:"primaryKey"`
	BlockID    uint `gorm:"primaryKey"`
	BlockOrder uint `gorm:"not_null;default:0;"`
}
