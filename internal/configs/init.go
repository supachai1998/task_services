package configs

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	AppName      string
	Port         string
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
}

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
}

var AppConfig *Config

func InitConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// Set default values
	viper.SetDefault("READ_TIMEOUT", 5)
	viper.SetDefault("WRITE_TIMEOUT", 30)
	viper.SetDefault("IDLE_TIMEOUT", 120)

	AppConfig = &Config{
		Server: ServerConfig{
			AppName:      viper.GetString("APP_NAME"),
			Port:         viper.GetString("SERVER_PORT"),
			ReadTimeout:  viper.GetInt("READ_TIMEOUT"),
			WriteTimeout: viper.GetInt("WRITE_TIMEOUT"),
		},
		Database: DatabaseConfig{
			Driver:   viper.GetString("DB_DRIVER"),
			Host:     viper.GetString("POSTGRES_HOST"),
			Port:     viper.GetInt("POSTGRES_PORT"),
			User:     viper.GetString("POSTGRES_USER"),
			Password: viper.GetString("POSTGRES_PASSWORD"),
			DbName:   viper.GetString("POSTGRES_DB"),
		},
	}

	// Log the loaded configuration (optional)
	log.Printf("Configuration loaded: %+v\n", AppConfig)
}
