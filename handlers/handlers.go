package handlers

import (
	"context"
	"fmt"
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

var loggedInUser models.User

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

	loggedInUser = user
	fmt.Printf("loggedInUser: %v\n", loggedInUser)
	cookie := new(fiber.Cookie)
	cookie.Name = "user_uuid"
	cookie.Value = initializers.GetUserUuid(user.Username)
	cookie.Expires = time.Now().Add(1 * time.Hour)
	ctx.Cookie(cookie)
	c := context.Background()
	initializers.SetToRedis(c, "user_uuid", user.Uuid)

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

	// Auto sign out if cache expired
	if len(username) == 0 {
		return ctx.Render("index", fiber.Map{"signInStatus": "0"})
	}

	fmt.Printf("username retrieved from session is %v\n", username)
	userUuid := initializers.GetUserUuid(username)
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

func AddNewTodo(c *fiber.Ctx) error {
	task_new := new(models.Task)
	fmt.Printf("Adding new Todo: %v\n", task_new)
	if err := c.BodyParser(task_new); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	fmt.Printf("Logged in user is: %v\n", &loggedInUser)
	initializers.DB.Db.Model(&loggedInUser).Association("Tasks").Append(task_new)
	fmt.Println(loggedInUser.ID)

	var task_temp = []models.Task{}
	initializers.DB.Db.Model(&loggedInUser).Association("Tasks").Find(&task_temp)

	return c.JSON(fiber.Map{
		"success": true,
		"Task":    task_new,
	})
}

func SignOutUser(ctx *fiber.Ctx) error {
	fmt.Println("Signing out user")
	c := context.Background()
	initializers.DeleteFromRedis(c, "username")
	return ctx.Render("index", fiber.Map{"signInStatus": "0"})
}
