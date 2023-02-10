package controllers

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/zeeshanz/TODO/database"
	"github.com/zeeshanz/TODO/models"
	"github.com/zeeshanz/TODO/repos"
)

func AddNewTodo(ctx *fiber.Ctx) error {
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

	var todo models.Todo
	todo.Uuid = uuid.Must(uuid.NewRandom()).String() // UUID will uniquely idenfiy the todo item
	todo.UserUuid = userUuid

	if err = ctx.BodyParser(&todo); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  404,
			"message": err.Error,
		})
	}

	result := database.DB.Db.Model(models.Todo{}).Create(&todo)
	if result.Error != nil {
		return errors.New("failed to create new todo")
	}

	return ctx.JSON(fiber.Map{
		"status":   200,
		"todoItem": todo.TodoItem,
		"uuid":     todo.Uuid,
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

	todoResponse, err := repos.GetTodosForUser(userUuid)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": err,
		})
	} else {
		return ctx.Render("tasks", fiber.Map{
			"Todos": todoResponse,
		})
	}
}

func DeleteTodo(ctx *fiber.Ctx) error {
	var todo models.Todo
	if err := ctx.BodyParser(&todo); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error,
		})
	}

	fmt.Printf("Deleting Todo uuid %v\n", todo.Uuid)
	todoId := repos.GetTodoId(todo.Uuid)
	result := database.DB.Db.Delete(&todo, todoId)

	if result.RowsAffected == 0 {
		return ctx.SendStatus(404)
	}

	return ctx.SendStatus(200)
}

func CompleteTodo(ctx *fiber.Ctx) error {
	var todo models.Todo
	if err := ctx.BodyParser(&todo); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error,
		})
	}

	fmt.Printf("Completing Todo uuid %v\n", todo.Uuid)
	todoItem, err := repos.GetTodoItem(todo.Uuid)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": err,
		})
	}

	var isCompleted = todoItem.Completed
	err = repos.UpdateTodoStatus(todoItem.Uuid, !isCompleted)
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
	var todo models.Todo
	if err := ctx.BodyParser(&todo); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error,
		})
	}

	// Stop execution any further if length of string is less than 4 characters
	if len(todo.TodoItem) < 4 {
		return ctx.SendStatus(403)
	}

	todoItem, err := repos.GetTodoItem(todo.Uuid)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": err,
		})
	}

	err = repos.UpdateTodoItem(todoItem.Uuid, todo.TodoItem)
	if err != nil {
		return ctx.SendStatus(404)
	}

	return ctx.SendStatus(200)
}
