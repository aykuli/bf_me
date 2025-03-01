package configs

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Configs struct {
	DatabaseURI string
	Address     string
}

func Parse() *Configs {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Configs{
		DatabaseURI: os.Getenv("DATABASE_URL"),
		Address:     fmt.Sprintf(":%s", os.Getenv("PORT")),
	}
}
