package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"time"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	state := newAppState(cfg)

	ln, err := startServer(state)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
	defer ln.Close()

	addr := fmt.Sprintf("http://localhost:%d", defaultPort)
	fmt.Printf("formseal-sync running at %s\n", addr)

	// Run sync immediately on launch
	go runSync(state, cfg)

	// Open browser
	time.Sleep(300 * time.Millisecond)
	openBrowser(addr)

	// Block forever — user closes the window/process to quit
	select {}
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	cmd.Start()
}