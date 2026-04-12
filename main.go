package main

import (
	"fmt"
	"os"

	"github.com/formseal/formseal-sync/config"
	"github.com/formseal/formseal-sync/logic"
	"github.com/fyne.io/fyne/v2"
	"github.com/fyne-io/fyne/v2/app"
	"github.com/fyne-io/fyne/v2/container"
	"github.com/fyne-io/fyne/v2/widget"
)

func main() {
	cfg := config.LoadConfig()
	a := app.New()
	a.SetTheme(nil)

	w := a.NewWindow("formseal-sync")
	w.Resize(fyne.NewSize(420, 580))

	appLogic := logic.New(&cfg)

	var providerSelect *widget.Select
	var cfTokenEntry, cfNamespaceEntry, sbUrlEntry, sbKeyEntry, sbTableEntry, outputEntry *widget.Entry
	var statTotalLabel, statNewLabel, statusLabel *widget.Label
	var logView *widget.Text

	providerSelect = widget.NewSelect([]string{"cloudflare", "supabase"}, func(s string) {
		cfTokenEntry.Hidden = s != "cloudflare"
		cfNamespaceEntry.Hidden = s != "cloudflare"
		sbUrlEntry.Hidden = s != "supabase"
		sbKeyEntry.Hidden = s != "supabase"
		sbTableEntry.Hidden = s != "supabase"
	})
	if cfg.Provider == "supabase" {
		providerSelect.SetSelected("supabase")
	}

	cfTokenEntry = widget.NewPasswordEntry()
	cfTokenEntry.PlaceHolder = "cfut_..."
	cfTokenEntry.Text = cfg.Cloudflare.Token

	cfNamespaceEntry = widget.NewEntry()
	cfNamespaceEntry.PlaceHolder = "0f2b2dcf..."
	cfNamespaceEntry.Text = cfg.Cloudflare.Namespace

	sbUrlEntry = widget.NewEntry()
	sbUrlEntry.PlaceHolder = "https://xxx.supabase.co"
	sbUrlEntry.Text = cfg.Supabase.URL
	sbUrlEntry.Hidden = true

	sbKeyEntry = widget.NewPasswordEntry()
	sbKeyEntry.PlaceHolder = "eyJ..."
	sbKeyEntry.Text = cfg.Supabase.Key
	sbKeyEntry.Hidden = true

	sbTableEntry = widget.NewEntry()
	sbTableEntry.PlaceHolder = "ciphertexts"
	sbTableEntry.Text = cfg.Supabase.Table
	sbTableEntry.Hidden = true

	outputEntry = widget.NewEntry()
	outputEntry.Text = cfg.OutputFolder
	outputEntry.ReadOnly = true

	browseBtn := widget.NewButton("Browse", func() {
		dialog := &fyne.FileDialog{Type: 0}
		dialog.SetOnFileSelected(func(f fyne.URI) {
			outputEntry.SetText(f.Path())
		})
	})

	saveBtn := widget.NewButton("Save", func() {
		appLogic.GetProvider = func() string { return providerSelect.Selected }
		appLogic.GetCfToken = func() string { return cfTokenEntry.Text }
		appLogic.GetCfNamespace = func() string { return cfNamespaceEntry.Text }
		appLogic.GetSbUrl = func() string { return sbUrlEntry.Text }
		appLogic.GetSbKey = func() string { return sbKeyEntry.Text }
		appLogic.GetSbTable = func() string { return sbTableEntry.Text }
		appLogic.GetOutput = func() string { return outputEntry.Text }
		appLogic.SetOutput = func(v string) { outputEntry.SetText(v) }
		appLogic.SetStatTotal = func(v string) { statTotalLabel.SetText(v) }
		appLogic.SetStatNew = func(v string) { statNewLabel.SetText(v) }
		appLogic.SetStatus = func(v string) { statusLabel.SetText(v) }
		appLogic.AppendLog = func(v string) { logView.SetText(logView.Text + v) }
		appLogic.SetLog = func(v string) { logView.SetText(v) }
		appLogic.SetFetchEnabled = func(v bool) { fetchBtn.Disable() }
		appLogic.Save()
	})

	fetchBtn := widget.NewButton("Fetch now", func() {
		appLogic.GetProvider = func() string { return providerSelect.Selected }
		appLogic.GetCfToken = func() string { return cfTokenEntry.Text }
		appLogic.GetCfNamespace = func() string { return cfNamespaceEntry.Text }
		appLogic.GetSbUrl = func() string { return sbUrlEntry.Text }
		appLogic.GetSbKey = func() string { return sbKeyEntry.Text }
		appLogic.GetSbTable = func() string { return sbTableEntry.Text }
		appLogic.GetOutput = func() string { return outputEntry.Text }
		appLogic.SetOutput = func(v string) { outputEntry.SetText(v) }
		appLogic.SetStatTotal = func(v string) { statTotalLabel.SetText(v) }
		appLogic.SetStatNew = func(v string) { statNewLabel.SetText(v) }
		appLogic.SetStatus = func(v string) { statusLabel.SetText(v) }
		appLogic.AppendLog = func(v string) { logView.SetText(logView.Text + v) }
		appLogic.SetLog = func(v string) { logView.SetText(v) }
		appLogic.SetFetchEnabled = func(v bool) {
			if v {
				fetchBtn.Enable()
			} else {
				fetchBtn.Disable()
			}
		}
		appLogic.Fetch()
	})

	statTotalLabel = widget.NewLabel("—")
	statNewLabel = widget.NewLabel("—")
	statusLabel = widget.NewLabel("idle")
	logView = widget.NewText()
	logView.ReadOnly = true

	appLogic.GetProvider = func() string { return providerSelect.Selected }
	appLogic.GetCfToken = func() string { return cfTokenEntry.Text }
	appLogic.GetCfNamespace = func() string { return cfNamespaceEntry.Text }
	appLogic.GetSbUrl = func() string { return sbUrlEntry.Text }
	appLogic.GetSbKey = func() string { return sbKeyEntry.Text }
	appLogic.GetSbTable = func() string { return sbTableEntry.Text }
	appLogic.GetOutput = func() string { return outputEntry.Text }
	appLogic.SetOutput = func(v string) { outputEntry.SetText(v) }
	appLogic.SetStatTotal = func(v string) { statTotalLabel.SetText(v) }
	appLogic.SetStatNew = func(v string) { statNewLabel.SetText(v) }
	appLogic.SetStatus = func(v string) { statusLabel.SetText(v) }
	appLogic.AppendLog = func(v string) { logView.SetText(logView.Text + v) }
	appLogic.SetLog = func(v string) { logView.SetText(v) }
	appLogic.SetFetchEnabled = func(v bool) {}
	appLogic.Load()

	content := container.NewVBox(
		widget.NewLabel("formseal-sync"),
		container.NewHBox(
			widget.NewLabel("ciphertexts:"),
			statTotalLabel,
		),
		container.NewHBox(
			widget.NewLabel("new:"),
			statNewLabel,
		),
		statusLabel,
		fetchBtn,
		widget.NewLabel("Provider"),
		providerSelect,
		widget.NewLabel("Cloudflare"),
		cfTokenEntry,
		cfNamespaceEntry,
		widget.NewLabel("Supabase"),
		sbUrlEntry,
		sbKeyEntry,
		sbTableEntry,
		widget.NewLabel("Output"),
		container.NewHBox(outputEntry, browseBtn),
		saveBtn,
		logView,
	)

	w.SetContent(content)
	w.ShowAndRun()
}