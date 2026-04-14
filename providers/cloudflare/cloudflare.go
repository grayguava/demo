package cloudflare

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/grayguava/formseal-sync/providers/types"
)

// Provider implements types.Provider for Cloudflare KV.
type Provider struct{}

const apiBase = "https://api.cloudflare.com/client/v4"

var httpClient = &http.Client{Timeout: 30 * time.Second}

func (p *Provider) Validate(cfg *types.ProviderConfig) error {
	if cfg.Token == "" {
		return fmt.Errorf("token is required (set the env var named in token_env)")
	}
	if cfg.Namespace == "" {
		return fmt.Errorf("cloudflare_namespace is required")
	}
	return nil
}

func (p *Provider) Fetch(cfg *types.ProviderConfig, outputPath string) (int, int, error) {
	if err := p.Validate(cfg); err != nil {
		return 0, 0, err
	}

	accountID, err := getAccountID(cfg.Token)
	if err != nil {
		return 0, 0, fmt.Errorf("auth failed: %w", err)
	}

	base := fmt.Sprintf("%s/accounts/%s/storage/kv/namespaces/%s", apiBase, accountID, cfg.Namespace)

	keys, err := listKeys(base, cfg.Token)
	if err != nil {
		return 0, 0, fmt.Errorf("listing keys: %w", err)
	}

	seen, err := loadSeen(outputPath)
	if err != nil {
		return 0, 0, fmt.Errorf("reading existing data: %w", err)
	}

	f, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, 0, fmt.Errorf("opening output file: %w", err)
	}
	defer f.Close()

	written, skipped := 0, 0

	for _, key := range keys {
		value, err := getValue(base, cfg.Token, key)
		if err != nil {
			return written, skipped, fmt.Errorf("fetching key %q: %w", key, err)
		}
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if seen[value] {
			skipped++
			continue
		}
		if _, err := fmt.Fprintln(f, value); err != nil {
			return written, skipped, fmt.Errorf("writing output: %w", err)
		}
		seen[value] = true
		written++
	}

	return written, skipped, nil
}

func getAccountID(token string) (string, error) {
	var result struct {
		Success bool `json:"success"`
		Result  []struct {
			ID string `json:"id"`
		} `json:"result"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	if err := apiGet(token, apiBase+"/accounts", &result); err != nil {
		return "", err
	}
	if !result.Success {
		if len(result.Errors) > 0 {
			return "", fmt.Errorf("%s", result.Errors[0].Message)
		}
		return "", fmt.Errorf("API returned failure")
	}
	if len(result.Result) == 0 {
		return "", fmt.Errorf("no accounts found for this token")
	}
	return result.Result[0].ID, nil
}

func listKeys(base, token string) ([]string, error) {
	var all []string
	cursor := ""

	for {
		endpoint := base + "/keys"
		if cursor != "" {
			endpoint += "?cursor=" + url.QueryEscape(cursor)
		}

		var result struct {
			Success bool `json:"success"`
			Result  []struct {
				Name string `json:"name"`
			} `json:"result"`
			ResultInfo struct {
				Cursor string `json:"cursor"`
			} `json:"result_info"`
			Errors []struct {
				Message string `json:"message"`
			} `json:"errors"`
		}

		if err := apiGet(token, endpoint, &result); err != nil {
			return nil, err
		}
		if !result.Success {
			if len(result.Errors) > 0 {
				return nil, fmt.Errorf("%s", result.Errors[0].Message)
			}
			return nil, fmt.Errorf("API returned failure")
		}

		for _, k := range result.Result {
			all = append(all, k.Name)
		}

		cursor = result.ResultInfo.Cursor
		if cursor == "" {
			break
		}
	}

	return all, nil
}

func getValue(base, token, key string) (string, error) {
	endpoint := base + "/values/" + url.PathEscape(key)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func apiGet(token, endpoint string, out interface{}) error {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(out)
}

func loadSeen(path string) (map[string]bool, error) {
	seen := make(map[string]bool)
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return seen, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			seen[line] = true
		}
	}
	return seen, scanner.Err()
}