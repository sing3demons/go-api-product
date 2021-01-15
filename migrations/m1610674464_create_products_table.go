package migrations

import (
	"app/models"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func m1610674464CreateTableProductsTable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1610674464",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.Product{}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.DropTable("products").Error
		},
	}
}
