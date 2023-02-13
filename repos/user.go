package repos

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/zeeshanz/TODO/database"
	"github.com/zeeshanz/TODO/models"
	"golang.org/x/crypto/bcrypt"
)

/*
 * Add a new user. Check for duplication. Hash the password
 */
func AddUser(username string, password string) error {
	// Check if username already exists
	var user models.User
	canAddThisUser := database.DB.Db.Where("username = ?", username).First(&user).Error
	if canAddThisUser == nil {
		return errors.New("this username already exists")
	}

	// Add new user
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err == nil {
		user.Uuid = uuid.Must(uuid.NewRandom()).String() // UUID will uniquely idenfiy the user
		user.Username = strings.TrimSpace(username)
		user.Password = string(hashedPass)
		err := database.DB.Db.Create(&user)
		if err.Error == nil {
			return nil
		}
		return err.Error
	}
	return err
}

/*
- Authentication as part of user sign in.
*/
func AuthenticateUser(username string, password string) error {
	var user models.User
	err := database.DB.Db.Where("username = ?", username).First(&user).Error
	if err == nil {
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err == nil {
			return nil
		} else {
			fmt.Println("incorrect password please try again")
			return errors.New("incorrect password please try again")
		}
	} else {
		fmt.Println("username not found")
		return errors.New("username not found")
	}
}

/*
 * Retrieve username
 */
func GetUserUuid(username string) string {
	var tempUser models.User
	err := database.DB.Db.Where("username = ?", username).First(&tempUser).Error
	if err != nil {
		return ""
	}
	return tempUser.Uuid
}
