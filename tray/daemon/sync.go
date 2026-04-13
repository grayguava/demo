package daemon

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/grayguava/formseal-sync/tray/providers"
)

var (
	syncLogFile = filepath.Join(userHome(), ".formsealdaemon", "sync.syncLog")
)

type SyncService struct {
	mu       sync.RWMutex
	running  bool
	lastSync time.Time
	cfg      *Config
	stopCh   chan struct{}
}

func NewSyncService() *SyncService {
	cfg, _ := LoadConfig()
	return &SyncService{
		cfg: cfg,
	}
}

func (s *SyncService) SetConfig(cfg *Config) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cfg = cfg
}

func (s *SyncService) GetConfig() *Config {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cfg
}

func (s *SyncService) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.stopCh = make(chan struct{})
	s.mu.Unlock()

	go s.runLoop()
	syncLog("INFO: Sync started")
}

func (s *SyncService) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.running {
		return
	}
	s.running = false
	close(s.stopCh)
	syncLog("INFO: Sync stopped")
}

func (s *SyncService) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

func (s *SyncService) runLoop() {
	ticker := time.NewTicker(time.Duration(s.cfg.SyncInterval) * time.Minute)
	defer ticker.Stop()

	s.syncOnce()

	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.syncOnce()
		}
	}
}

func (s *SyncService) syncOnce() {
	s.mu.RLock()
	cfg := s.cfg
	s.mu.RUnlock()

	if cfg.Provider == "" {
		syncLog("WARN: No provider configured")
		return
	}

	provider := providers.Get(cfg.Provider)
	if provider == nil {
		syncLog("ERROR: Unknown provider: " + cfg.Provider)
		return
	}

	// Create provider config
	pCfg := &providers.ProviderConfig{
		Token:     cfg.APIToken,
		Namespace: cfg.CloudflareNS,
		URL:       cfg.SupabaseURL,
		Table:     cfg.SupabaseTable,
	}

	// Ensure output folder exists
	os.MkdirAll(cfg.OutputFolder, 0755)
	outputPath := filepath.Join(cfg.OutputFolder, "ciphertexts.jsonl")

	written, skipped, err := provider.Fetch(pCfg, outputPath)
	if err != nil {
		syncLog("ERROR: Fetch failed: " + err.Error())
		return
	}

	s.mu.Lock()
	s.lastSync = time.Now()
	s.mu.Unlock()

	syncLog(fmt.Sprintf("INFO: Synced %d new, %d skipped", written, skipped))
}

func (s *SyncService) LastSyncTime() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastSync
}

func syncLog(msg string) {
	os.MkdirAll(filepath.Dir(syncLogFile), 0755)
	f, err := os.OpenFile(syncLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	f.WriteString(fmt.Sprintf("[%s] %s\n", timestamp, msg))
}