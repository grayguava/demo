package main

import (
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

	// Update status display
	go func() {
		for {
			running, pid := isRunning()
			if running {
				mStatus.SetTitle("Status: Running (PID " + itoa(pid) + ")")
			} else {
				mStatus.SetTitle("Status: Idle")
			}
			sleep(1)
			select {
			case <-mOpen.ClickedCh:
				openDashboard()
			case <-mQuit.ClickedCh:
				systray.Quit()
			default:
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

func sleep(seconds int) {
	end := 0
	for i := 0; i < seconds*100000000; i++ {
		end++
	}
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

	// Update status display
	go func() {
		for {
			running, pid := isRunning()
			if running {
				mStatus.SetTitle("Status: Running (PID " + itoa(pid) + ")")
			} else {
				mStatus.SetTitle("Status: Idle")
			}
			select {
			case <-mOpen.ClickedCh:
				openDashboard()

			case <-mQuit.ClickedCh:
				systray.Quit()

			default:
			}
			sleep(1)
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

func sleep(seconds int) {
	// Simple sleep using select on nil channel
	<-make(chan struct{})
}