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

	err = db.AutoMigrate(&models.User{}, &models.Session{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate tables %s", err)
	}

	err = db.SetupJoinTable(&models.Block{}, "Exercises", &models.ExerciseBlock{})
	if err != nil {
		return nil, fmt.Errorf("failed to set up join table between exercises and blocks tables %s", err)
	}
	err = db.AutoMigrate(&models.Exercise{}, &models.Block{}, &models.ExerciseBlock{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate tables %s", err)
	}

	err = db.SetupJoinTable(&models.Training{}, "Blocks", &models.TrainingBlock{})
	if err != nil {
		return nil, fmt.Errorf("failed to set up join table between exercises and blocks tables %s", err)
	}
	err = db.AutoMigrate(&models.Training{}, &models.TrainingBlock{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate tables %s", err)
	}

	//result = db.Exec("ALTER TABLE training_blocks ADD CONSTRAINT IF NOT EXISTS uniq_training_block UNIQUE (training_id, block_id)")
	//if result.Error != nil {
	//	return nil, fmt.Errorf("failed to create constraint: %s", result.Error)
	//}
	return db, err
}
