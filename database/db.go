package database

import (
	"log"
	"os"

	"github.com/sing3demons/app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DATABASE_URL")
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}
	database.AutoMigrate(&models.Product{})
	database.AutoMigrate(&models.Category{})
	database.AutoMigrate(&models.User{})

	db = database
}

func GetDB() *gorm.DB {
	return db
}
