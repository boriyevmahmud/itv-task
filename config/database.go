package config

import (
	"fmt"
	"itv-task/internal/models"
	"log"
	"time"

	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// DatabaseModule for Uber FX
var DatabaseModule = fx.Module("database",
	fx.Provide(NewDatabase),
)

// NewDatabase initializes the database connection
func NewDatabase(cfg *Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%v sslmode=disable",
		cfg.PostgresHost, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDatabase, cfg.PostgresPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	db.AutoMigrate(models.Movie{})

	log.Println("✅ Connected to database")
	DB = db

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ Failed to get SQL DB instance: %v", err)
	}

	sqlDB.SetMaxOpenConns(25)                 // Max open connections (tune based on DB capacity)
	sqlDB.SetMaxIdleConns(10)                 // Max idle connections (reduces resource usage)
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // Time a connection can be reused
	return db
}
