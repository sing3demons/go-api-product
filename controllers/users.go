package controllers

import (
	"net/http"

	"github.com/sing3demons/app/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

//FindAll - api/v1/users @GET
// FindAll godoc
// @Summary get an users
// @Description get by form users
// @Tags user
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} []userResponse
// @Router /api/v1/users [get]
func (u *Users) FindAll(ctx *gin.Context) {
	var users []models.User
	query := u.DB.Order("id desc").Find(&users)

	pagination := pagination{ctx: ctx, query: query, records: &users}
	paging := pagination.paginate()

	var serializedUsers []userResponse
	copier.Copy(&serializedUsers, &users)
	ctx.JSON(http.StatusOK, gin.H{
		"users": usersPaging{Items: serializedUsers, Paging: paging},
	})
}

//Create - api/v1/users @POST
// Create godoc
// @Summary add an user
// @Description add by form user
// @Tags user
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param createUserForm body createUserForm true "create-user-form"
// @Success 201 {object} userResponse
// @Router /api/v1/users [post]
func (u *Users) Create(ctx *gin.Context) {
	var form createUserForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	copier.Copy(&user, &form)
	user.Password = user.GenerateEncryptedPassword()

	if err := u.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusCreated, gin.H{"user": serializedUser})
}

//FindOne - api/v1/users/:id @GET
// FindOne godoc
// @Summary get an user
// @Description get by form user
// @Tags user
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path int true "ID"
// @Success 200 {object} userResponse
// @Router /api/v1/users/{id} [get]
func (u *Users) FindOne(ctx *gin.Context) {
	user, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})
}

//Update - api/v1/users/:id @PATCH
// Update godoc
// @Summary update an user
// @Description update by form user
// @Tags user
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path int true "ID"
// @Param updateUserForm body updateUserForm false "update-user-form"
// @Success 200 {object} userResponse
// @Router /api/v1/users/{id} [patch]
func (u *Users) Update(ctx *gin.Context) {
	var form updateUserForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	user, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	copier.Copy(&user, &form)

	if form.Password != "" {
		user.Password = user.GenerateEncryptedPassword()
	}

	if err := u.DB.Save(&user).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})

}

//Delete - api/v1/users/:id @DELETE
// Delete godoc
// @Summary delete an user
// @Description delete by form user
// @Tags user
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path int true "ID"
// @Success 204
// @Router /api/v1/users/{id} [delete]
func (u *Users) Delete(ctx *gin.Context) {
	user, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	u.DB.Unscoped().Delete(&user)

	ctx.Status(http.StatusNoContent)
}

//Promote - api/v1/users/:id/promote @PATCH
// Promote godoc
// @Summary update an user
// @Description update by form user
// @Tags user
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path int true "ID"
// @Success 200 {object} userResponse
// @Router /api/v1/users/{id}/promote [patch]
func (u *Users) Promote(ctx *gin.Context) {
	user, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	user.Promote()
	u.DB.Save(&user)

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})
}

//Demote - api/v1/users/:id/demote @PATCH
// Demote godoc
// @Summary update an user
// @Description update by form user
// @Tags user
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path int true "ID"
// @Success 200 {object} userResponse
// @Router /api/v1/users/{id}/demote [patch]
func (u *Users) Demote(ctx *gin.Context) {
	user, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	user.Demote()
	u.DB.Save(&user)

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})
}
