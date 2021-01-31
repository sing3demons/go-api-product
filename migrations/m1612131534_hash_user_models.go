package migrations

import (
	"app/models"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func m1612131534HashUserModels() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1612131534",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.User{}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.DropTable("ussers").Error
		},
	}
}
