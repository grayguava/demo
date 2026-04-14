package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds all non-sensitive configuration.
// Tokens are never stored here — always read from env vars.
type Config struct {
	Provider      string `json:"provider"`
	CloudflareNS  string `json:"cloudflare_namespace"`
	SupabaseURL   string `json:"supabase_url"`
	SupabaseTable string `json:"supabase_table"`
	OutputFolder  string `json:"output_folder"`
	TokenEnv      string `json:"token_env"`
}

var configPath = filepath.Join(mustHomeDir(), ".formsealdaemon", "config.json")

func mustHomeDir() string {
	h, err := os.UserHomeDir()
	if err != nil {
		panic("cannot determine home directory: " + err.Error())
	}
	return h
}

func loadConfig() (*Config, error) {
	data, err := os.ReadFile(configPath)
	if os.IsNotExist(err) {
		return defaultConfig(), nil
	}
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func saveConfig(cfg *Config) error {
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

func defaultConfig() *Config {
	return &Config{
		Provider:     "cloudflare",
		OutputFolder: filepath.Join(mustHomeDir(), "FormSealData"),
		TokenEnv:     "FSYNC_CF_TOKEN",
	}
}

// token reads the API token from the env var named in cfg.TokenEnv.
func (cfg *Config) token() string {
	if cfg.TokenEnv == "" {
		return ""
	}
	return os.Getenv(cfg.TokenEnv)
}