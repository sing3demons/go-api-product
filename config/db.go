package config

import (
	"app/models"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

//InitDB - connenct database
func InitDB() {
	var err error
	// db, err := gorm.Open("sqlite3", "./tmp/gorm.db")
	db, err = gorm.Open("postgres", os.Getenv("DATABASE_CONNECTION"))
	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(gin.Mode() == gin.ReleaseMode)
	db.AutoMigrate(&models.Product{})
}

//GetDB - return db
func GetDB() *gorm.DB {
	return db
}

// CloseDB - close database
func CloseDB() {
	db.Close()
}
