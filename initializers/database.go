package initializers

import (
	"errors"
	"fmt"
	"log"

	"github.com/zeeshanz/TODO/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

func ConnectDB(config *Config) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Canada/Eastern", config.DBHost, config.DBUserName, config.DBUserPassword, config.DBName, config.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the Database")
	}
	fmt.Println("? Connected Successfully to the Database")

	DB = Dbinstance{
		Db: db,
	}
}

/*
 * Add a new user. Check for duplication. Hash the password
 */
func AddUser(userInfo models.User) error {
	var tempUser models.User
	canAddThisUser := DB.Db.Where("username = ?", userInfo.Username).First(&tempUser).Error
	if canAddThisUser == nil {
		return errors.New("This username already exists.")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(userInfo.Password), bcrypt.DefaultCost)
	if err == nil {
		userInfo.Password = string(hashedPass)
		tempTasks := []models.Task{}
		userInfo.Tasks = tempTasks
		err := DB.Db.Create(&userInfo)
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
	err := DB.Db.Where("username = ?", userInfo.Username).First(&tempUser).Error
	if err == nil {
		err = bcrypt.CompareHashAndPassword([]byte(tempUser.Password), []byte(userInfo.Password))
		if err == nil {
			return nil
		} else {
			return errors.New("Incorrect password. Please try again.")
		}
	} else {
		return errors.New("Username not found.")
	}
}

/*
 * Retrieve username
 */
func GetUserId(userInfo models.User) uint {
	var tempUser models.User
	DB.Db.Where("username = ?", userInfo.Username).First(&tempUser)
	return tempUser.ID
}

func ReturnTasksWithID(ID uint) ([]models.Task, error) {
	tempTasks := []models.User{}
	// As the user model stores a task struct, and not TaskResponse, we need to create
	// another variable so we can return the TaskResponse
	resTasks := []models.Task{}
	err := DB.Db.Where("ID = ?", ID).First(&tempTasks).Error
	if err != nil {
		return resTasks, err
	}
	// If no error, we can copy the tasks into the resTasks. The copier function handles this for us
	// copier.Copy(&resTasks, &tempTasks)
	return resTasks, nil

}
