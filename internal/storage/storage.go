package storage

import (
	"bf_me/pkg/minio"
	"gorm.io/gorm"
)

type Storage struct {
	DB *gorm.DB
	S3 *minio.IS3Storage
}
