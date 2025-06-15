package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"
)

var ServiceVersion = "development"

const defaultPort = 3000

func LoadConfig() error {

	var err error
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(".conf")
	viper.AutomaticEnv()
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Info().Msg("config reload")
	})
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		log.Error().Err(err).Msg("Fatal error config file")
	}

	// set default
	viper.SetDefault("http.port", defaultPort)
	viper.SetDefault("log.level", 0)
	viper.Set("http.timeout", "120s")

	lvl, _ := zerolog.ParseLevel(viper.GetString("log.level"))

	if lvl == zerolog.NoLevel {
		lvl = zerolog.ErrorLevel
	}

	zerolog.SetGlobalLevel(lvl)
	return nil
}
