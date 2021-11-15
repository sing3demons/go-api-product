package main

import (
	"log"
	"os"

	"github.com/sing3demons/app/config"
	"github.com/sing3demons/app/migrations"
	"github.com/sing3demons/app/routes"
	"github.com/sing3demons/app/seeds"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file")
		}
	}

	config.InitDB()
	defer config.CloseDB()
	migrations.Migrate()
	seeds.Load()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization")

	r := gin.Default()
	r.Use(cors.New(corsConfig)) //cors
	r.Static("/uploads", "./uploads")

	//สร้าง folder
	uploadDirs := [...]string{"products", "users"}
	for _, dir := range uploadDirs {
		os.MkdirAll("uploads/"+dir, 0755)
	}

	routes.Serve(r)

	r.Run(":" + os.Getenv("PORT"))
}
