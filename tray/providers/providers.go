package providers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type ProviderConfig struct {
	Token     string
	Namespace string // Cloudflare
	URL       string // Supabase
	Table     string // Supabase
}

type Result struct {
	Written int
	Skipped int
}

var registry = make(map[string]func() Provider)

type Provider interface {
	Name() string
	Fetch(cfg *ProviderConfig, outputPath string) (written int, skipped int, err error)
}

func Register(name string, factory func() Provider) {
	registry[name] = factory
}

func Get(name string) Provider {
	if factory, ok := registry[name]; ok {
		return factory()
	}
	return nil
}

func List() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}

func init() {
	Register("cloudflare", func() Provider {
		return &CloudflareProvider{client: &http.Client{Timeout: 30 * time.Second}}
	})
	Register("supabase", func() Provider {
		return &SupabaseProvider{client: &http.Client{Timeout: 30 * time.Second}}
	})
}

type CloudflareProvider struct {
	client *http.Client
}

func (p *CloudflareProvider) Name() string {
	return "cloudflare"
}

func (p *CloudflareProvider) Fetch(cfg *ProviderConfig, outputPath string) (int, int, error) {
	if cfg.Token == "" || cfg.Namespace == "" {
		return 0, 0, fmt.Errorf("cloudflare: missing token or namespace")
	}

	// Cloudflare API: list keys in KV namespace
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/kv/namespaces/%s/keys", os.Getenv("CF_ACCOUNT_ID"), cfg.Namespace)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, 0, err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, 0, fmt.Errorf("cloudflare API error: %s", body)
	}

	var result struct {
		Result []struct {
			Name      string `json:"name"`
			Modified  string `json:"modified_on"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, 0, err
	}

	// Read existing keys
	existing := make(map[string]bool)
	if data, err := os.ReadFile(outputPath); err == nil {
		var lines []string
		json.Unmarshal(data, &lines)
		for _, line := range lines {
			if line != "" {
				existing[line] = true
			}
		}
	}

	written := 0
	skipped := 0

	// Append new keys
	f, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	for _, key := range result.Result {
		if existing[key.Name] {
			skipped++
			continue
		}
		f.WriteString(key.Name + "\n")
		written++
		existing[key.Name] = true
	}

	return written, skipped, nil
}

type SupabaseProvider struct {
	client *http.Client
}

func (p *SupabaseProvider) Name() string {
	return "supabase"
}

func (p *SupabaseProvider) Fetch(cfg *ProviderConfig, outputPath string) (int, int, error) {
	if cfg.Token == "" || cfg.URL == "" || cfg.Table == "" {
		return 0, 0, fmt.Errorf("supabase: missing token, url, or table")
	}

	url := fmt.Sprintf("%s/rest/v1/%s?order=created_at.desc&limit=100", cfg.URL, cfg.Table)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, 0, err
	}
	req.Header.Set("apikey", cfg.Token)
	req.Header.Set("Authorization", "Bearer "+cfg.Token)

	resp, err := p.client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, 0, fmt.Errorf("supabase API error: %s", body)
	}

	var result []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, 0, err
	}

	written := 0
	skipped := 0

	f, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	for _, row := range result {
		data, _ := json.Marshal(row)
		f.WriteString(string(data) + "\n")
		written++
	}

	return written, skipped, nil
}