package minio

import (
	"bf_me/internal/configs"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type IS3Storage interface {
	Ping() error
	Upload(dst string, src io.Reader) (string, error)
	//RemoveObject(path string) error
	//GetContent(src string) ([]byte, error)
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

func (s *S3Storage) Upload(dst string, src io.Reader, contentType string) (string, error) {
	client, err := s.newConn()
	if err != nil {
		return "", err
	}
	info, err := client.PutObject(context.Background(), s.config.Bucket, dst, src, -1, minio.PutObjectOptions{
		ContentType:          contentType,
		DisableContentSha256: true,
	},
	)
	if err != nil {
		return "", err
	}
	fmt.Printf("File %s was saved", s.config.Bucket+"/"+info.Key)
	return s.config.Bucket + "/" + info.Key, nil
}

func (s *S3Storage) Delete(fname string) error {
	client, err := s.newConn()
	if err != nil {
		return err
	}
	err = client.RemoveObject(context.Background(),
		s.config.Bucket,
		strings.TrimPrefix(fname, s.config.Bucket+"/"), minio.RemoveObjectOptions{
			GovernanceBypass: false, ForceDelete: true,
		},
	)
	if err == nil {
		fmt.Printf("File %s was succesfully removed", fname)
	}
	return err
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
