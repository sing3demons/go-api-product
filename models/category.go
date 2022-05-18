package models

import "gorm.io/gorm"

//Category - model
type Category struct {
	gorm.Model
	Name    string `gorm:"unique;not null"`
	Product []Product
}
