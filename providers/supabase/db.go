package supabase

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func FetchCiphertexts(url, key, table, outputPath string) (written, skipped int, err error) {
	if url == "" || key == "" {
		return 0, 0, fmt.Errorf("Supabase URL or key not set")
	}

	fetchURL := url + "/rest/v1/" + table + "?select=data"

	client := &http.Client{}
	req, err := http.NewRequest("GET", fetchURL, nil)
	if err != nil {
		return 0, 0, err
	}
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("apikey", key)

	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return 0, 0, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var rows []struct {
		Data string `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&rows); err != nil {
		return 0, 0, err
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

	for _, row := range rows {
		value := strings.TrimSpace(row.Data)
		if value == "" {
			continue
		}
		if seen[value] {
			skipped++
			continue
		}
		fmt.Fprintln(file, value)
		seen[value] = true
		written++
	}

	return written, skipped, nil
}

func readExisting(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var lines []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines, nil
}

func openAppend(path string) (*os.File, error) {
	dir := ""
	if i := strings.LastIndex(path, "/"); i > 0 {
		dir = path[:i]
	} else if i := strings.LastIndex(path, "\\"); i > 0 {
		dir = path[:i]
	}
	if dir != "" {
		os.MkdirAll(dir, 0755)
	}
	return os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
}