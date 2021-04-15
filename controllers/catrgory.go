package controllers

import (
	"app/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

//Category - Method Receiver
type Category struct {
	DB *gorm.DB
}

type categoryResponse struct {
	Name    string `json:"name"`
	Product []struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"products"`
}

type allCategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type createCategoryForm struct {
	Name string `json:"name" binding:"required"`
}

type updateCategoryForm struct {
	Name string `json:"name"`
}

// FindAll - query all categories
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

// FindOne - first query
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

// Create - create
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

//Update - update --> patch
func (c *Category) Update(ctx *gin.Context) {
	var form updateCategoryForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	category, err := c.findCategoryByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := c.DB.Model(&category).Update(&form).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error})
		return
	}

	var serializedProduct allCategoryResponse
	copier.Copy(&serializedProduct, &category)
	ctx.JSON(http.StatusOK, gin.H{"category": serializedProduct})

}

// Delete - remove category
func (c *Category) Delete(ctx *gin.Context) {
	category, err := c.findCategoryByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.DB.Unscoped().Delete(&category)

	ctx.Status(http.StatusNoContent)
}

func (c *Category) findCategoryByID(ctx *gin.Context) (*models.Category, error) {
	var category models.Category
	id := ctx.Param("id")

	if err := c.DB.Preload("Product").First(&category, id).Error; err != nil {
		return nil, err
	}

	return &category, nil
}
