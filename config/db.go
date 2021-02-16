package config

import (
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
	db, err = gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(gin.Mode() == gin.ReleaseMode)
	// db.AutoMigrate(&models.Product{})
	// db.DropTable("users", "migrations")
}

//GetDB - return db
func GetDB() *gorm.DB {
	return db
}

// CloseDB - close database
func CloseDB() {
	db.Close()
}
