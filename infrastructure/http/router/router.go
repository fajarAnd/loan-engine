package router

import (
	"github.com/fajar-andriansyah/loan-engine/controllers"
	"github.com/fajar-andriansyah/loan-engine/infrastructure/database"
	"github.com/fajar-andriansyah/loan-engine/repositories"
	"github.com/fajar-andriansyah/loan-engine/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
)

func GetRouter() chi.Router {
	r := chi.NewRouter()

	// Initialize dependencies
	db := database.GetConn()

	// Repositories
	authRepo := repositories.NewAuthRepository(db)

	// Usecases
	jwtSecret := viper.GetString("jwt.secret")
	authUsecase := usecase.NewAuthUsecase(authRepo, jwtSecret)

	// Controllers
	authController := controller.NewAuthController(authUsecase)

	// Routes
	r.Get("/__health", controller.GetHealth)

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Auth routes
		r.Post("/auth/login", authController.Login)
	})

	return r
}
