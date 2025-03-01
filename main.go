package main

import (
	"bf_me/internal/configs"
	"bf_me/internal/routes"
	"bf_me/internal/storage"
	"log"
	"net/http"
)

func main() {
	config := configs.Parse()

	// ------- DATABASE -------
	db, err := storage.New(config.DatabaseURI)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	routes.RegisterExercisesRoutes(mux, db)
	//routes.RegisterTagsRoutes(mux)
	log.Fatal(http.ListenAndServe("localhost:3000", mux))
}
