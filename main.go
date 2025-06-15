package main

import (
	"net/http"
	"os"

	"github.com/fajar-andriansyah/loan-engine/config"
	"github.com/fajar-andriansyah/loan-engine/infrastructure/http/router"
	"github.com/rs/zerolog/log"
)

func main() {
	_ = config.LoadConfig()

	// Get port dari environment atau default ke 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Setup routes
	router := router.GetRouter()

	// Start server
	addr := ":" + port
	log.Info().Msgf("🚀 OpenAI Vision Service starting on %s", addr)
	log.Info().Msgf("📚 API Documentation available at http://localhost%s/", addr)
	log.Info().Msgf("🔍 Health check at http://localhost%s/__health", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Error().Msgf("Server failed to start: %v", err)
	}
}
