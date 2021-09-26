package migrations

import (
	"github.com/sing3demons/app/models"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)


func m1612129431CreactUserable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1612129431",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.User{}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.DropTable("ussers").Error
		},
	}
}