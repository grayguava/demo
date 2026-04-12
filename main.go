package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/formseal/formseal-sync/config"
	"github.com/formseal/formseal-sync/logic"
	"github.com/formseal/formseal-sync/providers/cloudflare"
	"github.com/formseal/formseal-sync/providers/supabase"
)

func main() {
	cfg := config.LoadConfig()

	if len(os.Args) < 2 {
		printHelp()
		return
	}

	switch os.Args[1] {
	case "config":
		showConfig(&cfg)
	case "set":
		setConfig(&cfg, os.Args[2:])
	case "fetch":
		fetch(&cfg)
	case "status":
		status(&cfg)
	case "help", "--help", "-h":
		printHelp()
	default:
		fmt.Println("Unknown command:", os.Args[1])
		printHelp()
	}
}

func printHelp() {
	fmt.Println(`formseal-sync - fetch encrypted form submissions

Usage:
  formseal-sync config           Show current configuration
  formseal-sync set <key> <val>  Set configuration value
  formseal-sync fetch            Fetch ciphertexts from provider
  formseal-sync status           Show stats
  formseal-sync help             Show this help

Config keys:
  provider         - cloudflare or supabase
  cf-token         - Cloudflare API token
  cf-namespace     - Cloudflare KV namespace ID
  sb-url           - Supabase project URL
  sb-key           - Supabase service key
  sb-table         - Supabase table name (default: ciphertexts)
  output           - Output folder path

Examples:
  formseal-sync set provider cloudflare
  formseal-sync set cf-token cfun_xxx
  formseal-sync set output C:\Users\you\data
  formseal-sync fetch`)
}

func showConfig(cfg *config.Config) {
	fmt.Println("Provider:", cfg.Provider)
	fmt.Println("Cloudflare Token:", cfg.Cloudflare.Token)
	fmt.Println("Cloudflare Namespace:", cfg.Cloudflare.Namespace)
	fmt.Println("Supabase URL:", cfg.Supabase.URL)
	fmt.Println("Supabase Key:", cfg.Supabase.Key)
	fmt.Println("Supabase Table:", cfg.Supabase.Table)
	fmt.Println("Output Folder:", cfg.OutputFolder)
}

func setConfig(cfg *config.Config, args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: formseal-sync set <key> <value>")
		return
	}

	key, value := args[0], args[1]

	switch key {
	case "provider":
		cfg.Provider = value
	case "cf-token":
		cfg.Cloudflare.Token = value
	case "cf-namespace":
		cfg.Cloudflare.Namespace = value
	case "sb-url":
		cfg.Supabase.URL = value
	case "sb-key":
		cfg.Supabase.Key = value
	case "sb-table":
		cfg.Supabase.Table = value
	case "output":
		cfg.OutputFolder = value
	default:
		fmt.Println("Unknown key:", key)
		return
	}

	if err := config.SaveConfig(cfg); err != nil {
		fmt.Println("Error saving config:", err)
		return
	}
	fmt.Println("Saved.")
}

func fetch(cfg *config.Config) {
	if cfg.OutputFolder == "" {
		fmt.Println("Error: output folder not set. Run: formseal-sync set output <path>")
		return
	}

	fmt.Println("Fetching from", cfg.Provider+"...")

	outputPath := cfg.OutputFolder + "/ciphertexts.jsonl"

	var written, skipped int
	var err error

	switch cfg.Provider {
	case "cloudflare":
		if cfg.Cloudflare.Token == "" || cfg.Cloudflare.Namespace == "" {
			fmt.Println("Error: Cloudflare token or namespace not set")
			return
		}
		accountID, err := cloudflare.GetAccountID(cfg.Cloudflare.Token)
		if err != nil {
			fmt.Println("Error getting account:", err)
			return
		}
		written, skipped, err = cloudflare.FetchCiphertexts(
			cfg.Cloudflare.Token, cfg.Cloudflare.Namespace, accountID, outputPath)
	case "supabase":
		if cfg.Supabase.URL == "" || cfg.Supabase.Key == "" {
			fmt.Println("Error: Supabase URL or key not set")
			return
		}
		table := cfg.Supabase.Table
		if table == "" {
			table = "ciphertexts"
		}
		written, skipped, err = supabase.FetchCiphertexts(
			cfg.Supabase.URL, cfg.Supabase.Key, table, outputPath)
	default:
		fmt.Println("Error: provider not set. Run: formseal-sync set provider <cloudflare|supabase>")
		return
	}

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Done! %d new, %d duplicates\n", written, skipped)
}

func status(cfg *config.Config) {
	if cfg.OutputFolder == "" {
		fmt.Println("Output folder not set")
		return
	}

	path := cfg.OutputFolder + "/ciphertexts.jsonl"
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("0 ciphertexts")
		return
	}
	defer file.Close()

	count := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() != "" {
			count++
		}
	}

	fmt.Printf("%d ciphertexts stored\n", count)
}