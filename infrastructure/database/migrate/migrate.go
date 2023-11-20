package main

import (
	"fmt"
	"go-rest-api/domain"
	"go-rest-api/infrastructure/database"
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
	dbConn, err := database.OpenPostgres()
	if err != nil {
		log.Fatalln(err)
	}
	defer database.ClosePostgres(dbConn)
	if err := dbConn.AutoMigrate(&domain.User{}, &domain.Task{}); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Successfully Migrated")
	}
}
