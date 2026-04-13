package main

import (
	"time"

	"github.com/getlantern/systray"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icon)
	systray.SetTitle("formseal-sync")
	systray.SetTooltip("FormSeal-Sync")

	mStatus := systray.AddMenuItem("Status: Idle", "")
	mStatus.Disable()
	systray.AddSeparator()

	mOpen := systray.AddMenuItem("Open Dashboard", "Open browser-based dashboard")
	systray.AddSeparator()

	mQuit := systray.AddMenuItem("Quit", "Quit formseal-sync")

	// Status update loop - runs every 2 seconds
	go func() {
		for {
			running, pid := isRunning()
			if running {
				mStatus.SetTitle("Status: Running (PID " + itoa(pid) + ")")
			} else {
				mStatus.SetTitle("Status: Idle")
			}
			time.Sleep(2 * time.Second)
		}
	}()

	// Menu event handler
	go func() {
		for {
			select {
			case <-mOpen.ClickedCh:
				openDashboard()
			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()
}

func onExit() {}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	return result
}