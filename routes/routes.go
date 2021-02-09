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
	
	authenticate := middleware.Authenticate().MiddlewareFunc()

	authGroup := v1.Group("/auth")
	authController := controllers.Auth{DB: db}
	{
		authGroup.POST("/register", authController.Register)
		authGroup.POST("/login", middleware.Authenticate().LoginHandler)
		authGroup.GET("/profile", authenticate, authController.GetProfile)
		authGroup.PATCH("/profile", authenticate, authController.UpdateProfile)
	}

	usersController := controllers.Users{DB: db}
	usersGroup := v1.Group("users")
	usersGroup.Use(authenticate)
	{
		usersGroup.GET("", usersController.FindAll)
		usersGroup.POST("", usersController.Create)
		usersGroup.GET("/:id", usersController.FindOne)
		usersGroup.PATCH("/:id", usersController.Update)
		usersGroup.DELETE("/:id", usersController.Delete)
		usersGroup.PATCH("/:id/promote", usersController.Promote)
		usersGroup.PATCH("/:id/demote", usersController.Demote)
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
