package controllers

import (
	"github.com/sing3demons/app/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

func (c *Category) findCategoryByID(ctx *gin.Context) (*models.Category, error) {
	var category models.Category
	id := ctx.Param("id")

	if err := c.DB.Preload("Product").First(&category, id).Error; err != nil {
		return nil, err
	}

	return &category, nil
}
