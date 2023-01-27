package main

import (
	"fmt"
	"log"

	"github.com/zeeshanz/TODO/initializers"
	"github.com/zeeshanz/TODO/models"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)
}

func main() {
	initializers.DB.AutoMigrate(&models.Task{})
	initializers.DB.AutoMigrate(&models.TaskResponse{})
	initializers.DB.AutoMigrate(&models.User{})
	initializers.DB.AutoMigrate(&models.UserResponse{})
	fmt.Println("? Migration complete")
}