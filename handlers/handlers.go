package initializers

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/zeeshanz/TODO/models"
	"gorm.io/gorm"
)

var current_data []models.User
var LogInType string = "-2" //-2 means nothing, -1 means account doesn't exist, 0 exists but means wrong password
var newReload bool = true
var regLoad bool = false
var current_user *models.User
var userExists bool = false

const SecretKey = "secret"

var userLoggedIn models.User

func ListUsers(c *fiber.Ctx) error {

	users := []models.User{}
	database.DB.Db.Find(&users)
	current_data = users

	if newReload == true {
		LogInType = "-2"
		return c.Render("mainPage", fiber.Map{"LogInType": LogInType})
	} else {
		newReload = true
		var userExistTemp bool = userExists
		var regLoadTemp bool = regLoad
		userExists = false
		return c.Render("mainPage", fiber.Map{
			"LogInType":     LogInType,
			"usernameInput": current_user.Username,
			"passwordInput": current_user.Password,
			"userExists":    strconv.FormatBool(userExistTemp),
			"regLoad":       strconv.FormatBool(regLoadTemp),
		})

	}

}

func CreateUser(c *fiber.Ctx) error {
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
		database.DB.Db.Create(&user)
	}
	current_user = user
	regLoad = true
	newReload = false
	return c.Redirect("/")
}

func LogInLogic(c *fiber.Ctx) error {
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

				claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
					Issuer:    user1.Username,
					ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //1 day
				})

				token, err := claims.SignedString([]byte(SecretKey))

				if err != nil {

					return c.JSON(fiber.Map{
						"message": "could not login",
					})
				}

				cookie := fiber.Cookie{
					Name:     "jwt",
					Value:    token,
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

func HandleTaskPage(c *fiber.Ctx) error {

	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.SendString("Unauthenticated, please Sign In!")
	}

	claims := token.Claims.(*jwt.StandardClaims)
	if claims.Issuer != userLoggedIn.Username {
		return c.SendString("Unauthenticated, please Sign In!")
	}

	var UserTasks []models.Task
	database.DB.Db.Model(&userLoggedIn).Association("Tasks").Find(&UserTasks)

	sort.Slice(UserTasks, func(i, j int) bool {
		return UserTasks[i].ID < UserTasks[j].ID
	})

	return c.Render("taskPage", UserTasks)
}

func SignOut(c *fiber.Ctx) error {

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.Redirect("/")
}

func AddTask(c *fiber.Ctx) error {

	task_new := new(models.Task)
	if err := c.BodyParser(task_new); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	database.DB.Db.Model(&userLoggedIn).Association("Tasks").Append(task_new)
	fmt.Println(userLoggedIn.ID)

	var task_temp = []models.Task{}
	database.DB.Db.Model(&userLoggedIn).Association("Tasks").Find(&task_temp)

	return c.JSON(fiber.Map{
		"success": true,
		"Task":    task_new,
	})
}

func DeleteTask(c *fiber.Ctx) error {

	task_del := new(models.Task)
	if err := c.BodyParser(task_del); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var task_temp models.Task

	//change below done before
	error := database.DB.Db.Delete(&task_temp, task_del.ID).Error
	if errors.Is(error, gorm.ErrRecordNotFound) {
		return errors.New("task not found")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"deleted": true,
	})

}

func UpdateTask(c *fiber.Ctx) error {
	task_up := new(models.Task)
	if err := c.BodyParser(task_up); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	database.DB.Db.Model(&task_up).Select("Finished", "Detail").Where("id = ?", task_up.ID).Updates(models.Task{Finished: !task_up.Finished, Detail: task_up.Detail})

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"updated": true,
	})

}
