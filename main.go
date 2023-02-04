package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/zeeshanz/TODO/handlers"
	"github.com/zeeshanz/TODO/initializers"
)

func main() {

	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	ctx := context.TODO()
	initializers.ConnectDB(&config)
	initializers.ConnectRedis(ctx)

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/static", "./static")
	app.Static("/", "./views", fiber.Static{
		Index: "index.html",
	})

	setupRoutes(app)

	initializers.SetToRedis(ctx, "name", "redis-test-123")

	app.Listen(":3001")
}

func setupRoutes(app *fiber.App) {
	app.Post("/signInUser", handlers.SignInUser)
	app.Post("/signUpUser", handlers.SignUpUser)
	app.Get("/signOutUser", handlers.SignOutUser)
	app.Post("/addNewTodo", handlers.AddNewTodo)
	app.Get("/todos", handlers.ShowTodos)
}
