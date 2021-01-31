package migrations

import (
	"app/config"
	"log"

	"gopkg.in/gormigrate.v1"
)

func Migrate() {
	db := config.GetDB()
	m := gormigrate.New(
		db,
		gormigrate.DefaultOptions,
		[]*gormigrate.Migration{
			m1610674464CreateTableProductsTable(),
			m1610950017CreactCategoryTable(),
			m1612129431CreactUserable(),
		},
	)

	if err := m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
}
