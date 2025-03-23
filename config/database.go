package config

import (
	"fmt"
	"itv-task/internal/models"
	"log"

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
	return db
}
