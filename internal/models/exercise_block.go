package models

import "time"

type ExerciseBlock struct {
	ExerciseID uint `gorm:"primaryKey"`
	BlockID    uint `gorm:"primaryKey"`
	Order      uint `gorm:"not_null;default:0"`
	CreatedAt  time.Time
}
