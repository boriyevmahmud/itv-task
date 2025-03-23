package services

import (
	"itv-task/internal/models"
	"itv-task/internal/repositories"
)

type MovieService struct {
	repo *repositories.MovieRepository
}

func NewMovieService(repo *repositories.MovieRepository) *MovieService {
	return &MovieService{repo: repo}
}

func (s *MovieService) CreateMovie(movie *models.Movie) error {
	return s.repo.Create(movie)
}

func (s *MovieService) GetMovieByID(id uint) (*models.Movie, error) {
	return s.repo.GetByID(id)
}

func (s *MovieService) GetAllMovies(title, director string, year int, sortBy string, limit, offset int) ([]models.Movie, error) {
	return s.repo.GetAll(title, director, year, sortBy, limit, offset)
}

func (s *MovieService) UpdateMovie(movie *models.Movie) error {
	return s.repo.Update(movie)
}

func (s *MovieService) DeleteMovie(id uint) error {
	return s.repo.Delete(id)
}
