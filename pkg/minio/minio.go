package minio

import (
	"bf_me/internal/configs"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type IS3Storage interface {
	Ping() error
}

type S3Storage struct {
	config *configs.S3
}

func NewStorage(config *configs.S3) *S3Storage {
	return &S3Storage{config}
}

func (s *S3Storage) newConn() (*minio.Client, error) {
	return minio.New(
		s.config.URL,
		&minio.Options{Creds: credentials.NewStaticV4(s.config.AccessKey, s.config.SecretKey, "")})
}

// Ping For debug purpose
func (s *S3Storage) Ping() error {
	client, err := s.newConn()
	if err != nil {
		return fmt.Errorf("error connecting to minio: %s", err)
	}

	exists, err := client.BucketExists(context.Background(), s.config.Bucket)
	if err != nil || !exists {
		return fmt.Errorf("bucket doesnt exists: %s", err)
	}

	fmt.Println("Successfully connected to minio S3 storage")
	return nil
}
