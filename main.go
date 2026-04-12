package main

import (
	"fmt"
	"os"

	"github.com/formseal/formseal-sync/config"
	"github.com/formseal/formseal-sync/logic"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	cfg := config.LoadConfig()
	app := logic.New(&cfg)

	var providerCB *walk.ComboBox
	var cfTokenLE, cfNamespaceLE, sbURLLE, sbKeyLE, sbTableLE, outputFolderLE *walk.LineEdit
	var statTotal, statNew, statusBadge *walk.Label
	var logTV *walk.TextEdit
	var fetchBtn *walk.PushButton

	app.GetProvider = func() string { return providerCB.Text() }
	app.SetProvider = func(v string) { providerCB.SetText(v) }
	app.GetCfToken = func() string { return cfTokenLE.Text() }
	app.GetCfNamespace = func() string { return cfNamespaceLE.Text() }
	app.GetSbUrl = func() string { return sbURLLE.Text() }
	app.GetSbKey = func() string { return sbKeyLE.Text() }
	app.GetSbTable = func() string { return sbTableLE.Text() }
	app.GetOutput = func() string { return outputFolderLE.Text() }
	app.SetOutput = func(v string) { outputFolderLE.SetText(v) }
	app.SetStatTotal = func(v string) { statTotal.SetText(v) }
	app.SetStatNew = func(v string) { statNew.SetText(v) }
	app.SetStatus = func(v string) { statusBadge.SetText(v) }
	app.AppendLog = func(v string) { logTV.AppendText(v) }
	app.SetLog = func(v string) { logTV.SetText(v) }
	app.SetFetchEnabled = func(v bool) { fetchBtn.SetEnabled(v) }

	MainWindow(
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
							Label{Text: "—", AssignTo: &statTotal, Font: Font{PointSize: 24, Bold: true}},
							Label{Text: "ciphertexts stored", Font: Font{PointSize: 11}},
						},
					},
					Composite{
						Layout: VBox{Spacing: 2},
						Children: []Widget{
							Label{Text: "—", AssignTo: &statNew, Font: Font{PointSize: 24, Bold: true}},
							Label{Text: "new since last fetch", Font: Font{PointSize: 11}},
						},
					},
				},
			},
			Composite{
				Layout: HBox{Spacing: 10},
				Children: []Widget{
					Label{Text: "idle", AssignTo: &statusBadge, Font: Font{PointSize: 11}},
					Label{Text: "", Font: Font{PointSize: 11}},
				},
			},
			PushButton{
				Text:      "Fetch now",
				AssignTo:  &fetchBtn,
				OnClicked: func() { app.Fetch() },
			},
			GroupBox{
				Title:  "Provider",
				Layout: VBox{Spacing: 8},
				Children: []Widget{
					ComboBox{
						AssignTo: &providerCB,
						Model:    []string{"cloudflare", "supabase"},
					},
				},
			},
			GroupBox{
				Title:  "Cloudflare",
				Layout: VBox{Spacing: 8},
				Children: []Widget{
					Label{Text: "API Token"},
					LineEdit{AssignTo: &cfTokenLE, PasswordMode: true},
					Label{Text: "KV Namespace ID"},
					LineEdit{AssignTo: &cfNamespaceLE},
				},
			},
			GroupBox{
				Title:     "Supabase",
				Layout:    VBox{Spacing: 8},
				Visible:   false,
				Children: []Widget{
					Label{Text: "Project URL"},
					LineEdit{AssignTo: &sbURLLE},
					Label{Text: "Service Key"},
					LineEdit{AssignTo: &sbKeyLE, PasswordMode: true},
					Label{Text: "Table name"},
					LineEdit{AssignTo: &sbTableLE},
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
							LineEdit{AssignTo: &outputFolderLE, ReadOnly: true},
							PushButton{
								Text: "Browse",
								OnClicked: func() {
									dlg := new(walk.FileDialog)
									dlg.Title = "Select output folder"
									dlg.Flags = walk.FlagPullDown
									if ok, err := dlg.ShowFolderPicker(nil); err == nil && ok {
										outputFolderLE.SetText(dlg.FilePath)
									}
								},
							},
						},
					},
				},
			},
			PushButton{
				Text:      "Save",
				OnClicked: func() { app.Save() },
			},
			TextEdit{
				AssignTo: &logTV,
				ReadOnly: true,
				MinSize:  Size{0, 120},
			},
		},
	).Run()

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}