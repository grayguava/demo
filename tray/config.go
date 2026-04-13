package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

var configFile = filepath.Join(homeDir(), ".formsealdaemon", "config.json")

type Config struct {
	Provider        string `json:"provider,omitempty"`
	CloudflareNS    string `json:"cloudflare.namespace,omitempty"`
	SupabaseURL     string `json:"supabase.url,omitempty"`
	SupabaseTable   string `json:"supabase.table,omitempty"`
	OutputFolder    string `json:"output_folder,omitempty"`
	SyncInterval    int    `json:"sync_interval,omitempty"`
}

func loadConfig() (*Config, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{SyncInterval: 15}, nil
		}
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func saveConfig(cfg *Config) error {
	os.MkdirAll(filepath.Dir(configFile), 0755)
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configFile, data, 0644)
}