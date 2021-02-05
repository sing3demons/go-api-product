package controllers

import (
	"app/models"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

type Auth struct {
	DB *gorm.DB
}

type authForm struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type updateProfileForm struct {
	Name   string                `form:"name"`
	Email  string                `form:"email" `
	Avatar *multipart.FileHeader `form:"avatar"`
}

type authResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name" `
}

//GetProfile - /auth/profile => JWT => sub (UserID) => User => User
func (a *Auth) GetProfile(ctx *gin.Context) {
	sub, _ := ctx.Get("sub")
	user := sub.(*models.User)
	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})
}

//UpdateProfile - patch
func (a *Auth) UpdateProfile(ctx *gin.Context) {
	var form updateProfileForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	sub, _ := ctx.Get("sub")
	user := sub.(*models.User)

	setUserImage(ctx, user)
	if err := a.DB.Model(user).Update(&form).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})
}

//Register - sign up
func (a *Auth) Register(ctx *gin.Context) {
	var form authForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error(), "message": "ลงทะเบียนไม่สำเร็จ"})
		return
	}
	var user models.User
	copier.Copy(&user, &form)
	user.Password = user.GenerateEncryptedPassword()
	if err := a.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error(), "message": "ลงทะเบียนไม่สำเร็จ"})
		return
	}

	var serializedUser authResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusCreated, gin.H{"user": serializedUser, "message": "ลงทะเบียนสำเร็จ"})
}
