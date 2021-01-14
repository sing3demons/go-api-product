package routes

import (
	"kp-app/controllers"

	"github.com/gin-gonic/gin"
)

//Serve - middleware
func Serve(r *gin.Engine) {

	productGroup := r.Group("/api/v1/products")
	productController := controllers.Product{}

	{
		productGroup.GET("", productController.FindAll)
		productGroup.GET("/:id", productController.FindOne)
		productGroup.POST("", productController.Update)
	}
}
