package models

import "github.com/jackc/pgx/v5/pgtype"

type Session struct {
	ID     pgtype.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID uint        `gorm:"not null;index"`
	User   User        `gorm:"foreignKey:UserID;not null"`
}
