package main

import (
	"bf_me/internal/configs"
	"bf_me/internal/routes"
	"bf_me/internal/storage"
	"bf_me/pkg/database"
	"bf_me/pkg/minio"
	"log"
	"net/http"
)

func main() {
	// ------- CONFIGs from .env file -------
	config := configs.Parse()
	var err error

	// ------- S3 MINIO -------
	s3 := minio.NewStorage(&config.S3)
	err = s3.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// ------- DATABASE -------
	db, err := database.New(config.DatabaseURI)
	if err != nil {
		log.Fatal(err)
	}

	st := &storage.Storage{DB: db, S3: s3}
	mux := http.NewServeMux()

	// ------- ROUTES -------
	routes.RegisterExercisesRoutes(mux, st)
	routes.RegisterTagsRoutes(mux, st.DB)

	// ------- RSERVER -------

	log.Fatal(http.ListenAndServe(config.Address, mux))
}
