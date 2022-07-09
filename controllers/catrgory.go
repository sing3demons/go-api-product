package controllers

import (
	"net/http"

	"github.com/sing3demons/app/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

// FindAll godoc
// @Summary Show an categories
// @Description get by form categories
// @Tags categories
// @Accept  json
// @Produce  json
// @Success 200 {object} []allCategoryResponse
// @Router /api/v1/categories [get]
func (c *Category) FindAll(ctx *gin.Context) {
	var categories []models.Category

	if err := c.DB.Preload("Product").Find(&categories).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if s := ctx.Query("name"); s != "" {
		if err := c.DB.Where("name = ?", s).Find(&categories).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

	}

	serializedCategory := []allCategoryResponse{}
	copier.Copy(&serializedCategory, &categories)
	ctx.JSON(http.StatusOK, gin.H{"category": serializedCategory})
}

// FindOne godoc
// @Summary Show an category
// @Description get by form category
// @Tags categories
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 200 {object} allCategoryResponse
// @Router /api/v1/categories/{id} [get]
func (c *Category) FindOne(ctx *gin.Context) {
	category, err := c.findCategoryByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var serializedCategory categoryResponse
	copier.Copy(&serializedCategory, &category)
	ctx.JSON(http.StatusOK, gin.H{"category": serializedCategory})
}

// Create godoc
// @Summary add an category
// @Description add by form category
// @Tags categories
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param name body createCategoryForm true "name"
// @Success 201 {object} allCategoryResponse
// @Router /api/v1/categories [post]
func (c *Category) Create(ctx *gin.Context) {
	var form createCategoryForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var categories models.Category
	copier.Copy(&categories, &form)

	if err := c.DB.Create(&categories).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	serializedCategory := allCategoryResponse{}
	copier.Copy(&serializedCategory, &categories)
	ctx.JSON(http.StatusCreated, gin.H{"category": serializedCategory})
}

// UpdateAll godoc
// @Summary update an category
// @Description update by form category
// @Tags categories
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path string true "id"
// @Param name body updateCategoryForm true "name"
// @Success 200 {object} allCategoryResponse
// @Router /api/v1/categories/{id} [patch]
func (c *Category) Update(ctx *gin.Context) {
	var form updateCategoryForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	category, err := c.findCategoryByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	copier.Copy(&category, &form)
	if err := c.DB.Save(&category).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error})
		return
	}

	var serializedProduct allCategoryResponse
	copier.Copy(&serializedProduct, &category)
	ctx.JSON(http.StatusOK, gin.H{"category": serializedProduct})

}

// Delete godoc
// @Summary	delete an category
// @Description	delete by json category
// @Tags	categories
// @Accept	json
// @Produce	json
// @Security BearerAuth
// @Param id path string true "id"
// @Success 204
// @Failure	422  {object} string "Bad Request"
// @Failure	404  {object}  map[string]any	"{"error": "not found"}"
// @Router /api/v1/categories/{id} [delete]
func (c *Category) Delete(ctx *gin.Context) {
	category, err := c.findCategoryByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.DB.Unscoped().Delete(&category)

	ctx.Status(http.StatusNoContent)
}
