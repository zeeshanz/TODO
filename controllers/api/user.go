package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/zeeshanz/TODO/database"
	"github.com/zeeshanz/TODO/models"
	"github.com/zeeshanz/TODO/repos"
)

/*
 * Sign up a new user.
 */
func SignUpUser(c *fiber.Ctx) error {
	var userDTO models.UserDTO

	// Parse ctx to receive the credentials
	if err := c.BodyParser(&userDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  404,
			"message": err.Error,
		})
	}

	if err := repos.AddUser(userDTO.Username, userDTO.Password); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Could not add new user.",
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

	err = repos.AuthenticateUser(user.Username, user.Password)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Sign in was unsuccessful.",
		})
	}

	var userUuid = repos.GetUserUuid(user.Username)
	fmt.Printf("loggedInUser: %v\n", user)
	cookie := new(fiber.Cookie)
	sessionId := uuid.New().String()
	cookie.Name = "session-id"
	cookie.Value = sessionId
	cookie.Expires = time.Now().Add(1 * time.Hour)
	ctx.Cookie(cookie)
	c := context.Background()
	database.SetToRedis(c, sessionId, userUuid)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Sign-in successful.",
	})
}

func SignOutUser(ctx *fiber.Ctx) error {
	fmt.Println("Signing out user")
	c := context.Background()
	database.DeleteFromRedis(c, "username")
	return ctx.Render("index", fiber.Map{"signInStatus": "0"})
}
