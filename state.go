package main

import (
	"sync"
	"time"
)

type SyncResult struct {
	Done    bool      `json:"done"`
	Error   string    `json:"error,omitempty"`
	Written int       `json:"written"`
	Skipped int       `json:"skipped"`
	RunAt   time.Time `json:"run_at"`
}

type AppState struct {
	mu     sync.RWMutex
	cfg    *Config
	result SyncResult
}

func newAppState(cfg *Config) *AppState {
	return &AppState{cfg: cfg}
}

func (s *AppState) setResult(r SyncResult) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.result = r
}

func (s *AppState) getResult() SyncResult {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.result
}

func (s *AppState) getConfig() *Config {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cfg
}

func (s *AppState) setConfig(cfg *Config) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cfg = cfg
}