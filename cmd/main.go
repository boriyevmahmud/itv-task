package main

import (
	"itv-task/config"
	"itv-task/internal/handlers"
	"itv-task/internal/repositories"
	"itv-task/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	cfg := config.Load()
	config.InitDB(&cfg)

	movieRepo := repositories.NewMovieRepository(config.DB)
	movieService := services.NewMovieService(movieRepo)
	movieHandler := handlers.NewMovieHandler(movieService)

	r := gin.Default()

	// Movie Routes
	r.POST("/movies", movieHandler.CreateMovie)
	r.GET("/movies", movieHandler.GetAllMovies)
	r.GET("/movies/:id", movieHandler.GetMovieByID)
	r.PUT("/movies/:id", movieHandler.UpdateMovie)
	r.DELETE("/movies/:id", movieHandler.DeleteMovie)

	url := ginSwagger.URL("swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// Start Server
	r.Run(":8080")
}
