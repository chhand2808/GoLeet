package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	APIKey string `json:"api_key"`
}

var configFilePath = filepath.Join(os.Getenv("HOME"), ".goleet_config.json")

func SaveAPIKey(key string) error {
	config := Config{APIKey: key}

	file, err := os.Create(configFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(config)
}
