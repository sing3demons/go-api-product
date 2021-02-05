package models

import "github.com/jinzhu/gorm"

//Category - model
type Category struct {
	gorm.Model
	Name    string `gorm:"unique;not null"`
	Product []Product
}
