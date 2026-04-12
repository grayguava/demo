package cloudflare

import (
	"os"
	"strings"
)

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