package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/formseal/formseal-sync/cmd/formseal-sync/providers/cloudflare"
	"github.com/formseal/formseal-sync/cmd/formseal-sync/providers/supabase"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type MainWindow struct {
	*walk.MainWindow
	providerCB     *walk.ComboBox
	cfTokenLE      *walk.LineEdit
	cfNamespaceLE  *walk.LineEdit
	sbURLLE        *walk.LineEdit
	sbKeyLE        *walk.LineEdit
	sbTableLE      *walk.LineEdit
	outputFolderLE *walk.LineEdit
	statTotal      *walk.Label
	statNew        *walk.Label
	statusBadge    *walk.Label
	statusTS       *walk.Label
	logTV          *walk.TextEdit
	fetchBtn       *walk.PushButton
}

func main() {
	cfg := LoadConfig()

	mw := &MainWindow{}
	if err := runMainWindow(mw, &cfg); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func runMainWindow(mw *MainWindow, cfg *Config) error {
	cfFields := new(walk.GroupBox)
	sbFields := new(walk.GroupBox)

	mw.MainWindow, _ = MainWindow(
		Title:    "formseal-sync",
		Width:    420,
		Height:   580,
		MinSize:  Size{420, 580},
		MaxSize:  Size{420, 580},
		Layout:   VBox{},
		Font:     Font{PointSize: 13},
		Children: []Widget{
			Composite{
				Layout: HBox{Spacing: 10},
				Children: []Widget{
					TextLabel{Text: "formseal-sync", Font: Font{PointSize: 14, Bold: true}},
				},
			},
			Composite{
				Layout: HBox{Spacing: 5},
				Children: []Widget{
					PushButton{Text: "Home", OnClicked: func() {}},
					PushButton{Text: "Settings", OnClicked: func() {}},
				},
			},
			GroupBox{
				Title:  "Stats",
				Layout: HBox{Spacing: 20},
				Children: []Widget{
					Composite{
						Layout: VBox{Spacing: 2},
						Children: []Widget{
							Label{Text: "—", AssignTo: &mw.statTotal, Font: Font{PointSize: 24, Bold: true}},
							Label{Text: "ciphertexts stored", Font: Font{PointSize: 11}},
						},
					},
					Composite{
						Layout: VBox{Spacing: 2},
						Children: []Widget{
							Label{Text: "—", AssignTo: &mw.statNew, Font: Font{PointSize: 24, Bold: true}},
							Label{Text: "new since last fetch", Font: Font{PointSize: 11}},
						},
					},
				},
			},
			Composite{
				Layout: HBox{Spacing: 10},
				Children: []Widget{
					Label{Text: "idle", AssignTo: &mw.statusBadge, Font: Font{PointSize: 11}},
					Label{Text: "", AssignTo: &mw.statusTS, Font: Font{PointSize: 11}},
				},
			},
			PushButton{
				Text:      "Fetch now",
				AssignTo:  &mw.fetchBtn,
				OnClicked: func() { mw.runFetch(cfg) },
			},
			GroupBox{
				Title:  "Provider",
				Layout: VBox{Spacing: 8},
				Children: []Widget{
					ComboBox{
						AssignTo: &mw.providerCB,
						Model:    []string{"cloudflare", "supabase"},
						OnCurrentIndexChanged: func() {
							idx := mw.providerCB.CurrentIndex()
							cfFields.SetVisible(idx == 0)
							sbFields.SetVisible(idx == 1)
						},
					},
				},
			},
			GroupBox{
				AssignTo:  &cfFields,
				Title:     "Cloudflare",
				Layout:    VBox{Spacing: 8},
				Children: []Widget{
					Label{Text: "API Token"},
					LineEdit{AssignTo: &mw.cfTokenLE, Text: cfg.Cloudflare.Token, PasswordMode: true},
					Label{Text: "KV Namespace ID"},
					LineEdit{AssignTo: &mw.cfNamespaceLE, Text: cfg.Cloudflare.Namespace},
				},
			},
			GroupBox{
				AssignTo:  &sbFields,
				Title:     "Supabase",
				Layout:    VBox{Spacing: 8},
				Visible:   false,
				Children: []Widget{
					Label{Text: "Project URL"},
					LineEdit{AssignTo: &mw.sbURLLE, Text: cfg.Supabase.URL},
					Label{Text: "Service Key"},
					LineEdit{AssignTo: &mw.sbKeyLE, Text: cfg.Supabase.Key, PasswordMode: true},
					Label{Text: "Table name"},
					LineEdit{AssignTo: &mw.sbTableLE, Text: cfg.Supabase.Table},
				},
			},
			GroupBox{
				Title:  "Output",
				Layout: VBox{Spacing: 8},
				Children: []Widget{
					Label{Text: "Output folder"},
					Composite{
						Layout: HBox{Spacing: 8},
						Children: []Widget{
							LineEdit{AssignTo: &mw.outputFolderLE, Text: cfg.OutputFolder, ReadOnly: true},
							PushButton{Text: "Browse", OnClicked: func() { mw.browseFolder(cfg) }},
						},
					},
				},
			},
			PushButton{
				Text:      "Save",
				OnClicked: func() { mw.saveConfig(cfg) },
			},
			TextEdit{
				AssignTo: &mw.logTV,
				ReadOnly: true,
				MinSize:  Size{0, 120},
			},
		},
	).Create()

	if mw.providerCB != nil {
		mw.providerCB.SetCurrentIndex(0)
		if cfg.Provider == "supabase" {
			mw.providerCB.SetCurrentIndex(1)
		}
	}

	mw.logTV.SetText("Ready.")
	mw.refreshStats(cfg)

	return nil
}

func (mw *MainWindow) browseFolder(cfg *Config) {
	dlg := new(walk.FileDialog)
	dlg.Title = "Select output folder"
	dlg.Flags = walk.FlagPullDown
	if ok, err := dlg.ShowFolderPicker(mw); err == nil && ok {
		mw.outputFolderLE.SetText(dlg.FilePath)
	}
}

func (mw *MainWindow) saveConfig(cfg *Config) {
	cfg.Provider = mw.providerCB.Text()
	cfg.Cloudflare.Token = mw.cfTokenLE.Text()
	cfg.Cloudflare.Namespace = mw.cfNamespaceLE.Text()
	cfg.Supabase.URL = mw.sbURLLE.Text()
	cfg.Supabase.Key = mw.sbKeyLE.Text()
	cfg.Supabase.Table = mw.sbTableLE.Text()
	cfg.OutputFolder = mw.outputFolderLE.Text()

	if err := SaveConfig(cfg); err != nil {
		mw.logTV.AppendText("\nSave error: " + err.Error())
		return
	}
	mw.logTV.AppendText("\nSaved.")
}

func (mw *MainWindow) runFetch(cfg *Config) {
	mw.fetchBtn.SetEnabled(false)
	mw.logTV.SetText("")
	mw.logTV.AppendText("Connecting...")
	mw.statusBadge.SetText("fetching...")

	cfg.Provider = mw.providerCB.Text()
	cfg.Cloudflare.Token = mw.cfTokenLE.Text()
	cfg.Cloudflare.Namespace = mw.cfNamespaceLE.Text()
	cfg.Supabase.URL = mw.sbURLLE.Text()
	cfg.Supabase.Key = mw.sbKeyLE.Text()
	cfg.Supabase.Table = mw.sbTableLE.Text()
	cfg.OutputFolder = mw.outputFolderLE.Text()

	if cfg.OutputFolder == "" {
		mw.logTV.AppendText("\nOutput folder not set.")
		mw.statusBadge.SetText("error")
		mw.fetchBtn.SetEnabled(true)
		return
	}

	outputPath := filepath.Join(cfg.OutputFolder, "ciphertexts.jsonl")

	var written, skipped int
	var err error

	switch cfg.Provider {
	case "cloudflare":
		if cfg.Cloudflare.Token == "" || cfg.Cloudflare.Namespace == "" {
			mw.logTV.AppendText("\nCloudflare token or namespace not set.")
			mw.statusBadge.SetText("error")
			mw.fetchBtn.SetEnabled(true)
			return
		}
		accountID, err := cloudflare.GetAccountID(cfg.Cloudflare.Token)
		if err != nil {
			mw.logTV.AppendText("\nAccount error: " + err.Error())
			mw.statusBadge.SetText("error")
			mw.fetchBtn.SetEnabled(true)
			return
		}
		written, skipped, err = cloudflare.FetchCiphertexts(
			cfg.Cloudflare.Token, cfg.Cloudflare.Namespace, accountID, outputPath)
	case "supabase":
		if cfg.Supabase.URL == "" || cfg.Supabase.Key == "" {
			mw.logTV.AppendText("\nSupabase URL or key not set.")
			mw.statusBadge.SetText("error")
			mw.fetchBtn.SetEnabled(true)
			return
		}
		table := cfg.Supabase.Table
		if table == "" {
			table = "ciphertexts"
		}
		written, skipped, err = supabase.FetchCiphertexts(
			cfg.Supabase.URL, cfg.Supabase.Key, table, outputPath)
	default:
		mw.logTV.AppendText("\nNo provider set.")
		mw.statusBadge.SetText("error")
		mw.fetchBtn.SetEnabled(true)
		return
	}

	if err != nil {
		mw.logTV.AppendText("\nError: " + err.Error())
		mw.statusBadge.SetText("error")
	} else {
		mw.logTV.AppendText(fmt.Sprintf("\n%d new, %d duplicates", written, skipped))
		mw.statusBadge.SetText("ok")
		mw.statNew.SetText(fmt.Sprintf("%d", written))
	}

	mw.fetchBtn.SetEnabled(true)
	mw.refreshStats(cfg)
}

func (mw *MainWindow) refreshStats(cfg *Config) {
	if cfg.OutputFolder == "" {
		mw.statTotal.SetText("—")
		return
	}
	path := filepath.Join(cfg.OutputFolder, "ciphertexts.jsonl")
	data, err := os.ReadFile(path)
	if err != nil {
		mw.statTotal.SetText("0")
		return
	}
	count := 0
	for _, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) != "" {
			count++
		}
	}
	mw.statTotal.SetText(fmt.Sprintf("%d", count))
}