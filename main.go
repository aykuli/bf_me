package main

import (
	"bf_me/internal/configs"
	"bf_me/internal/storage"
	"fmt"
	"log"
)

func main() {
	config := configs.Parse()

	// ------- DATABASE -------
	db, err := storage.New(config.DatabaseURI)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(db)
}
