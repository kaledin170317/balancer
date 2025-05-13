package config

import (
	"balancer/internal/logger"
	"flag"
	"fmt"
)

func LoadFinalConfig() (*Config, error) {
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to config file (JSON or YAML)")
	flag.Parse()

	cmdCfg := ParseFromFlags()
	if isValidConfig(cmdCfg) {
		logger.Info(nil, "Loaded config from command-line flags")
		return cmdCfg, nil
	}

	if configPath != "" {
		fileCfg, err := LoadFromFile(configPath)
		if err != nil {
			logger.Error(nil, "Error loading config from file", "err", err, "path", configPath)
			return nil, fmt.Errorf("error loading config from file: %w", err)
		}
		logger.Info(nil, "Loaded config from file", "path", configPath)
		return fileCfg, nil
	}

	envCfg, err := LoadFromEnv()
	if err != nil {
		logger.Error(nil, "Error loading config from environment", "err", err)
		return nil, fmt.Errorf("error loading config from env: %w", err)
	}
	logger.Info(nil, "Loaded config from environment variables")
	return envCfg, nil
}

func isValidConfig(cfg *Config) bool {
	return cfg.ListenAddr != "" || len(cfg.Backends) > 0 || cfg.Algorithm != ""
}
