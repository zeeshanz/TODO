package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/zeeshanz/TODO/database"
	apiRoutes "github.com/zeeshanz/TODO/routes/api"
)

func main() {

	config, err := database.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	ctx := context.TODO()
	database.ConnectDB(&config)
	database.ConnectRedis(ctx)

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/static", "./static")
	app.Static("/", "./views", fiber.Static{
		Index: "index.html",
	})

	setupRoutes(app)

	database.SetToRedis(ctx, "name", "redis-test-123")

	app.Listen(":3001")
}

func setupRoutes(app *fiber.App) {
	api := app.Group("")
	apiRoutes.TodoRoute(api.Group(""))
	apiRoutes.UserRoute(api.Group(""))
}
