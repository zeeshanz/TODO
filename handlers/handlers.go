package handlers

import (
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
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	for i := 0; i < len(current_data); i++ {
		if current_data[i].Username == user.Username {
			userExists = true
		}
	}

	if !userExists {
		initializers.DB.Db.Create(&user)
	}
	current_user = user
	regLoad = true
	newReload = false
	// return c.Redirect("/")

	// Return status 200 OK.
	var response = c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Added user " + current_user.Username,
	})
	return response
	// return c.JSON(fiber.Map{
	// 	"error":      false,
	// 	"msg":        nil,
	// 	"username":   current_user.Username,
	// 	"statusCode": 200,
	// })
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
