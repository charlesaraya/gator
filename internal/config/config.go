package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

const (
	CONFIG_FILE string = ".gatorconfig.json"
)

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home dir: %w", err)
	}
	return filepath.Join(home, CONFIG_FILE), nil
}

type Config struct {
	DBUrl    string `json:"db_url"`
	UserName string `json:"current_user_name"`
}

func (c *Config) SetUser(name string) error {
	c.UserName = name
	data, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal confiration: %w", err)
	}
	configFile, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("failed to get config file path: %w", err)
	}
	if err = os.WriteFile(configFile, data, 0664); err != nil {
		return fmt.Errorf("failed to write configuration file: %w", err)
	}
	return nil
}

func Read() (Config, error) {
	config := Config{}
	configFile, err := getConfigFilePath()
	if err != nil {
		return config, fmt.Errorf("failed to get configuration file: %w", err)
	}
	data, err := os.ReadFile(configFile)
	if err != nil {
		return config, fmt.Errorf("failed to read configuration file: %w", err)
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}
	return config, nil
}

func LoadDB(cfg *Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}
	return db, nil
}
