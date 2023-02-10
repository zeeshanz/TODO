package controllers

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/zeeshanz/TODO/database"
	"github.com/zeeshanz/TODO/models"
	"github.com/zeeshanz/TODO/repos"
)

func CreateTodo(ctx *fiber.Ctx) error {
	fmt.Println("Adding new todo")
	c := context.Background()
	sessionId := ctx.Cookies("session-id")
	userUuid, err := database.GetFromRedis(c, sessionId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": 500,
			"message": err.Error,
		})
	}

	var todoDTO models.TodoDTO
	todoDTO.Uuid = uuid.Must(uuid.NewRandom()).String() // UUID will uniquely idenfiy the todo item
	todoDTO.UserUuid = userUuid

	if err = ctx.BodyParser(&todoDTO); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  404,
			"message": err.Error,
		})
	}

	aTodo, err := repos.CreateTodo(todoDTO.Uuid, todoDTO.UserUuid, todoDTO.Completed, todoDTO.TodoItem)
	if err != nil {
		return errors.New("failed to create new todo")
	}

	result := &models.TodoResponse{}
	if err := copier.Copy(&result, &aTodo); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Cannot map results",
		})
	}

	return ctx.JSON(fiber.Map{
		"status":   200,
		"todoItem": todoDTO.TodoItem,
		"uuid":     todoDTO.Uuid,
	})
}

/*
 * Retrieve todos from the database.
 */
func GetTodos(ctx *fiber.Ctx) error {
	sessionId := ctx.Cookies("session-id")
	c := context.Background()
	userUuid, err := database.GetFromRedis(c, sessionId)
	if err != nil {
		fmt.Println(err)
	}

	todoResponse := []models.TodoResponse{}
	todos, _ := repos.GetAllTodos(userUuid)
	if err := copier.Copy(&todoResponse, &todos); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Cannot map results",
		})
	}

	return ctx.Render("todos", fiber.Map{
		"Todos": todoResponse,
	})

}

func DeleteTodo(ctx *fiber.Ctx) error {
	var todoDTO models.TodoDTO
	if err := ctx.BodyParser(&todoDTO); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error,
		})
	}

	fmt.Printf("Deleting Todo uuid %v\n", todoDTO.Uuid)
	err := repos.DeleteTodo(todoDTO.Uuid)
	if err != nil {
		return ctx.SendStatus(404)
	}

	return ctx.SendStatus(200)
}

func CompleteTodo(ctx *fiber.Ctx) error {
	var todoDTO models.TodoDTO
	if err := ctx.BodyParser(&todoDTO); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error,
		})
	}

	fmt.Printf("Completing Todo uuid %v\n", todoDTO.Uuid)
	todoItem, err := repos.GetTodoItem(todoDTO.Uuid)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": err,
		})
	}

	var isCompleted = todoItem.Completed
	err = repos.UpdateTodoStatus(todoDTO.Uuid, !isCompleted)
	if err != nil {
		return ctx.SendStatus(404)
	}

	if isCompleted {
		return ctx.SendStatus(201) // meaning Todo is updated to not completed
	} else {
		return ctx.SendStatus(202) // meaning Todo is updated to completed
	}
}

func UpdateTodo(ctx *fiber.Ctx) error {
	var todoDTO models.TodoDTO
	if err := ctx.BodyParser(&todoDTO); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error,
		})
	}

	// Stop execution any further if length of string is less than 4 characters
	if len(todoDTO.TodoItem) < 4 {
		return ctx.SendStatus(403)
	}

	err := repos.UpdateTodoItem(todoDTO.Uuid, todoDTO.TodoItem)
	if err != nil {
		return ctx.SendStatus(404)
	}

	return ctx.SendStatus(200)
}
