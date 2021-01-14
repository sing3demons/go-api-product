package controllers

import (
	"kp-app/models"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Product struct {
	ID uint
}

type createProductForm struct {
	Title string                `form:"title" binding:"required"`
	Body  string                `form:"body" binding:"required"`
	Image *multipart.FileHeader `form:"image" binding:"required"`
}

var products []models.Product = []models.Product{
	{ID: 1, Title: "Title#1", Body: "Body#1"},
	{ID: 2, Title: "Title#2", Body: "Body#2"},
	{ID: 3, Title: "Title#3", Body: "Body#3"},
}

func (p *Product) FindAll(ctx *gin.Context) {
	result := products

	if limit := ctx.Query("limit"); limit != "" {
		n, _ := strconv.Atoi(limit)

		result = result[:n]
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

func (p *Product) FindOne(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	for _, item := range products {
		if item.ID == uint(id) {
			ctx.JSON(http.StatusOK, gin.H{"data": item})
			return
		}
	}
	ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
}

func (p *Product) Update(ctx *gin.Context) {
	var form createProductForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := models.Product{
		ID:    uint(len(products) + 1),
		Title: form.Title,
		Body:  form.Body,
	}

	// Get file
	file, _ := ctx.FormFile("image")

	// Create file
	path := "uploads/products/" + strconv.Itoa(int(product.ID)) // ID => 8, uploads/articles/8/image.png
	os.MkdirAll(path, 0755)                                     // -> uploads/products/8

	// Upload file
	filename := path + "/" + file.Filename
	if err := ctx.SaveUploadedFile(file, filename); err != nil {
		log.Fatal(err.Error())
	}

	// Attach file to products
	product.Image = os.Getenv("HOST") + "/" + filename

	products = append(products, product)

	ctx.JSON(http.StatusOK, gin.H{"message": product})
}
