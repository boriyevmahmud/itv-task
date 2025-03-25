package services

import (
	"itv-task/internal/models"
	"itv-task/internal/repositories"
	"itv-task/pkg/logger"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MovieService struct {
	repo *repositories.MovieRepository
	log  logger.Logger
}

func NewMovieService(repo *repositories.MovieRepository, log logger.Logger) *MovieService {
	return &MovieService{repo: repo, log: log}
}

// Provide the service to the Fx container
var MovieServiceModule = fx.Module("movieService",
	fx.Provide(NewMovieService),
)

func (s *MovieService) CreateMovie(movie *models.CreateMovieRequest) (uint, error) {
	s.log.Info("Creating movie", zap.Any("request", movie))
	id, err := s.repo.Create(movie)
	if err != nil {
		s.log.Error("Failed to create movie", zap.Any("request", movie), zap.Error(err))
		return 0, err
	}

	return id, nil
}

func (s *MovieService) GetMovieByID(id uint) (*models.MovieResponse, error) {
	s.log.Info("getting movie", zap.Any("request", id))

	movie, err := s.repo.GetByID(id)
	if err != nil {
		s.log.Error("Failed to fetch movie", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}

	return movie, nil
}

func (s *MovieService) GetAllMovies(title, director string, year int, sortBy, sortOrder string, limit, offset int) (models.MovieListResponse, error) {
	s.log.Info("Getting movies", zap.Any("request", map[string]interface{}{
		"title":     title,
		"director":  director,
		"year":      year,
		"sortBy":    sortBy,
		"sortOrder": sortOrder,
		"limit":     limit,
		"offset":    offset}))
	movies, err := s.repo.GetAll(title, director, year, sortBy, sortOrder, limit, offset)
	if err != nil {
		s.log.Error("Failed to fetch movies", zap.Any("request", map[string]interface{}{
			"title":     title,
			"director":  director,
			"year":      year,
			"sortBy":    sortBy,
			"sortOrder": sortOrder,
			"limit":     limit,
			"offset":    offset,
		}), zap.Error(err))
		return models.MovieListResponse{}, err
	}

	return movies, nil
}

func (s *MovieService) UpdateMovie(movie *models.UpdateMovieRequest) error {
	s.log.Info("Updating movie", zap.Any("request", movie))

	err := s.repo.Update(movie)
	if err != nil {
		s.log.Error("Failed to update movie", zap.Any("request", movie), zap.Error(err))
		return err
	}

	s.log.Info("Movie updated successfully", zap.Uint("id", movie.ID))
	return nil
}

func (s *MovieService) DeleteMovie(id uint) error {
	s.log.Info("Deleting movie", zap.Uint("request", id))

	err := s.repo.Delete(id)
	if err != nil {
		s.log.Error("Failed to delete movie", zap.Uint("id", id), zap.Error(err))
		return err
	}

	return nil
}

func (s *MovieService) BulkInsertMovies(movies *models.BulkInsertMoviesRequest) error {
	s.log.Info("Creating movie", zap.Any("request", movies))
	err := s.repo.BulkInsertMovies(movies)
	if err != nil {
		s.log.Error("Failed to create movie", zap.Any("request", movies), zap.Error(err))
		return err
	}

	return nil
}

func (s *MovieService) GetMovieByTitle(title string) (*models.MovieResponse, error) {
	s.log.Info("getting movie by title", zap.Any("request", title))

	movie, err := s.repo.GetByTitle(title)
	if err != nil {
		s.log.Error("Failed to fetch movie by titkle", zap.String("title", title), zap.Error(err))
		return nil, err
	}

	return movie, nil
}
