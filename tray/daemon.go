package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/getlantern/systray"
)

var (
	pidFile = filepath.Join(homeDir(), ".formsealdaemon", "sync.pid")
	logFile = filepath.Join(homeDir(), ".formsealdaemon", "sync.log")
)

func homeDir() string {
	h, _ := os.UserHomeDir()
	return h
}

func workerPath() string {
	exe, err := os.Executable()
	if err != nil {
		return "worker.py"
	}
	return filepath.Join(filepath.Dir(exe), "worker.py")
}

func isRunning() (bool, int) {
	data, err := os.ReadFile(pidFile)
	if err != nil {
		return false, 0
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return false, 0
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false, 0
	}
	err = proc.Signal(syscall.Signal(0))
	if err != nil {
		os.Remove(pidFile)
		return false, 0
	}
	return true, pid
}

func startDaemon() {
	if running, _ := isRunning(); running {
		return
	}

	python := findPython()
	cmd := exec.Command(python, workerPath())
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil

	if err := cmd.Start(); err != nil {
		showError("Failed to start sync: " + err.Error())
		log("ERROR: Failed to start sync: " + err.Error())
		return
	}

	os.MkdirAll(filepath.Dir(pidFile), 0755)
	os.WriteFile(pidFile, []byte(strconv.Itoa(cmd.Process.Pid)), 0644)
	log("INFO: Started sync daemon (PID " + strconv.Itoa(cmd.Process.Pid) + ")")
}

func stopDaemon() {
	running, pid := isRunning()
	if !running {
		return
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return
	}
	proc.Kill()
	os.Remove(pidFile)
	log("INFO: Stopped sync daemon (PID " + strconv.Itoa(pid) + ")")
}

func refreshStatus(mStatus *systray.MenuItem, mStart *systray.MenuItem, mStop *systray.MenuItem) {
	running, pid := isRunning()
	if running {
		mStatus.SetTitle(fmt.Sprintf("Running (PID %d)", pid))
		mStart.Disable()
		mStop.Enable()
	} else {
		mStatus.SetTitle("Not running")
		mStart.Enable()
		mStop.Disable()
	}
}

func findPython() string {
	candidates := []string{"pythonw", "python3", "python", "py"}

	for _, name := range candidates {
		p, err := exec.LookPath(name)
		if err != nil {
			continue
		}

		out, err := exec.Command(p, "--version").Output()
		if err != nil {
			continue
		}

		version := strings.TrimSpace(string(out))
		major, minor, ok := parsePythonVersion(version)
		if !ok {
			continue
		}

		if major < 3 || (major == 3 && minor < 8) {
			showError(fmt.Sprintf(
				"formseal-sync requires Python 3.8 or newer.\nFound: %s\n\nInstall from https://python.org",
				version,
			))
			os.Exit(1)
		}

		return p
	}

	showError("Python is required to run formseal-sync.\n\nInstall from https://python.org\nWindows: winget install Python.Python.3")
	os.Exit(1)
	return ""
}

func parsePythonVersion(output string) (major, minor int, ok bool) {
	parts := strings.Fields(output)
	if len(parts) < 2 {
		return 0, 0, false
	}
	nums := strings.Split(parts[1], ".")
	if len(nums) < 2 {
		return 0, 0, false
	}
	maj, err1 := strconv.Atoi(nums[0])
	min, err2 := strconv.Atoi(nums[1])
	if err1 != nil || err2 != nil {
		return 0, 0, false
	}
	return maj, min, true
}

func showError(msg string) {
	exec.Command("msg", "*", msg).Run()
}

func log(msg string) {
	os.MkdirAll(filepath.Dir(logFile), 0755)
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	f.WriteString(msg + "\n")
}