package config

import (
	"fmt"
	//"os"

	"github.com/spf13/viper"
)

//! \struct Config
//! \brief Holds application configuration settings.
type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
}

//! \fn Load() (*Config, error)
//! \brief Loads configuration from environment and file.
//! \return Config instance and error (if any).
func Load() (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.SetEnvPrefix("APP")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	cfg := &Config{
		Port:        v.GetString("port"),
		DatabaseURL: v.GetString("database_url"),
		JWTSecret:   v.GetString("jwt_secret"),
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("database_url is required")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("jwt_secret is required")
	}

	return cfg, nil
}