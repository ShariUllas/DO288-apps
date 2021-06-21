package config

import (
	"os"
	"strings"
)

// Config holds env data
type Config struct {
	DBUserName            string
	DBPassword            string
	DBHost                string
	DBPort                string
	DBName                string
	Port                  string
	DatadogEnableAPM      bool
	DatadogEnableProfiler bool
	DatadogServiceName    string
}

// NewConfig initializes config
func NewConfig() *Config {
	var config Config
	config.DBUserName = os.Getenv("POSTGRES_USER")
	config.DBPassword = os.Getenv("POSTGRES_PASSWORD")
	config.DBHost = os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	if !strings.HasPrefix(dbPort, ":") {
		config.DBPort = ":" + dbPort
	}
	config.DBName = os.Getenv("POSTGRES_DB")
	apiPort := os.Getenv("API_PORT")
	if !strings.HasPrefix(apiPort, ":") {
		config.Port = ":" + apiPort
	}
	config.DatadogEnableAPM = true
	config.DatadogEnableProfiler = true
	config.DatadogServiceName = "catalogue"
	return &config
}
