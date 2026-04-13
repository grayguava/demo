package dash

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/grayguava/formseal-sync/tray/daemon"
	ui "github.com/grayguava/formseal-sync/tray/dash/ui"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, ui.IndexHTML)
}

func handleStatus(svc interface {
	IsRunning() bool
	LastSyncTime() time.Time
	GetConfig() interface{}
}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg, _ := daemon.LoadConfig()
		logPath := filepath.Join(os.Getenv("USERPROFILE"), ".formsealdaemon", "sync.log")

		var logs string
		if data, err := os.ReadFile(logPath); err == nil {
			lines := strings.Split(string(data), "\n")
			if len(lines) > 20 {
				lines = lines[len(lines)-20:]
			}
			logs = strings.Join(lines, "\n")
		}

		msgCount := 0
		if cfg != nil && cfg.OutputFolder != "" {
			outputPath := filepath.Join(cfg.OutputFolder, "ciphertexts.jsonl")
			if data, err := os.ReadFile(outputPath); err == nil {
				msgCount = len(strings.Split(string(data), "\n")) - 1
				if msgCount < 0 {
					msgCount = 0
				}
			}
		}

		var lastSync string
		if svc.LastSyncTime().IsZero() {
			lastSync = "Never"
		} else {
			lastSync = svc.LastSyncTime().Format("2006-01-02 15:04:05")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"running":   svc.IsRunning(),
			"msgCount":   msgCount,
			"lastSync":   lastSync,
			"logs":       logs,
			"config":     cfg,
		})
	}
}

func handleStart(svc interface {
	Start()
}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc.Start()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "started"})
	}
}

func handleStop(svc interface {
	Stop()
}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc.Stop()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "stopped"})
	}
}

func handleSave(svc interface {
	SetConfig(*daemon.Config)
	GetConfig() interface{}
}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var cfg daemon.Config
		if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		if err := daemon.SaveConfig(&cfg); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		svc.SetConfig(&cfg)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "saved"})
	}
}

var _ = fmt.Sprint // suppress unused