package models

import (
	"time"

	"gorm.io/gorm"
)

type Movie struct {
	gorm.Model
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     string         `gorm:"type:varchar(255);not null" json:"title"`
	Director  string         `gorm:"type:varchar(255);not null" json:"director"`
	Year      int            `gorm:"not null" json:"year"`
	Plot      string         `gorm:"type:text" json:"plot"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
