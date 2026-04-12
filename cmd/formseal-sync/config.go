package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type ProviderConfig struct {
	Token     string `json:"token"`
	Namespace string `json:"namespace"`
	URL       string `json:"url"`
	Key       string `json:"key"`
	Table     string `json:"table"`
}

type Config struct {
	Provider      string          `json:"provider"`
	Cloudflare    ProviderConfig `json:"cloudflare"`
	Supabase      ProviderConfig `json:"supabase"`
	OutputFolder  string          `json:"output_folder"`
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".formseal-sync", "config.json"), nil
}

func LoadConfig() Config {
	path, err := configPath()
	if err != nil {
		return Config{}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}
	}
	return cfg
}

func SaveConfig(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}