package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zeeshanz/TODO/initializers"
	"github.com/zeeshanz/TODO/models"
)

var current_data []models.User
var LogInType string = "-2" //-2 means nothing, -1 means account doesn't exist, 0 exists but means wrong password
var newReload bool = true
var regLoad bool = false
var current_user *models.User
var userExists bool = false

const SecretKey = "secret"

var userLoggedIn models.User

/*
 * Sign up a new user.
 */
func SignUpUser(c *fiber.Ctx) error {
	var creds models.User
	// First we need to parse the variable ctx to receive the credentials
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
			"message": "No account found. Please signup first.",
		})
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "username"
	cookie.Value = strconv.FormatUint(uint64(initializers.GetUserId(user.Username)), 10)
	cookie.Expires = time.Now().Add(24 * time.Hour)
	ctx.Cookie(cookie)
	c := context.Background()
	initializers.SetToRedis(c, "username", user.Username)

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
	username := initializers.GetFromRedis(c, "username")
	fmt.Printf("username retrieved from session is %v\n", username)
	userId := initializers.GetUserId(username)
	taskResponse, err := initializers.GetTasksForUser(userId)

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
