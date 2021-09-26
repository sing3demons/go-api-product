package migrations

import (
	"github.com/sing3demons/app/models"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func m11613010616UpdateCollumProducts() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1613010616",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.Product{}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.DropTable("products").Error
		},
	}
}
