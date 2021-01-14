package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type product struct {
	ID    uint
	Title string
	Body  string
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
}
