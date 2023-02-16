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

type UserDTO struct {
	Uuid     string `json:"uuid" gorm:"uniqueIndex"`
	Username string `json:"username" validate:"omitempty,max=64" gorm:"uniqueIndex"`
	Password string `json:"password" validate:"omitempty,min=8,alphanum"`
}

type UserResponse struct {
	Uuid     string `json:"uuid" gorm:"uniqueIndex"`
	Username string `json:"username" validate:"omitempty,max=64" gorm:"uniqueIndex"`
	Password string `json:"password" validate:"omitempty,min=8,alphanum"`
}
