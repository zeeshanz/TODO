package handlers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/zeeshanz/TODO/initializers"
	"github.com/zeeshanz/TODO/models"
)

/*
 * Sign up a new user.
 */
func SignUpUser(c *fiber.Ctx) error {
	var creds models.User
	// Parse ctx to receive the credentials
	err := c.BodyParser(&creds)
	if err != nil {
		fmt.Println("Error with parsing credentials")
	}

	err = initializers.AddUser(creds)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "User already exists.",
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "User successfully added.",
		})
	}
}

/*
 * Sign in the user. Create a cookie to remember the user.
 */
func SignInUser(ctx *fiber.Ctx) error {
	var user models.User
	err := ctx.BodyParser(&user)
	if err != nil {
		fmt.Println("Error with parsing credentials")
	}

	err = initializers.AuthenticateUser(user)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Sign in was unsuccessful.",
		})
	}

	var userUuid = initializers.GetUserUuid(user.Username)
	fmt.Printf("loggedInUser: %v\n", user)
	cookie := new(fiber.Cookie)
	sessionId := uuid.New().String()
	cookie.Name = "session-id"
	cookie.Value = sessionId
	cookie.Expires = time.Now().Add(1 * time.Hour)
	ctx.Cookie(cookie)
	c := context.Background()
	fmt.Printf("session-id saved in cookie: %v\n", sessionId)
	initializers.SetToRedis(c, sessionId, userUuid)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Sign-in successful.",
	})
}

/*
 * Retrieve todos from the database.
 */
func ShowTodos(ctx *fiber.Ctx) error {
	sessionId := ctx.Cookies("session-id")
	fmt.Printf("session-id retrieved from cookie: %v\n", sessionId)
	c := context.Background()
	userUuid, err := initializers.GetFromRedis(c, sessionId)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("user_uuid retrieved from Redis is %v\n", userUuid)
	todoResponse, err := initializers.GetTodosForUser(userUuid)

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

func AddNewTodo(ctx *fiber.Ctx) error {
	c := context.Background()
	userUuid, err := initializers.GetFromRedis(c, "user_uuid")
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

	result := initializers.DB.Db.Model(models.Todo{}).Create(&todo)
	if result.Error != nil {
		return errors.New("failed to create new todo")
	}

	return ctx.JSON(fiber.Map{
		"status":   200,
		"todoItem": todo.TodoItem,
		"uuid":     todo.Uuid,
	})
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
	todoId := initializers.GetTodoId(todo.Uuid)
	result := initializers.DB.Db.Delete(&todo, todoId)

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
	todoItem, err := initializers.GetTodoItem(todo.Uuid)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": err,
		})
	}

	var isCompleted = todoItem.Completed
	err = initializers.UpdateTodo(todoItem.Uuid, !isCompleted)
	if err != nil {
		return ctx.SendStatus(404)
	}

	if isCompleted {
		return ctx.SendStatus(201) // meaning Todo is updated to not completed
	} else {
		return ctx.SendStatus(202) // meaning Todo is updated to completed
	}
}

func SignOutUser(ctx *fiber.Ctx) error {
	fmt.Println("Signing out user")
	c := context.Background()
	initializers.DeleteFromRedis(c, "username")
	return ctx.Render("index", fiber.Map{"signInStatus": "0"})
}
