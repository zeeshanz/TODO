package repos

import (
	"fmt"

	"github.com/zeeshanz/TODO/database"
	"github.com/zeeshanz/TODO/models"
)

func CreateTodo(uuid string, userUuid string, completed bool, todoItem string) (*models.Todo, error) {
	todo := &models.Todo{
		Uuid:      uuid,
		UserUuid:  userUuid,
		Completed: completed,
		TodoItem:  todoItem,
	}

	if err := database.DB.Db.Create(&todo).Error; err != nil {
		return nil, err
	}
	return todo, nil

}

/*
 * Retrieve Todo item's id
 */
func getTodoId(uuid string) uint {
	var TodoItem models.Todo
	err := database.DB.Db.Where("uuid = ?", uuid).First(&TodoItem).Error
	if err != nil {
		return 0
	}
	return TodoItem.ID
}

/*
 * Retrieve a Todo item by uuid
 */
func GetTodoItem(uuid string) (models.Todo, error) {
	var todoItem models.Todo
	err := database.DB.Db.Where("uuid = ?", uuid).First(&todoItem).Error
	if err != nil {
		return todoItem, err
	}
	return todoItem, nil
}

/*
 * Retrieve all Todos for a given user where user is identified by its uuid
 */
func GetAllTodos(userUuid string) ([]models.Todo, error) {
	todos := []models.Todo{}
	fmt.Println("Querying for todos")
	err := database.DB.Db.Where("user_uuid = ?", userUuid).Find(&todos).Error
	if err != nil {
		fmt.Println(err)
		return todos, err
	}
	return todos, nil
}

/*
 * Mark complete status as true of false for a given Todo item
 */
func UpdateTodoStatus(todoUuid string, completed bool) error {
	todo := new(models.Todo)
	err := database.DB.Db.Model(&todo).Select("Completed").Where("uuid = ?", todoUuid).Updates(models.Todo{Completed: completed})
	if err != nil {
		fmt.Println(err)
		return err.Error
	}
	return nil
}

/*
 * Update the todo item with new text
 */
func UpdateTodoItem(todoUuid string, updatedTodo string) error {
	todo := new(models.Todo)
	err := database.DB.Db.Model(&todo).Where("uuid = ?", todoUuid).Updates(models.Todo{TodoItem: updatedTodo})
	if err != nil {
		fmt.Println(err)
		return err.Error
	}
	return nil
}

/*
 * Update the todo item with new text
 */
func DeleteTodo(todoUuid string) error {
	todoId := getTodoId(todoUuid) // Getting the primary key because it is used in the Delete function.
	fmt.Printf("Deleting todo with id: %v\n", todoId)
	var todo models.Todo
	err := database.DB.Db.Delete(&todo, todoId)

	if err != nil {
		fmt.Println(err)
		return err.Error
	}
	return nil
}
