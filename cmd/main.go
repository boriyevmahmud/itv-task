package main

import (
	"context"
	"fmt"
	"itv-task/config"
	"itv-task/internal/handlers"
	"itv-task/internal/repositories"
	"itv-task/internal/services"
	"itv-task/pkg/logger"
	utils "itv-task/pkg/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"

	_ "itv-task/docs" // Swagger documentation
)

// NewRouter initializes the Gin router
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func NewRouter(cfg *config.Config, movieHandler *handlers.MovieHandler, authHandler *handlers.AuthHandler) *gin.Engine {
	r := gin.Default()

	// Middleware
	r.Use(gin.Recovery())     // Handles panics
	r.Use(utils.AuthLogger()) // Example logging middleware

	// Public Routes
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("swagger/doc.json")))
	r.GET("/movies", movieHandler.GetAllMovies)
	r.GET("/movies/:id", movieHandler.GetMovieByID)
	r.POST("/auth/login", authHandler.Login)
	r.POST("/auth/refresh", authHandler.RefreshToken)

	// Protected Routes (Require Auth)
	authRoutes := r.Group("/movies")
	fmt.Println("cfg: ", cfg.JWTAccessSecret)
	authRoutes.Use(utils.AuthMiddleware(cfg)) // Apply token validation
	{
		authRoutes.POST("/", movieHandler.CreateMovie)
		authRoutes.PUT("/:id", movieHandler.UpdateMovie)
		authRoutes.DELETE("/:id", movieHandler.DeleteMovie)
		authRoutes.POST("/bulk-insert", movieHandler.BulkInsertMovies)
	}

	return r
}

// StartServer starts the HTTP server with Uber FX lifecycle
func StartServer(lc fx.Lifecycle, router *gin.Engine) {
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Println("🚀 Server is running on port 8080")
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatalf("❌ Server error: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("⚠️ Shutting down server...")
			shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			return server.Shutdown(shutdownCtx)
		},
	})
}

func main() {
	app := fx.New(
		fx.Provide(
			func() *config.Config { cfg := config.Load(); return &cfg }, // Provide config first
		),
		config.DatabaseModule, // Ensure database module comes after config
		fx.Provide(
			func() logger.Logger { log := logger.New("itv", "Movies"); return log },
			repositories.NewMovieRepository,
			services.NewMovieService,
			handlers.NewMovieHandler,
			services.NewAuthService,
			handlers.NewAuthHandler,
			NewRouter,
		),
		fx.Invoke(StartServer), // Start server
	)

	// Handle OS Signals for Graceful Shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		log.Println("🛑 Application shutting down...")
		app.Stop(context.Background())
	}()

	app.Run()
}
