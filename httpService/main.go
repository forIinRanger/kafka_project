package main

import (
	"github.com/joho/godotenv"
	"httpservice/cmd"
	"log"
	"os"
)

func main() {
	dockerEnv := os.Getenv("DOCKER_ENV")
	log.Printf("DOCKER_ENV: '%s'", dockerEnv)

	if dockerEnv == "" {
		// Мы НЕ в Docker, загрузи .env файл
		if err := godotenv.Load("../.env"); err != nil {
			log.Println("No .env file found")
		}
	}
	cmd.RunServer()
}
