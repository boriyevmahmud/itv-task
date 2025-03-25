package models

import (
	"time"

	"gorm.io/gorm"
)

type Movie struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	Title     string         `gorm:"type:varchar(255);not null;uniqueIndex:idx_movies_title"` // Unique index for title
	Director  string         `gorm:"type:varchar(255);not null;index:idx_movies_director"`    // Index for director
	Year      int            `gorm:"index:idx_movies_year"`                                   // Index for faster search by year
	Plot      string         `gorm:"type:text"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index:idx_movies_deleted_at"` // Index for soft deletes
}

type CreateMovieRequest struct {
	Title    string `json:"title" binding:"required,max=255" example:"Inception"`
	Director string `json:"director" binding:"required,max=255" example:"Christopher Nolan"`
	Year     int    `json:"year" binding:"required,gte=1888,lte=2025" example:"2010"`
	Plot     string `json:"plot" example:"A skilled thief is given a chance to erase his criminal past by performing an impossible task."`
}

type BulkInsertMoviesRequest struct {
	Movies []CreateMovieRequest `json:"movies" binding:"required,dive,required"`
}

type UpdateMovieRequest struct {
	ID       uint   `json:"-"`
	Title    string `json:"title" binding:"max=255" example:"Inception"`
	Director string `json:"director" binding:"max=255" example:"Christopher Nolan"`
	Year     int    `json:"year" binding:"gte=1888,lte=2025" example:"2010"`
	Plot     string `json:"plot" example:"A skilled thief is given a chance to erase his criminal past by performing an impossible task."`
}

type MovieResponse struct {
	ID        uint      `json:"id" example:"1"`
	Title     string    `json:"title" example:"Inception"`
	Director  string    `json:"director" example:"Christopher Nolan"`
	Year      int       `json:"year" example:"2010"`
	Plot      string    `json:"plot" example:"A skilled thief is given a chance to erase his criminal past by performing an impossible task."`
	CreatedAt time.Time `json:"created_at" example:"2025-03-22T15:04:05Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2025-03-22T15:04:05Z"`
}

type MovieListResponse struct {
	Movies []MovieResponse `json:"movies"`
	Count  int             `json:"count" example:"100"`
}
