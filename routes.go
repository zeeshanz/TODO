package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zeeshanz/TODO/handlers"
)

func setupRoutes(app *fiber.App) {

	app.Post("/loginUser", handlers.SignInUser)
	app.Post("/user", handlers.SignUpUser)
}
