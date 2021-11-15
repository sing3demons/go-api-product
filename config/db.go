package config

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

//InitDB - connenct database
func InitDB() {
	var err error
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	db_port := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")
	// db, err := gorm.Open("sqlite3", "./tmp/gorm.db")
	DATABASE_URL := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s  sslmode=disable TimeZone=Asia/Bangkok", dbHost, user, password, dbname, db_port)

	// fmt.Print("url: ", DATABASE_URL)
	db, err = gorm.Open("postgres", DATABASE_URL)
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
