package router

import (
	_healthController "github.com/fajar-andriansyah/loan-engine/controllers"
	"github.com/go-chi/chi/v5"
)

func GetRouter() chi.Router {
	r := chi.NewRouter()

	// Middleware
	// TODO: Set Middleware

	// Routes
	r.Get("/__health", _healthController.GetHealth)

	return r
}
