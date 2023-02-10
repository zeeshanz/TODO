package routes

import (
	"github.com/gofiber/fiber/v2"
	apiControllers "github.com/zeeshanz/TODO/controllers/api"
)

func TodoRoute(route fiber.Router) {
	route.Post("/addNewTodo", apiControllers.AddNewTodo)
	route.Post("/deleteTodo", apiControllers.DeleteTodo)
	route.Post("/completeTodo", apiControllers.CompleteTodo)
	route.Post("/updateTodo", apiControllers.UpdateTodo)
	route.Get("/todos", apiControllers.GetTodos)
}
