package repositories

import (
	"itv-task/internal/models"
	"log"
	"time"

	"gorm.io/gorm"
)

type MovieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) *MovieRepository {
	return &MovieRepository{db: db}
}

func (r *MovieRepository) Create(movie *models.CreateMovieRequest) (uint, error) {
	gormModel := models.Movie{
		Title:    movie.Title,
		Director: movie.Director,
		Year:     movie.Year,
		Plot:     movie.Plot,
	}
	if err := r.db.Table("movies").Create(&gormModel).Error; err != nil {
		log.Println("❌ Failed to create movie:", err)
		return 0, err
	}
	return gormModel.ID, nil
}

func (r *MovieRepository) GetByID(id uint) (*models.MovieResponse, error) {
	var movie models.MovieResponse
	if err := r.db.Table("movies").First(&movie, "id = ? AND deleted_at IS NULL ", id).Error; err != nil {
		log.Println("❌ Movie not found:", err)
		return nil, err
	}
	return &movie, nil
}

func (r *MovieRepository) GetByTitle(title string) (*models.MovieResponse, error) {
	var movie models.MovieResponse
	if err := r.db.Table("movies").First(&movie, "title = ? AND deleted_at IS NULL ", title).Error; err != nil {
		log.Println("❌ Movie not found:", err)
		return nil, err
	}
	return &movie, nil
}
func (r *MovieRepository) GetAll(title, director string, year int, sortBy, sortOrder string, limit, offset int) (models.MovieListResponse, error) {
	var movies []models.MovieResponse
	var totalCount int64
	query := r.db.Model(&models.Movie{})

	// Apply filters
	if title != "" {
		query = query.Where("title ILIKE ?", "%"+title+"%")
	}
	if director != "" {
		query = query.Where("director ILIKE ?", "%"+director+"%")
	}
	if year > 0 {
		query = query.Where("year = ?", year)
	}

	// Get total count before applying limit & offset
	if err := query.Count(&totalCount).Error; err != nil {
		log.Println("❌ Failed to count movies:", err)
		return models.MovieListResponse{}, err
	}

	switch sortOrder {
	case "asc":
		sortOrder = " AS "
	case "desc":
		sortOrder = " DESC "
	default:
		sortOrder = " DESC "
	}

	// Apply sorting
	switch sortBy {
	case "title":
		query = query.Order("title " + sortBy)
	case "year":
		query = query.Order("year " + sortOrder)
	case "created_at":
		query = query.Order("created_at " + sortOrder)
	case "director":
		query = query.Order("director " + sortBy)
	default:
		query = query.Order("id " + sortBy)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	query = query.Offset(offset)

	if err := query.Find(&movies).Error; err != nil {
		log.Println("❌ Failed to retrieve movies:", err)
		return models.MovieListResponse{}, err
	}

	return models.MovieListResponse{
		Movies: movies,
		Count:  int(totalCount),
	}, nil
}

func (r *MovieRepository) Update(movie *models.UpdateMovieRequest) error {
	gormModel := models.Movie{
		Title:     movie.Title,
		Director:  movie.Director,
		Year:      movie.Year,
		Plot:      movie.Plot,
		UpdatedAt: time.Now(),
	}
	if err := r.db.Table("movies").Save(gormModel).Error; err != nil {
		log.Println("❌ Failed to update movie:", err)
		return err
	}
	return nil
}

func (r *MovieRepository) Delete(id uint) error {
	if err := r.db.Table("movies").Where("id = ?", id).Delete(&models.Movie{}).Error; err != nil {
		log.Println("❌ Failed to soft delete movie:", err)
		return err
	}
	return nil
}

func (s *MovieRepository) BulkInsertMovies(movies *models.BulkInsertMoviesRequest) error {
	tx := s.db.Table("movies").Begin() // Start transaction
	if tx.Error != nil {
		log.Println("❌ Failed to start transaction:", tx.Error)
		return tx.Error
	}
	defer func() {
		tx.Rollback()
	}()

	var gormModel = models.Movie{}
	for _, movie := range movies.Movies {
		gormModel = models.Movie{
			Title:    movie.Title,
			Director: movie.Director,
			Year:     movie.Year,
			Plot:     movie.Plot,
		}
		if err := tx.Create(&gormModel).Error; err != nil {
			tx.Rollback() // Rollback on failure
			log.Println("❌ Failed to bulk insert movies:", err)
			return err
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		log.Println("❌ Failed to commit transaction:", err)
		return err
	}

	return nil
}
