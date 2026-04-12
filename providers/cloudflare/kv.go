package cloudflare

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type UserResponse struct {
	Success bool   `json:"success"`
	Result  Result `json:"result"`
	Errors  []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type Result struct {
	Accounts []Account `json:"accounts"`
}

type Account struct {
	ID string `json:"id"`
}

func GetAccountID(token string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.cloudflare.com/client/v4/user", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	var data UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	if !data.Success {
		if len(data.Errors) > 0 {
			return "", fmt.Errorf(data.Errors[0].Message)
		}
		return "", fmt.Errorf("unknown error")
	}

	if len(data.Result.Accounts) == 0 {
		return "", fmt.Errorf("no accounts found")
	}

	return data.Result.Accounts[0].ID, nil
}

func FetchCiphertexts(token, namespaceID, accountID, outputPath string) (written, skipped int, err error) {
	client := &http.Client{}
	baseURL := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/storage/kv/namespaces/%s",
		accountID, namespaceID)

	var allKeys []string
	cursor := ""

	for {
		url := baseURL + "/keys"
		if cursor != "" {
			url += "?cursor=" + cursor
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return 0, 0, err
		}
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := client.Do(req)
		if err != nil {
			return 0, 0, err
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode != 200 {
			return 0, 0, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
		}

		var data struct {
			Success    bool   `json:"success"`
			Result     []struct {
				Name string `json:"name"`
			} `json:"result"`
			ResultInfo struct {
				Cursor string `json:"cursor"`
			} `json:"result_info"`
		}
		if err := json.Unmarshal(body, &data); err != nil {
			return 0, 0, err
		}

		for _, k := range data.Result {
			allKeys = append(allKeys, k.Name)
		}

		cursor = data.ResultInfo.Cursor
		if cursor == "" {
			break
		}
	}

	if len(allKeys) == 0 {
		return 0, 0, nil
	}

	seen := make(map[string]bool)
	if existing, err := readExisting(outputPath); err == nil {
		for _, v := range existing {
			seen[v] = true
		}
	}

	file, err := openAppend(outputPath)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	for _, key := range allKeys {
		url := baseURL + "/values/" + strings.ReplaceAll(key, " ", "%20")

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			continue
		}
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := client.Do(req)
		if err != nil {
			continue
		}

		value, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode != 200 {
			continue
		}

		valueStr := strings.TrimSpace(string(value))
		if valueStr == "" {
			continue
		}

		if seen[valueStr] {
			skipped++
			continue
		}

		fmt.Fprintln(file, valueStr)
		seen[valueStr] = true
		written++
	}

	return written, skipped, nil
}