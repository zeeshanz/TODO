package models

import (
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Uuid      string `json:"uuid" gorm:"uniqueIndex"`
	TodoItem  string `json:"todo_item"`
	Completed bool   `gorm:"default:false" json:"completed"`
	UserUuid  string `json:"user_uuid"`
}

type TodoDTO struct {
	Uuid      string `json:"uuid" gorm:"uniqueIndex"`
	TodoItem  string `json:"todo_item"`
	Completed bool   `gorm:"default:false" json:"completed"`
	UserUuid  string `json:"user_uuid"`
}

type TodoResponse struct {
	Uuid      string `json:"uuid" gorm:"uniqueIndex"`
	TodoItem  string `json:"todo_item"`
	Completed bool   `gorm:"default:false" json:"completed"`
	UserUuid  string `json:"user_uuid"`
}
