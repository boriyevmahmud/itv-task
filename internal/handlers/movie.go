package handlers

import (
	"itv-task/internal/models"
	"itv-task/internal/services"
	"itv-task/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MovieHandler struct {
	service *services.MovieService
}

func NewMovieHandler(service *services.MovieService) *MovieHandler {
	return &MovieHandler{service: service}
}

// @Security ApiKeyAuth
// CreateMovie creates a new movie
// @Summary Create a new movie
// @Description Add a new movie to the database
// @Tags movies
// @Accept json
// @Produce json
// @Param movie body models.CreateMovieRequest true "Movie data"
// @Success 201 {object} models.CreateMovieRequest
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /movies [post]
func (h *MovieHandler) CreateMovie(c *gin.Context) {
	var movie models.CreateMovieRequest
	if err := c.ShouldBindJSON(&movie); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request", "Failed to parse request body")
		return
	}

	if len(movie.Title) == 0 || len(movie.Title) > 255 {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid title", "Title is required and must be <= 255 characters")
		return
	}
	if len(movie.Director) == 0 || len(movie.Director) > 255 {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid director", "Director is required and must be <= 255 characters")
		return
	}
	if movie.Year < 1888 || movie.Year > 2025 {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid year", "Year must be between 1888 and 2025")
		return
	}

	if _, err := h.service.CreateMovie(&movie); err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Internal server error", "Failed to create movie")
		return
	}

	c.JSON(http.StatusCreated, movie)
}

// GetAllMovies retrieves all movies
// @Summary Get all movies
// @Description Retrieve a list of movies with optional filters
// @Tags movies
// @Produce json
// @Param title query string false "Filter by title"
// @Param director query string false "Filter by director"
// @Param year query int false "Filter by year"
// @Param limit query int false "Limit results"
// @Param offset query int false "Offset results"
// @Param sort_by query string false "Sort by field (title, year, created_at, director)"
// @Param sort_order query string false "Sort order (asc, desc)"
// @Success 200 {array} models.MovieListResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /movies [get]
func (h *MovieHandler) GetAllMovies(c *gin.Context) {
	title := c.Query("title")
	director := c.Query("director")
	yearStr := c.Query("year")
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")
	sortBy := c.Query("sort_by")
	sortOrder := c.Query("sort_order")

	var year, limit, offset int
	var err error

	if yearStr != "" {
		year, err = strconv.Atoi(yearStr)
		if err != nil || year < 1888 || year > 2025 {
			utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid year", "Year must be between 1888 and 2025")
			return
		}
	}
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid limit", "Limit must be a positive number")
			return
		}
	} else {
		limit = 10
	}
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid offset", "Offset must be a non-negative number")
			return
		}
	} else {
		offset = 0
	}
	if sortBy != "" {
		if sortBy != "title" && sortBy != "year" && sortBy != "created_at" && sortBy != "director" {
			utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid sort_by", "Invalid sort_by value")
			return
		}
	}
	if sortOrder != "" {
		if sortOrder != "asc" && sortOrder != "desc" {
			utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid sort_order", "Invalid sort_order value")
			return
		}
	}

	movies, err := h.service.GetAllMovies(title, director, year, sortBy, sortOrder, limit, offset)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Internal server error", "Failed to retrieve movies")
		return
	}

	c.JSON(http.StatusOK, movies)
}

// GetMovieByID retrieves a single movie by ID
// @Summary Get a movie by ID
// @Description Retrieve a movie using its ID
// @Tags movies
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} models.MovieResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /movies/{id} [get]
func (h *MovieHandler) GetMovieByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid ID", "Movie ID must be a positive integer")
		return
	}

	movie, err := h.service.GetMovieByID(uint(id))
	if err != nil {
		if err.Error() == "record not found" {
			utils.SendErrorResponse(c, http.StatusNotFound, "Movie not found", "No movie found with the given ID")
		} else {
			utils.SendErrorResponse(c, http.StatusInternalServerError, "Internal server error", "Failed to retrieve movie")
		}
		return
	}

	c.JSON(http.StatusOK, movie)
}

// @Security ApiKeyAuth
// UpdateMovie updates an existing movie
// @Summary Update a movie
// @Description Modify an existing movie
// @Tags movies
// @Accept json
// @Produce json
// @Param id path int true "Movie ID"
// @Param movie body models.UpdateMovieRequest true "Updated movie data"
// @Success 200 {object} models.UpdateMovieRequest
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /movies/{id} [put]
func (h *MovieHandler) UpdateMovie(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid ID", "Movie ID must be a positive integer")
		return
	}

	var movie models.UpdateMovieRequest
	if err := c.ShouldBindJSON(&movie); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request", "Failed to parse request body")
		return
	}

	_, err = h.service.GetMovieByID(uint(id))
	if err != nil {
		if err.Error() == "record not found" {
			utils.SendErrorResponse(c, http.StatusNotFound, "Movie not found", "No movie found with the given ID")
		} else {
			utils.SendErrorResponse(c, http.StatusInternalServerError, "Internal server error", "Failed to retrieve movie")
		}
		return
	}

	if len(movie.Title) == 0 || len(movie.Title) > 255 {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid title", "Title is required and must be <= 255 characters")
		return
	}
	if len(movie.Director) == 0 || len(movie.Director) > 255 {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid director", "Director is required and must be <= 255 characters")
		return
	}
	if movie.Year < 1888 || movie.Year > 2025 {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid year", "Year must be between 1888 and 2025")
		return
	}
	movie.ID = uint(id)

	if err := h.service.UpdateMovie(&movie); err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Internal server error", "Failed to update movie")
		return
	}
	c.JSON(http.StatusOK, movie)
}

// @Security ApiKeyAuth
// DeleteMovie deletes a movie by ID
// @Summary Delete a movie
// @Description Remove a movie from the database
// @Tags movies
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /movies/{id} [delete]
func (h *MovieHandler) DeleteMovie(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid ID", "Movie ID must be a positive integer")
		return
	}

	if err := h.service.DeleteMovie(uint(id)); err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Internal server error", "Failed to delete movie")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie deleted"})
}
