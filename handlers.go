package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/grayguava/formseal-sync/providers"
	"github.com/grayguava/formseal-sync/providers/types"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(indexHTML)
}

func handleStatus(state *AppState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg := state.getConfig()
		result := state.getResult()

		msgCount := 0
		if cfg.OutputFolder != "" {
			outPath := filepath.Join(cfg.OutputFolder, "ciphertexts.jsonl")
			if data, err := os.ReadFile(outPath); err == nil {
				for _, b := range data {
					if b == '\n' {
						msgCount++
					}
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"result":   result,
			"msgCount": msgCount,
			"tokenSet": cfg.token() != "",
			"provider": cfg.Provider,
		})
	}
}

func handleConfig(state *AppState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(state.getConfig())

		case http.MethodPost:
			var cfg Config
			if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if err := saveConfig(&cfg); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			state.setConfig(&cfg)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"status": "saved"})

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func handleSync(state *AppState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		cfg := state.getConfig()
		go runSync(state, cfg)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "started"})
	}
}

func runSync(state *AppState, cfg *Config) {
	state.setResult(SyncResult{Done: false, RunAt: time.Now()})

	p := providers.Get(cfg.Provider)
	if p == nil {
		state.setResult(SyncResult{Done: true, Error: "unknown provider: " + cfg.Provider, RunAt: time.Now()})
		return
	}

	pCfg := &types.ProviderConfig{
		Token:     cfg.token(),
		Namespace: cfg.CloudflareNS,
		URL:       cfg.SupabaseURL,
		Table:     cfg.SupabaseTable,
	}

	if err := p.Validate(pCfg); err != nil {
		state.setResult(SyncResult{Done: true, Error: err.Error(), RunAt: time.Now()})
		return
	}

	if err := os.MkdirAll(cfg.OutputFolder, 0755); err != nil {
		state.setResult(SyncResult{Done: true, Error: "cannot create output folder: " + err.Error(), RunAt: time.Now()})
		return
	}

	outPath := filepath.Join(cfg.OutputFolder, "ciphertexts.jsonl")
	written, skipped, err := p.Fetch(pCfg, outPath)
	if err != nil {
		state.setResult(SyncResult{Done: true, Error: err.Error(), RunAt: time.Now()})
		return
	}

	state.setResult(SyncResult{
		Done:    true,
		Written: written,
		Skipped: skipped,
		RunAt:   time.Now(),
	})
}