package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load() //загрузка .env файла
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
