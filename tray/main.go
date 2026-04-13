package main

import (
	"log"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/getlantern/systray"

	"github.com/grayguava/formseal-sync/tray/daemon"
	"github.com/grayguava/formseal-sync/tray/dash"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icon)
	systray.SetTitle("formseal-sync")
	systray.SetTooltip("FormSeal-Sync")

	svc := daemon.NewSyncService()

	mStatus := systray.AddMenuItem("Status: Idle", "")
	mStatus.Disable()
	systray.AddSeparator()

	mStart := systray.AddMenuItem("Start Sync", "Start the sync service")
	mStop := systray.AddMenuItem("Stop Sync", "Stop the sync service")
	mStop.Disable()
	systray.AddSeparator()

	mDash := systray.AddMenuItem("Open Dashboard", "Open browser dashboard")
	systray.AddSeparator()

	mQuit := systray.AddMenuItem("Quit", "Quit formseal-sync")

	go func() {
		server := dash.New(svc)
		if err := server.Start(); err != nil {
			log.Printf("Dashboard server error: %v", err)
		}
	}()

	go func() {
		for {
			if svc.IsRunning() {
				mStatus.SetTitle("Status: Running")
				mStart.Disable()
				mStop.Enable()
			} else {
				mStatus.SetTitle("Status: Idle")
				mStart.Enable()
				mStop.Disable()
			}
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		for {
			select {
			case <-mStart.ClickedCh:
				svc.Start()
			case <-mStop.ClickedCh:
				svc.Stop()
			case <-mDash.ClickedCh:
				openDashboard()
			case <-mQuit.ClickedCh:
				svc.Stop()
				systray.Quit()
			}
		}
	}()
}

func onExit() {}

func openDashboard() {
	time.Sleep(500 * time.Millisecond)
	exec.Command("cmd", "/c", "start", "http://localhost:3847").Start()
}

var _ = filepath.Join // suppress unused