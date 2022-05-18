package routes

import (
	"github.com/sing3demons/app/cache"
	"github.com/sing3demons/app/database"
	"github.com/sing3demons/app/controllers"
	"github.com/sing3demons/app/middleware"

	"github.com/gin-gonic/gin"
)

func NewCacherConfig() *cache.CacherConfig {
	return &cache.CacherConfig{}
}

//Serve - middleware
func Serve(r *gin.Engine) {
	db := database.GetDB()
	cacher := cache.NewCacher(NewCacherConfig())
	v1 := r.Group("/api/v1")

	authenticate := middleware.Authenticate().MiddlewareFunc()
	authorize := middleware.Authorize()

	authGroup := v1.Group("/auth")
	authController := controllers.Auth{DB: db}
	{
		authGroup.POST("/register", authController.Register)
		authGroup.POST("/login", middleware.Authenticate().LoginHandler)
		authGroup.GET("/profile", authenticate, authController.GetProfile)
		authGroup.PUT("/profile", authenticate, authController.UpdateProfile)
	}

	usersController := controllers.Users{DB: db}
	usersGroup := v1.Group("users")
	usersGroup.Use(authenticate, authorize)
	{
		usersGroup.GET("", usersController.FindAll)
		usersGroup.POST("", usersController.Create)
		usersGroup.GET("/:id", usersController.FindOne)
		usersGroup.PATCH("/:id", usersController.Update)
		usersGroup.DELETE("/:id", usersController.Delete)
		usersGroup.PATCH("/:id/promote", usersController.Promote)
		usersGroup.PATCH("/:id/demote", usersController.Demote)
	}

	productController := controllers.Product{DB: db, Cacher: cacher}
	productGroup := v1.Group("/products")
	productGroup.GET("", productController.FindAll)
	productGroup.GET("/:id", productController.FindOne)

	productGroup.Use(authenticate, authorize)
	{
		productGroup.POST("", productController.Create)
		productGroup.PATCH("/:id", productController.Update)
		productGroup.DELETE("/:id", productController.Delete)
	}

	categoryController := controllers.Category{DB: db}
	categoryGroup := v1.Group("/categories")
	categoryGroup.GET("", categoryController.FindAll)
	categoryGroup.GET("/:id", categoryController.FindOne)
	categoryGroup.Use(authenticate, authorize)
	{
		categoryGroup.POST("", categoryController.Create)
		categoryGroup.PATCH("/:id", categoryController.Update)
		categoryGroup.DELETE("/:id", categoryController.Delete)
	}
}
