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
	doesExist := DB.Db.Where("username = ?", userInfo.Username).First(&tempUser).Error
	if doesExist != nil {
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
