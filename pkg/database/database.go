package database

import (
	"bf_me/internal/models"
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database interface {
	Ping(ctx context.Context) error
	Close() error
}

func New(uri string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database %s", err)
	}
	fmt.Println("Successfully connected to database")
	// Enable UUID generation function for PostgreSQL
	result := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if result.Error != nil {
		return nil, fmt.Errorf("failed to enable extension for uuid: %s", err)
	}
	err = db.AutoMigrate(&models.Exercise{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate exercises table %s", err)
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate user table %s", err)
	}

	err = db.AutoMigrate(&models.Session{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate session table %s", err)
	}

	err = db.AutoMigrate(&models.Tag{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate tag table %s", err)
	}

	return db, err
}
