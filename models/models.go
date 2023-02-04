package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Uuid     string `json:"uuid" gorm:"uniqueIndex"`
	Username string `json:"username" validate:"omitempty,max=64" gorm:"uniqueIndex"`
	Password string `json:"password" validate:"omitempty,min=8,alphanum"`
}

type Task struct {
	gorm.Model
	TaskName  string `json:"task_name"`
	Completed bool   `gorm:"default:false" json:"completed"`
	UserUuid  string `json:"user_uuid"`
}
