package repositories

import (
	"itv-task/internal/models"
	"log"

	"gorm.io/gorm"
)

type MovieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) *MovieRepository {
	return &MovieRepository{db: db}
}

func (r *MovieRepository) Create(movie *models.Movie) error {
	if err := r.db.Create(movie).Error; err != nil {
		log.Println("❌ Failed to create movie:", err)
		return err
	}
	return nil
}

func (r *MovieRepository) GetByID(id uint) (*models.Movie, error) {
	var movie models.Movie
	if err := r.db.First(&movie, id).Error; err != nil {
		log.Println("❌ Movie not found:", err)
		return nil, err
	}
	return &movie, nil
}

func (r *MovieRepository) GetAll(title, director string, year int, sortBy string, limit, offset int) ([]models.Movie, error) {
	var movies []models.Movie
	query := r.db

	if title != "" {
		query = query.Where("title ILIKE ?", "%"+title+"%")
	}
	if director != "" {
		query = query.Where("director ILIKE ?", "%"+director+"%")
	}
	if year > 0 {
		query = query.Where("year = ?", year)
	}

	if sortBy == "title" {
		query = query.Order("title ASC")
	} else if sortBy == "year" {
		query = query.Order("year DESC")
	} else {
		query = query.Order("id ASC")
	}

	query = query.Limit(limit).Offset(offset)

	if err := query.Find(&movies).Error; err != nil {
		log.Println("❌ Failed to retrieve movies:", err)
		return nil, err
	}
	return movies, nil
}

func (r *MovieRepository) Update(movie *models.Movie) error {
	if err := r.db.Save(movie).Error; err != nil {
		log.Println("❌ Failed to update movie:", err)
		return err
	}
	return nil
}

func (r *MovieRepository) Delete(id uint) error {
	if err := r.db.Where("id = ?", id).Delete(&models.Movie{}).Error; err != nil {
		log.Println("❌ Failed to soft delete movie:", err)
		return err
	}
	return nil
}
