package migrations

import (
	"github.com/sing3demons/app/models"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func m1610950017CreactCategoryTable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1610950017",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.Category{}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.DropTable("categories").Error
		},
	}
}
