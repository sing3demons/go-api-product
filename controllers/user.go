package controllers

import "github.com/jinzhu/gorm"

type Users struct {
	DB *gorm.DB
}

type createUserForm struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

type updateUserForm struct {
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=8"`
	Name     string `json:"name"`
}

type userResponse struct {
	ID     uint   `json:"id"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
	Name   string `json:"name"`
	Role   string `json:"role"`
}