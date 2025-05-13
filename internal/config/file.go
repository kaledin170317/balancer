package config

import (
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

func LoadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	switch ext := strings.ToLower(filepath.Ext(path)); ext {
	case ".json":
		err = json.Unmarshal(data, &cfg)
	case ".yml", ".yaml":
		err = yaml.Unmarshal(data, &cfg)
	default:
		err = errors.New("unsupported config format: " + ext)
	}
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
