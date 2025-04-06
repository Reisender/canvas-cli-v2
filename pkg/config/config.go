package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config contains Canvas API configuration
type Config struct {
	APIKey  string `mapstructure:"api_key"`
	BaseURL string `mapstructure:"base_url"`
}

// Global config instance
var (
	AppConfig Config
)

// InitConfig initializes the configuration
func InitConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error finding home directory:", err)
		return
	}

	// Config file path
	configDir := filepath.Join(home, ".config", "canvas-cli")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Println("Error creating config directory:", err)
		return
	}

	// Set up viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	// Set defaults
	viper.SetDefault("base_url", "https://canvas.instructure.com/api/v1")

	// Read config from file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, create it
			if err := viper.SafeWriteConfig(); err != nil {
				fmt.Println("Error creating config file:", err)
			}
		} else {
			fmt.Println("Error reading config file:", err)
		}
	}

	// Bind environment variables
	viper.SetEnvPrefix("CANVAS")
	viper.BindEnv("api_key")
	viper.BindEnv("base_url")

	// Unmarshal config
	if err := viper.Unmarshal(&AppConfig); err != nil {
		fmt.Println("Error parsing config:", err)
	}
}

// SaveConfig saves the current configuration
func SaveConfig() error {
	return viper.WriteConfig()
}

// GetConfig returns the current config
func GetConfig() Config {
	return AppConfig
}

// UpdateConfig updates the configuration with new values
func UpdateConfig(key string, value string) error {
	viper.Set(key, value)
	AppConfig = Config{}
	if err := viper.Unmarshal(&AppConfig); err != nil {
		return err
	}
	return SaveConfig()
}
