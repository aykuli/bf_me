package main

import (
	"bf_me/internal/configs"
	"bf_me/internal/routes"
	"bf_me/internal/storage"
	"bf_me/pkg/minio"
	"log"
	"net/http"
)

func main() {
	config := configs.Parse()
	var err error

	// S3 Storage
	var s3Storage minio.IS3Storage
	s3Storage = minio.NewStorage(&config.S3)
	err = s3Storage.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// ------- DATABASE -------
	db, err := storage.New(config.DatabaseURI)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	routes.RegisterExercisesRoutes(mux, db)
	routes.RegisterTagsRoutes(mux, db)

	log.Fatal(http.ListenAndServe(config.Address, mux))
}
