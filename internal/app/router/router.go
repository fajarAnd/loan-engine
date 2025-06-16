package router

import (
	"github.com/fajar-andriansyah/loan-engine/internal/app/constants"
	"github.com/fajar-andriansyah/loan-engine/internal/app/controllers"
	"github.com/fajar-andriansyah/loan-engine/internal/app/database"
	"github.com/fajar-andriansyah/loan-engine/internal/app/middleware"
	repositories2 "github.com/fajar-andriansyah/loan-engine/internal/app/repositories"
	usecase2 "github.com/fajar-andriansyah/loan-engine/internal/app/usecase"
	"github.com/fajar-andriansyah/loan-engine/internal/pkg/pdf"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"
)

func GetRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)

	db := database.GetConn()

	pdfGenerator := pdf.NewPDFGenerator()

	// Repositories
	authRepo := repositories2.NewAuthRepository(db)
	loanRepo := repositories2.NewLoanRepository(db)
	fileRepo := repositories2.NewFileRepository(db)
	investmentRepo := repositories2.NewInvestmentRepository(db)

	// Usecases
	jwtSecret := viper.GetString("jwt.secret")
	authUsecase := usecase2.NewAuthUsecase(authRepo, jwtSecret)
	loanUsecase := usecase2.NewLoanUsecase(loanRepo, pdfGenerator)
	fileUsecase := usecase2.NewFileUsecase(fileRepo)
	investmentUsecase := usecase2.NewInvestmentUsecase(investmentRepo, pdfGenerator)

	// Controllers
	authController := controller.NewAuthController(authUsecase)
	loanController := controller.NewLoanController(loanUsecase)
	fileController := controller.NewFileController(fileUsecase)
	investmentController := controller.NewInvestmentController(investmentUsecase)

	workDir, _ := filepath.Abs(".")
	filesDir := http.Dir(filepath.Join(workDir, "uploads"))
	FileServer(r, "/uploads", filesDir)

	// Routes
	r.Get("/__health", controller.GetHealth)

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public auth routes
		r.Post("/auth/login", authController.Login)

		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuthMiddleware())

			// Employee only routes
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireUserType(constants.USER_EMPLOYEE))

				// Field validator routes
				r.Group(func(r chi.Router) {
					r.Use(middleware.RequireRole(constants.ROLE_FIELD_VALIDATOR))
					r.Post("/files/upload", fileController.UploadSurveyDocument)
				})

				// Field officer routes
				r.Group(func(r chi.Router) {
					r.Use(middleware.RequireRole(constants.ROLE_FIELD_OFFICER))
					r.Put("/loans/{id}/approve", loanController.ApproveLoan)
					r.Put("/loans/{id}/disburse", loanController.DisburseLoan)
				})

			})

			// Borrower routes
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireUserType(constants.USER_BORROWER))
				r.Post("/loans", loanController.CreateLoanProposal)
			})

			// Investor routes
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireUserType(constants.USER_INVESTOR))
				r.Post("/loans/{id}/investments", investmentController.CreateInvestment)
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
