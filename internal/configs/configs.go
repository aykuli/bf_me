package configs

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type S3 struct {
	AccessKey string
	SecretKey string
	URL       string
	Bucket    string
}
type Configs struct {
	DatabaseURI string
	Address     string
	S3
}

func Parse() *Configs {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Configs{
		DatabaseURI: os.Getenv("DATABASE_URL"),
		Address:     fmt.Sprintf(":%s", os.Getenv("PORT")),
		S3: S3{
			AccessKey: os.Getenv("MINIO_ACCESS_KEY"),
			SecretKey: os.Getenv("MINIO_SECRET_KEY"),
			URL:       os.Getenv("MINIO_URL"),
			Bucket:    os.Getenv("MINIO_BUCKET"),
		},
	}
}
