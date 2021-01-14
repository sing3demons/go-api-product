package routes

import (
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type product struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
	Image string `json:"image"`
}

type createProductForm struct {
	Title string                `form:"title" binding:"required"`
	Body  string                `form:"body" binding:"required"`
	Image *multipart.FileHeader `form:"image" binding:"required"`
}

//Serve - middleware
func Serve(r *gin.Engine) {
	products := []product{
		{ID: 1, Title: "Title#1", Body: "Body#1"},
		{ID: 2, Title: "Title#2", Body: "Body#2"},
		{ID: 3, Title: "Title#3", Body: "Body#3"},
	}

	productGroup := r.Group("/api/v1/products")

	productGroup.GET("", func(ctx *gin.Context) {
		result := products

		if limit := ctx.Query("limit"); limit != "" {
			n, _ := strconv.Atoi(limit)

			result = result[:n]
		}

		ctx.JSON(http.StatusOK, gin.H{"data": result})
	})

	//:id
	productGroup.GET("/:id", func(ctx *gin.Context) {
		id, _ := strconv.Atoi(ctx.Param("id"))
		for _, item := range products {
			if item.ID == uint(id) {
				ctx.JSON(http.StatusOK, gin.H{"data": item})
				return
			}
		}
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
	})

	productGroup.POST("", func(ctx *gin.Context) {
		var form createProductForm
		if err := ctx.ShouldBind(&form); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		p := product{
			ID:    uint(len(products) + 1),
			Title: form.Title,
			Body:  form.Body,
		}

		// Get file
		file, _ := ctx.FormFile("image")

		// Create file
		path := "uploads/products/" + strconv.Itoa(int(p.ID)) // ID => 8, uploads/articles/8/image.png
		os.MkdirAll(path, 0755)                               // -> uploads/products/8

		// Upload file
		filename := path + "/" + file.Filename
		if err := ctx.SaveUploadedFile(file, filename); err != nil {
			log.Fatal(err.Error())
		}

		// Attach file to products
		p.Image = "http://localhost:8080/" + filename

		products = append(products, p)

		ctx.JSON(http.StatusOK, gin.H{"message": p})
	})
}
