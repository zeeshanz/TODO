package handlers

import (
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

var userLoggedIn models.User

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

func SignInUser(c *fiber.Ctx) error {
	LogInType = "-1"

	user1 := new(models.User)
	if err := c.BodyParser(user1); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	for i := 0; i < len(current_data); i++ {
		if current_data[i].Username == user1.Username {

			LogInType = "0"
			if current_data[i].Password == user1.Password {

				LogInType = "-2"
				userLoggedIn = current_data[i]
				cookie := fiber.Cookie{
					Name:     "jwt",
					Value:    "token-to-be",
					Expires:  time.Now().Add(time.Hour * 24),
					HTTPOnly: true,
				}
				c.Cookie(&cookie)

				return c.Redirect("/taskPage")
			}

		}
	}

	newReload = false
	regLoad = false
	current_user = user1
	return c.Redirect("/")
}
