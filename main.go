package main

import (
	db "github.com/fajar-andriansyah/loan-engine/internal/app/database"
	"github.com/fajar-andriansyah/loan-engine/internal/app/router"
	"github.com/spf13/viper"
	"net/http"
	"os"

	"github.com/fajar-andriansyah/loan-engine/config"
	"github.com/rs/zerolog/log"
)

func main() {
	_ = config.LoadConfig()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := db.InitDB(viper.GetString("database.dsn")); err != nil {
		log.Warn().Err(err).Msg("")
	}

	router := router.GetRouter()

	// Start server
	addr := ":" + port
	log.Info().Msgf("Health check: http://localhost%s/__health", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Error().Msgf("Server failed to start: %v", err)
	}
}
