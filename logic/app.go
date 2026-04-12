package logic

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/formseal/formseal-sync/config"
	"github.com/formseal/formseal-sync/providers/cloudflare"
	"github.com/formseal/formseal-sync/providers/supabase"
)

type App struct {
	cfg *config.Config

	GetProvider    func() string
	SetProvider    func(string)
	GetCfToken     func() string
	GetCfNamespace func() string
	GetSbUrl       func() string
	GetSbKey       func() string
	GetSbTable     func() string
	GetOutput      func() string
	SetOutput      func(string)
	SetStatTotal   func(string)
	SetStatNew     func(string)
	SetStatus      func(string)
	AppendLog      func(string)
	SetLog         func(string)
	SetFetchEnabled func(bool)
}

func New(cfg *config.Config) *App {
	return &App{cfg: cfg}
}

func (a *App) Load() {
	if a.cfg.Provider == "supabase" {
		a.SetProvider("supabase")
	} else {
		a.SetProvider("cloudflare")
	}

	a.SetStatTotal("—")
	a.SetLog("Ready.")
	a.RefreshStats()
}

func (a *App) Save() {
	a.cfg.Provider = a.GetProvider()
	a.cfg.Cloudflare.Token = a.GetCfToken()
	a.cfg.Cloudflare.Namespace = a.GetCfNamespace()
	a.cfg.Supabase.URL = a.GetSbUrl()
	a.cfg.Supabase.Key = a.GetSbKey()
	a.cfg.Supabase.Table = a.GetSbTable()
	a.cfg.OutputFolder = a.GetOutput()

	if err := config.SaveConfig(a.cfg); err != nil {
		a.AppendLog("\nSave error: " + err.Error())
		return
	}
	a.AppendLog("\nSaved.")
}

func (a *App) Fetch() {
	a.SetFetchEnabled(false)
	a.SetLog("")
	a.AppendLog("Connecting...")
	a.SetStatus("fetching...")

	a.cfg.Provider = a.GetProvider()
	a.cfg.Cloudflare.Token = a.GetCfToken()
	a.cfg.Cloudflare.Namespace = a.GetCfNamespace()
	a.cfg.Supabase.URL = a.GetSbUrl()
	a.cfg.Supabase.Key = a.GetSbKey()
	a.cfg.Supabase.Table = a.GetSbTable()
	a.cfg.OutputFolder = a.GetOutput()

	if a.cfg.OutputFolder == "" {
		a.AppendLog("\nOutput folder not set.")
		a.SetStatus("error")
		a.SetFetchEnabled(true)
		return
	}

	outputPath := filepath.Join(a.cfg.OutputFolder, "ciphertexts.jsonl")

	var written, skipped int
	var err error

	switch a.cfg.Provider {
	case "cloudflare":
		if a.cfg.Cloudflare.Token == "" || a.cfg.Cloudflare.Namespace == "" {
			a.AppendLog("\nCloudflare token or namespace not set.")
			a.SetStatus("error")
			a.SetFetchEnabled(true)
			return
		}
		accountID, err := cloudflare.GetAccountID(a.cfg.Cloudflare.Token)
		if err != nil {
			a.AppendLog("\nAccount error: " + err.Error())
			a.SetStatus("error")
			a.SetFetchEnabled(true)
			return
		}
		written, skipped, err = cloudflare.FetchCiphertexts(
			a.cfg.Cloudflare.Token, a.cfg.Cloudflare.Namespace, accountID, outputPath)

	case "supabase":
		if a.cfg.Supabase.URL == "" || a.cfg.Supabase.Key == "" {
			a.AppendLog("\nSupabase URL or key not set.")
			a.SetStatus("error")
			a.SetFetchEnabled(true)
			return
		}
		table := a.cfg.Supabase.Table
		if table == "" {
			table = "ciphertexts"
		}
		written, skipped, err = supabase.FetchCiphertexts(
			a.cfg.Supabase.URL, a.cfg.Supabase.Key, table, outputPath)

	default:
		a.AppendLog("\nNo provider set.")
		a.SetStatus("error")
		a.SetFetchEnabled(true)
		return
	}

	if err != nil {
		a.AppendLog("\nError: " + err.Error())
		a.SetStatus("error")
	} else {
		a.AppendLog(fmt.Sprintf("\n%d new, %d duplicates", written, skipped))
		a.SetStatus("ok")
		a.SetStatNew(fmt.Sprintf("%d", written))
	}

	a.SetFetchEnabled(true)
	a.RefreshStats()
}

func (a *App) RefreshStats() {
	if a.cfg.OutputFolder == "" {
		a.SetStatTotal("—")
		return
	}
	path := filepath.Join(a.cfg.OutputFolder, "ciphertexts.jsonl")
	data, err := os.ReadFile(path)
	if err != nil {
		a.SetStatTotal("0")
		return
	}
	count := 0
	for _, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) != "" {
			count++
		}
	}
	a.SetStatTotal(fmt.Sprintf("%d", count))
}