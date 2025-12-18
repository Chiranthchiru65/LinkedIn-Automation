package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// Config holds all configuration for our bot
// We nest structs to match the YAML structure (Target, Connection, Settings)
type Config struct {
	Credentials Credentials // From .env
	Target      Target      `yaml:"target"`     // From yaml
	Connection  Connection  `yaml:"connection"` // From yaml
	Settings    Settings    `yaml:"settings"`   // From yaml
}

type Credentials struct {
	Username string
	Password string
}

type Target struct {
	Keywords string `yaml:"keywords"`
	Location string `yaml:"location"`
}

type Connection struct {
	MessageTemplate string `yaml:"message_template"`
	LimitPerDay     int    `yaml:"limit_per_day"`
}

type Settings struct {
	Headless    bool `yaml:"headless"`
	StealthMode bool `yaml:"stealth_mode"`
}

// LoadConfig reads .env and config.yaml and returns a populated Config struct
func LoadConfig() (*Config, error) {
	// 1. Load .env file variables into the system environment
	// We use "godotenv" library for this.
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: Error loading .env file (might be in production environment)")
	}

	// 2. Read the YAML file
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read config.yaml: %w", err)
	}

	// 3. Parse YAML into the Config struct
	var cfg Config
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config.yaml: %w", err)
	}

	// 4. Manually load credentials from environment variables
	// We do this separately because passwords should NEVER be in the YAML file
	cfg.Credentials.Username = os.Getenv("LINKEDIN_USERNAME")
	cfg.Credentials.Password = os.Getenv("LINKEDIN_PASSWORD")

	// 5. Basic Validation
	if cfg.Credentials.Username == "" || cfg.Credentials.Password == "" {
		return nil, fmt.Errorf("missing credentials in .env file")
	}

	return &cfg, nil
}