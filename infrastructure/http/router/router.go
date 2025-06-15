package router

import (
	"github.com/fajar-andriansyah/loan-engine/controllers"
	"github.com/fajar-andriansyah/loan-engine/infrastructure/database"
	"github.com/fajar-andriansyah/loan-engine/infrastructure/http/middleware"
	"github.com/fajar-andriansyah/loan-engine/repositories"
	"github.com/fajar-andriansyah/loan-engine/usecase"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"
)

func GetRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)

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
		// Public auth routes
		r.Post("/auth/login", authController.Login)

		// Protected routes (require authentication)
		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuthMiddleware())

			// Employee only routes
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireUserType("employee"))

				// Field validator routes
				r.Group(func(r chi.Router) {
					r.Use(middleware.RequireRole("FIELD_VALIDATOR"))
					// TODO: Add field validator specific routes
					// r.Post("/loans/{id}/survey", surveyController.UploadSurvey)
				})

				// Field officer routes
				r.Group(func(r chi.Router) {
					r.Use(middleware.RequireRole("FIELD_OFFICER"))
					// TODO: Add field officer specific routes
					// r.Put("/loans/{id}/approve", loanController.ApproveLoan)
					// r.Put("/loans/{id}/disburse", loanController.DisburseLoan)
				})

			})

			// Borrower routes
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireUserType("borrower"))
				// TODO: Add borrower specific routes
				// r.Post("/loans", loanController.CreateLoan)
				// r.Get("/loans/{id}", loanController.GetLoanDetails)
			})

			// Investor routes
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireUserType("investor"))
				// TODO: Add investor specific routes
				// r.Get("/loans/available", loanController.GetAvailableLoans)
				// r.Post("/loans/{id}/investments", investmentController.MakeInvestment)
				// r.Get("/portfolio", investmentController.GetPortfolio)
			})
		})

	})

	return r
}
