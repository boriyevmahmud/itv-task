package services

import (
	"itv-task/internal/models"
	"itv-task/internal/repositories"

	"go.uber.org/fx"
)

type MovieService struct {
	repo *repositories.MovieRepository
}

func NewMovieService(repo *repositories.MovieRepository) *MovieService {
	return &MovieService{repo: repo}
}

// Provide the service to the Fx container
var MovieServiceModule = fx.Module("movieService",
	fx.Provide(NewMovieService),
)

func (s *MovieService) CreateMovie(movie *models.CreateMovieRequest) (uint, error) {
	return s.repo.Create(movie)
}

func (s *MovieService) GetMovieByID(id uint) (*models.MovieResponse, error) {
	return s.repo.GetByID(id)
}

func (s *MovieService) GetAllMovies(title, director string, year int, sortBy, sortOrder string, limit, offset int) (models.MovieListResponse, error) {
	return s.repo.GetAll(title, director, year, sortBy, sortOrder, limit, offset)
}

func (s *MovieService) UpdateMovie(movie *models.UpdateMovieRequest) error {
	return s.repo.Update(movie)
}

func (s *MovieService) DeleteMovie(id uint) error {
	return s.repo.Delete(id)
}