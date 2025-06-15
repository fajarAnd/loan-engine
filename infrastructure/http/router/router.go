package router

import (
	"net/http"
	"path/filepath"
	"strings"

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
	loanRepo := repositories.NewLoanRepository(db)
	fileRepo := repositories.NewFileRepository(db)

	// Usecases
	jwtSecret := viper.GetString("jwt.secret")
	authUsecase := usecase.NewAuthUsecase(authRepo, jwtSecret)
	loanUsecase := usecase.NewLoanUsecase(loanRepo)
	fileUsecase := usecase.NewFileUsecase(fileRepo)

	// Controllers
	authController := controller.NewAuthController(authUsecase)
	loanController := controller.NewLoanController(loanUsecase)
	fileController := controller.NewFileController(fileUsecase)

	// Static file serving for uploaded documents
	workDir, _ := filepath.Abs(".")
	filesDir := http.Dir(filepath.Join(workDir, "uploads"))
	FileServer(r, "/uploads", filesDir)

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
					r.Post("/files/upload", fileController.UploadSurveyDocument)
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
				r.Post("/loans", loanController.CreateLoanProposal)
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

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
