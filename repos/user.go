package repos

import (
	"errors"

	"github.com/google/uuid"
	"github.com/zeeshanz/TODO/database"
	"github.com/zeeshanz/TODO/models"
	"golang.org/x/crypto/bcrypt"
)

/*
 * Add a new user. Check for duplication. Hash the password
 */
func AddUser(userInfo models.User) error {
	// Check if username already exists
	var tempUser models.User
	canAddThisUser := database.DB.Db.Where("username = ?", userInfo.Username).First(&tempUser).Error
	if canAddThisUser == nil {
		return errors.New("this username already exists")
	}

	// Add new user
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(userInfo.Password), bcrypt.DefaultCost)
	if err == nil {
		userInfo.Uuid = uuid.Must(uuid.NewRandom()).String() // UUID will uniquely idenfiy the user
		userInfo.Password = string(hashedPass)
		err := database.DB.Db.Create(&userInfo)
		if err.Error == nil {
			return nil
		}
		return err.Error
	}
	return err
}

/*
 * Authentication as part of user sign in.
 */
func AuthenticateUser(userInfo models.User) error {
	var tempUser models.User
	err := database.DB.Db.Where("username = ?", userInfo.Username).First(&tempUser).Error
	if err == nil {
		err = bcrypt.CompareHashAndPassword([]byte(tempUser.Password), []byte(userInfo.Password))
		if err == nil {
			return nil
		} else {
			return errors.New("incorrect password please try again")
		}
	} else {
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
