// config/config.go
package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port    string
		Mode    string // development, production
		TimeOut int    // in seconds
	}
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
		SSLMode  string
		TimeZone string
	}
	OAuth struct {
		Google struct {
			ClientID     string
			ClientSecret string
			RedirectURL  string
		}
	}
	JWT struct {
		Secret        string
		ExpiresInHrs  int
		RefreshSecret string
	}
	Redis struct {
		Host     string
		Port     string
		Password string
		DB       int
	}
}

// LoadConfig reads configuration from file or environment variables
func LoadConfig() (*Config, error) {
	var config Config

	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name

	// Path to look for the config file
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("$HOME/.appname")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set default values
	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; use defaults and environment variables
			fmt.Println("No config file found. Using environment variables and defaults.")
		} else {
			// Config file was found but another error was produced
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Unmarshal config into struct
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	// Validate required fields
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "development")
	viper.SetDefault("server.timeout", 30)

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.timezone", "UTC")

	// JWT defaults
	viper.SetDefault("jwt.expiresinhrs", 24)

	// Redis defaults
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", "6379")
	viper.SetDefault("redis.db", 0)
}

func validateConfig(config *Config) error {
	// Database validation
	if config.Database.User == "" {
		return fmt.Errorf("database user is required")
	}
	if config.Database.Password == "" {
		return fmt.Errorf("database password is required")
	}
	if config.Database.DBName == "" {
		return fmt.Errorf("database name is required")
	}

	// OAuth validation
	if config.OAuth.Google.ClientID == "" {
		return fmt.Errorf("google OAuth client ID is required")
	}
	if config.OAuth.Google.ClientSecret == "" {
		return fmt.Errorf("google OAuth client secret is required")
	}
	if config.OAuth.Google.RedirectURL == "" {
		return fmt.Errorf("google OAuth redirect URL is required")
	}

	// JWT validation
	if config.JWT.Secret == "" {
		return fmt.Errorf("JWT secret is required")
	}

	return nil
}

// GetDSN returns the formatted database connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
		c.Database.TimeZone,
	)
}
