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

	mStatus := systray.AddMenuItem("Status: checking...", "")
	mStatus.Disable()
	systray.AddSeparator()

	mStart := systray.AddMenuItem("Start Sync", "Start background sync")
	mStop := systray.AddMenuItem("Stop Sync", "Stop background sync")
	mStop.Disable()
	systray.AddSeparator()

	mOpenCLI := systray.AddMenuItem("Open CLI Setup", "Run fsync setup quick")
	systray.AddSeparator()

	mQuit := systray.AddMenuItem("Quit", "Quit formseal-sync")

	// Initial status refresh
	go refreshStatus(mStatus, mStart, mStop)

	go func() {
		for {
			select {
			case <-mStart.ClickedCh:
				startDaemon()
				refreshStatus(mStatus, mStart, mStop)

			case <-mStop.ClickedCh:
				stopDaemon()
				refreshStatus(mStatus, mStart, mStop)

			case <-mOpenCLI.ClickedCh:
				openCLI()

			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()
}

func onExit() {}