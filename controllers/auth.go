package controllers

import (
	"mime/multipart"
	"net/http"

	"github.com/sing3demons/app/models"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

//Auth - receiver adater
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
// GetProfile godoc
// @Summary get an user profile
// @Description get by form user
// @Tags auth
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} userResponse
// @Router /api/v1/auth/profile [get]
func (a *Auth) GetProfile(ctx *gin.Context) {
	sub, _ := ctx.Get("sub")
	user := sub.(*models.User)
	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})
}

// UpdateProfile godoc
// @Summary update an user
// @Description add by form user
// @Tags auth
// @Accept  mpfd
// @Produce  json
// @Security BearerAuth
// @Param name formData string true "name"
// @Param email formData string true "email"
// @Param avatar formData file true "avatar"
// @Success 200 {object} userResponse
// @Router /api/v1/auth/profile [put]
func (a *Auth) UpdateProfile(ctx *gin.Context) {
	form := updateProfileForm{}
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	sub, _ := ctx.Get("sub")
	user := sub.(*models.User)
	copier.Copy(&user, &form)

	img, err := setUserImage(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	user.Avatar = *img

	if err := a.DB.Save(&user).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})
}

// Register godoc
// @Summary add an user
// @Description add by form user
// @Tags auth
// @Accept  json
// @Produce  json
// @Param authForm body authForm true "authForm"
// @Success 201 {object} authResponse
// @Router /api/v1/auth/register [post]
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
