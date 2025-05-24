package config

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

const (
	CONFIG_FILE string = ".gatorconfig.json"
)

func getConfigFile() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
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
		return err
	}
	configFile, err := getConfigFile()
	if err != nil {
		return err
	}
	if err = os.WriteFile(configFile, data, 0664); err != nil {
		return err
	}
	return nil
}

func Read() (Config, error) {
	config := Config{}
	configFile, err := getConfigFile()
	if err != nil {
		return config, err
	}
	data, err := os.ReadFile(configFile)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

func LoadDB(cfg *Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		return nil, err
	}
	return db, nil
}
