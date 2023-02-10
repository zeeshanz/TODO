package main

import (
	"fmt"
	"log"

	"github.com/zeeshanz/TODO/database"
	"github.com/zeeshanz/TODO/models"
)

func init() {
	config, err := database.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	database.ConnectDB(&config)
}

func main() {
	database.DB.Db.AutoMigrate(&models.Todo{})
	database.DB.Db.AutoMigrate(&models.User{})
	fmt.Println("? Migration complete")
}
