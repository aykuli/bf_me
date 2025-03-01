package storage

import (
	"bf_me/internal/models"
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage interface {
	Ping(ctx context.Context) error
	Close() error
}

func New(uri string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database %s", err)

	}

	err = db.AutoMigrate(&models.Exercise{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate exercises table %s", err)
	}

	err = db.AutoMigrate(&models.Tag{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate tag table %s", err)
	}

	return db, err
}
