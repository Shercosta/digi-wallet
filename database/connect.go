package database

import (
	"fmt"
	"log"
	"os"

	"github.com/Shercosta/digi-wallet/database/migrations"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect() *gorm.DB {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	fmt.Println("connected to database")

	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		migrations.CreateBalanceTable(),
		migrations.CreateUserTable(),
		migrations.AddUserIDToBalance(),
		migrations.AddLevelToUser(),
	})

	if err := m.Migrate(); err != nil {
		log.Fatal("failed to migrate database", err)
	}

	log.Println("database migrated")

	return db
}
