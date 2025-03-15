package models

import (
	"gorm.io/gorm"
)

type ExerciseBlock struct {
	gorm.Model
	ExerciseID    uint `gorm:"primaryKey"`
	BlockID       uint `gorm:"primaryKey"`
	ExerciseOrder uint `gorm:"not_null;default:0;"`
}
