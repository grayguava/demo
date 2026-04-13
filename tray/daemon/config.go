package daemon

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Provider        string `json:"provider"`
	CloudflareNS    string `json:"cloudflare_namespace"`
	SupabaseURL     string `json:"supabase_url"`
	SupabaseTable   string `json:"supabase_table"`
	OutputFolder    string `json:"output_folder"`
	SyncInterval    int    `json:"sync_interval"`
	APIToken        string `json:"api_token"`
}

var configFile = filepath.Join(userHome(), ".formsealdaemon", "config.json")

func userHome() string {
	h, _ := os.UserHomeDir()
	return h
}

func LoadConfig() (*Config, error) {
	os.MkdirAll(filepath.Dir(configFile), 0755)
	data, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{
				SyncInterval: 15,
				OutputFolder: filepath.Join(userHome(), "FormSeal"),
			}, nil
		}
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func SaveConfig(cfg *Config) error {
	os.MkdirAll(filepath.Dir(configFile), 0755)
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configFile, data, 0644)
}