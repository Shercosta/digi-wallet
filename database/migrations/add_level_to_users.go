package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AddLevelToUsers() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "addleveltousers", 
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec(`
				ALTER TABLE users 
				ADD COLUMN IF NOT EXISTS level INT DEFAULT 1 NOT NULL;
				UPDATE users SET level = 1 WHERE level IS NULL;
			`).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec("ALTER TABLE users DROP COLUMN level").Error
		},
	}
}
