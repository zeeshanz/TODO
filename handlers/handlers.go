package handlers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
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

	var userId = initializers.GetUserUuid(user.Username)
	fmt.Printf("loggedInUser: %v\n", user)
	cookie := new(fiber.Cookie)
	cookie.Name = "user_uuid"
	cookie.Value = userId
	cookie.Expires = time.Now().Add(1 * time.Hour)
	ctx.Cookie(cookie)
	c := context.Background()
	initializers.SetToRedis(c, "user_uuid", userId)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Sign-in successful.",
	})
}

/*
 * Retrieve tasks from the database.
 */
func ShowTasks(ctx *fiber.Ctx) error {
	c := context.Background()
	userUuid, err := initializers.GetFromRedis(c, "user_uuid")

	// Auto sign out if cache expired
	if len(userUuid) == 0 {
		return ctx.Render("index", fiber.Map{"signInStatus": "0"})
	}

	fmt.Printf("user_uuid retrieved from Redis is %v\n", userUuid)
	taskResponse, err := initializers.GetTasksForUser(userUuid)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": err,
		})
	} else {
		return ctx.Render("tasks", fiber.Map{
			"Tasks": taskResponse,
		})
	}
}

func AddNewTodo(ctx *fiber.Ctx) error {
	c := context.Background()
	userUuid, err := initializers.GetFromRedis(c, "user_uuid")
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error,
		})
	}

	var todo models.Task
	todo.UserUuid = userUuid

	if err = ctx.BodyParser(&todo); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error,
		})
	}

	result := initializers.DB.Db.Model(models.Task{}).Create(&todo)
	if result.Error != nil {
		return errors.New("failed to create new task")
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Task created successfully.",
	})
}

func SignOutUser(ctx *fiber.Ctx) error {
	fmt.Println("Signing out user")
	c := context.Background()
	initializers.DeleteFromRedis(c, "username")
	return ctx.Render("index", fiber.Map{"signInStatus": "0"})
}
