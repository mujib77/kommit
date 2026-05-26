package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Provider string
	APIKey   string
	Model	string
	Style	string
	MaxLength int
	Language string
}

func Load() Config {
	home, _ := os.UserHomeDir()
	configDir := filepath.Join(home, ".kommit")
	os.MkdirAll(configDir, 0755)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	viper.SetDefault("provider", "gemini")
	viper.SetDefault("style", "conventional")
	viper.SetDefault("max_length", 72)
	viper.SetDefault("language", "english")

	viper.ReadInConfig()

	apiKey := viper.GetString("api_key")
if apiKey == "" {
    apiKey = os.Getenv("OPENAI_API_KEY")
}
if apiKey == "" {
    apiKey = os.Getenv("ANTHROPIC_API_KEY")
}
if apiKey == "" {
    apiKey = os.Getenv("GEMINI_API_KEY")
}

	return Config{
		Provider: viper.GetString("provider"),
		APIKey:   apiKey,
		Model:    viper.GetString("model"),
		Style:    viper.GetString("style"),
		MaxLength: viper.GetInt("max_length"),
		Language: viper.GetString("language"),
	}
}

func ConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".kommit", "config.yaml")
}

func PrintSetup() {
	fmt.Printf(`
kommit setup

Create a config file at: %s

Example config:
  provider: openai
  api_key: sk-
  model: gpt-4o
  style: conventional
  language: english
 
Or set environment variables:
  export OPENAI_API_KEY=sk-...
  export ANTHROPIC_API_KEY=sk-...
  
 `, ConfigPath())
}