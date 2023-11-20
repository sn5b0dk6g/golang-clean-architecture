package main

import (
	"go-rest-api/infrastructure"
	logf "go-rest-api/infrastructure/log"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("GO_ENV") == "dev" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalln(err)
		}
	}

	var app = infrastructure.NewConfig().
		Name(os.Getenv("APP_NAME")).
		Logger(logf.InstanceSlogLogger).
		Validator().
		DbSQL().
		DbNoSQL()

	app.WebServerPort(os.Getenv("APP_PORT")).
		WebServer().
		Start()
}
