package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AddLevelToUser() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "addleveltouser",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec(`ALTER TABLE users ADD COLUMN level INT NOT NULL DEFAULT 1`).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec(`ALTER TABLE users DROP COLUMN level`).Error
		},
	}
}
