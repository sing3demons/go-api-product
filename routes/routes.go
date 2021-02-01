package routes

import (
	"app/config"
	"app/controllers"
	"app/middleware"

	"github.com/gin-gonic/gin"
)

//Serve - middleware
func Serve(r *gin.Engine) {
	db := config.GetDB()
	v1 := r.Group("/api/v1")

	userGroup := v1.Group("/auth")
	userController := controllers.Auth{DB: db}
	{
		userGroup.POST("/register", userController.Register)
		userGroup.POST("/login", middleware.Authenticate().LoginHandler)
	}

	productGroup := v1.Group("/products")
	productController := controllers.Product{DB: db}
	{
		productGroup.GET("", productController.FindAll)
		productGroup.GET("/:id", productController.FindOne)
		productGroup.POST("", productController.Create)
		productGroup.PATCH("/:id", productController.Update)
		productGroup.DELETE("/:id", productController.Delete)
	}

	categoryGroup := v1.Group("/categories")
	categoryController := controllers.Category{DB: db}
	{
		categoryGroup.GET("", categoryController.FindAll)
		categoryGroup.GET("/:id", categoryController.FindOne)
		categoryGroup.POST("", categoryController.Create)
		categoryGroup.PATCH("/:id", categoryController.Update)
		categoryGroup.DELETE("/:id", categoryController.Delete)
	}
}
