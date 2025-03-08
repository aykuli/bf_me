package models

import "github.com/jackc/pgx/v5/pgtype"

type Session struct {
	ID   pgtype.UUID `gorm:"primaryKey"`
	User User        `gorm:"not null"`
}
