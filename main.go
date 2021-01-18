package main

import (
	"app/config"
	"app/migrations"
	"app/routes"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.InitDB()
	defer config.CloseDB()
	migrations.Migrate()

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
