package initializers

import (
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
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
	// Check if username already exists
	var tempUser models.User
	canAddThisUser := DB.Db.Where("username = ?", userInfo.Username).First(&tempUser).Error
	if canAddThisUser == nil {
		return errors.New("this username already exists")
	}

	// Add new user
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(userInfo.Password), bcrypt.DefaultCost)
	if err == nil {
		userInfo.Uuid = uuid.Must(uuid.NewRandom()).String() // UUID will uniquely idenfiy the user
		userInfo.Password = string(hashedPass)
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
	err := DB.Db.Where("username = ?", username).First(&tempUser).Error
	if err != nil {
		return ""
	}
	return tempUser.Uuid
}

func GetTodosForUser(userUuid string) ([]models.Todo, error) {
	todos := []models.Todo{}
	fmt.Println("Querying for todos")
	err := DB.Db.Where("user_uuid = ?", userUuid).Find(&todos).Error
	if err != nil {
		fmt.Println(err)
		return todos, err
	}
	return todos, nil
}
