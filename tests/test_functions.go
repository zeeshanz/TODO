package tests

import (
	"github.com/zeeshanz/TODO/models"
	"gorm.io/gorm"
)

func FindUser(userTemp *models.User, user string, db1 *gorm.DB) (err error) {
	err1 := db1.Where("username = ?", user).First(userTemp).Error
	return err1
}

func CreateUser(userUuid string, username string, password string, db1 *gorm.DB) (user *models.User, err error) {
	user1 := &models.User{Uuid: userUuid, Username: username, Password: password}
	err1 := db1.Create(user1).Error
	if err != nil {
		return nil, err1
	}
	return user1, err1
}

func CreateTodo(uuid string, todoItem string, completed bool, userUuid string, db1 *gorm.DB) (todo *models.Todo, err error) {
	todo1 := &models.Todo{Uuid: uuid, TodoItem: todoItem, Completed: completed, UserUuid: userUuid}
	err1 := db1.Create(todo1).Error
	if err != nil {
		return nil, err1
	}
	return todo1, err1
}

func FindUserTasks(user *models.User, usertask *[]models.Todo, db1 *gorm.DB) (err error) {
	return db1.Model(user).Association("Todos").Find(usertask)
}
