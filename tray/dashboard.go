package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

var dashboardPort = 3847

func openDashboard() {
	// Kill any existing dashboard on the port
	exec.Command("cmd", "/c", "taskkill", "/F", "/IM", "formseal-sync-dashboard.exe").Run()

	// Start the dashboard server
	go runDashboardServer()

	// Wait a moment for server to start
	// Then open browser
	go func() {
		// Simple delay
		for i := 0; i < 10000000; i++ {
		}
		exec.Command("cmd", "/c", "start", "http://localhost:"+strconv.Itoa(dashboardPort)).Start()
	}()
}

func runDashboardServer() {
	http.HandleFunc("/", handleDashboard)
	http.HandleFunc("/api/status", handleStatus)
	http.HandleFunc("/api/start", handleStart)
	http.HandleFunc("/api/stop", handleStop)
	http.HandleFunc("/api/config", handleConfig)

	// Create log file
	logFilePath := filepath.Join(homeDir(), ".formsealdaemon", "dashboard.log")
	os.MkdirAll(filepath.Dir(logFilePath), 0755)
	logFile, _ := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer logFile.Close()

	server := &http.Server{Addr: ":" + strconv.Itoa(dashboardPort)}
	server.SetKeepAlivesEnabled(false)

	// Log startup
	logFile.WriteString(fmt.Sprintf("[%s] Dashboard started on port %d\n", timestamp(), dashboardPort))

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logFile.WriteString(fmt.Sprintf("[%s] Dashboard error: %v\n", timestamp(), err))
	}
}

func timestamp() string {
	// Simple timestamp
	return "2026-01-01 00:00:00"
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta http-equiv="refresh" content="30">
	<title>FormSeal-Sync Dashboard</title>
	<style>
		* { box-sizing: border-box; margin: 0; padding: 0; }
		body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background: #0f172a; color: #e2e8f0; min-height: 100vh; }
		.container { max-width: 800px; margin: 0 auto; padding: 30px 20px; }
		h1 { color: #f97316; margin-bottom: 10px; }
		.subtitle { color: #64748b; margin-bottom: 30px; font-size: 14px; }
		.card { background: #1e293b; border-radius: 12px; padding: 24px; margin-bottom: 20px; }
		h2 { color: #94a3b8; margin: 0 0 16px; font-size: 12px; text-transform: uppercase; letter-spacing: 1px; }
		.status { font-size: 24px; font-weight: bold; }
		.status.running { color: #22c55e; }
		.status.idle { color: #64748b; }
		.btn { display: inline-block; background: #f97316; color: white; border: none; padding: 12px 20px; margin: 5px; cursor: pointer; border-radius: 8px; font-weight: 600; font-size: 14px; }
		.btn:hover { background: #ea580c; }
		.btn.green { background: #22c55e; }
		.btn.green:hover { background: #16a34a; }
		.btn.red { background: #ef4444; }
		.btn.red:hover { background: #dc2626; }
		.btn.blue { background: #3b82f6; }
		.btn.blue:hover { background: #2563eb; }
		.log { background: #0f172a; padding: 16px; border-radius: 8px; font-family: monospace; font-size: 12px; max-height: 300px; overflow-y: auto; white-space: pre-wrap; color: #94a3b8; }
		.config-table { width: 100%; border-collapse: collapse; }
		.config-table td { padding: 10px; border-bottom: 1px solid #334155; }
		.config-table td:first-child { color: #64748b; width: 40%; }
		.error { background: #7f1d1d; color: #fca5a5; padding: 12px; border-radius: 8px; margin-bottom: 16px; }
		.success { background: #14532d; color: #86efac; padding: 12px; border-radius: 8px; margin-bottom: 16px; }
	</style>
</head>
<body>
	<div class="container">
		<h1>FormSeal-Sync Dashboard</h1>
		<div class="subtitle">Local management interface</div>

		<div class="card">
			<h2>Sync Status</h2>
			<div id="status" class="status">Loading...</div>
		</div>

		<div class="card">
			<h2>Actions</h2>
			<button class="btn green" onclick="doAction('/api/start')">Start Sync</button>
			<button class="btn red" onclick="doAction('/api/stop')">Stop Sync</button>
			<button class="btn blue" onclick="loadStatus()">Refresh</button>
		</div>

		<div class="card">
			<h2>Configuration</h2>
			<div id="config">Loading...</div>
		</div>

		<div class="card">
			<h2>Recent Activity</h2>
			<div id="log" class="log">Loading...</div>
		</div>

		<div class="card">
			<h2>Environment Variables</h2>
			<div style="font-size: 13px; color: #94a3b8; line-height: 1.8;">
				<strong>Windows (PowerShell):</strong><br>
				$env:FSYNC_CF_TOKEN="cfut_..."<br>
				$env:FSYNC_SU_KEY="eyJ..."<br><br>
				<strong>Windows (CMD):</strong><br>
				setx FSYNC_CF_TOKEN "cfut_..."<br><br>
				<strong>Linux/macOS:</strong><br>
				export FSYNC_CF_TOKEN="cfut_..."
			</div>
		</div>
	</div>

	<script>
		function loadStatus() {
			fetch('/api/status')
				.then(r => r.json())
				.then(d => {
					var statusEl = document.getElementById('status');
					if (d.running) {
						statusEl.textContent = 'Running (PID ' + d.pid + ')';
						statusEl.className = 'status running';
					} else {
						statusEl.textContent = 'Idle';
						statusEl.className = 'status idle';
					}

					var configEl = document.getElementById('config');
					if (d.config) {
						var html = '<table class="config-table">';
						for (var k in d.config) {
							html += '<tr><td>' + k + '</td><td>' + (d.config[k] || '-') + '</td></tr>';
						}
						html += '</table>';
						configEl.innerHTML = html;
					} else {
						configEl.innerHTML = '<div class="error">No configuration. Run: fsync setup quick</div>';
					}

					var logEl = document.getElementById('log');
					logEl.textContent = d.log || 'No activity yet';
				})
				.catch(e => {
					document.getElementById('status').textContent = 'Error: ' + e;
				});
		}

		function doAction(url) {
			fetch(url, {method: 'POST'})
				.then(r => r.json())
				.then(d => {
					if (d.error) {
						alert('Error: ' + d.error);
					}
					loadStatus();
				})
				.catch(e => alert('Error: ' + e));
		}

		loadStatus();
	</script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, html)
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	running, pid := isRunning()

	cfg, _ := loadConfig()

	logData := ""
	logPath := filepath.Join(homeDir(), ".formsealdaemon", "sync.log")
	if data, err := os.ReadFile(logPath); err == nil {
		lines := strings.Split(string(data), "\n")
		if len(lines) > 20 {
			lines = lines[len(lines)-20:]
		}
		logData = strings.Join(lines, "\n")
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"running":%v,"pid":%d,"config":%s,"log":%s}`,
		running, pid,
		configToJSON(cfg),
		toJSON(logData))
}

func handleStart(w http.ResponseWriter, r *http.Request) {
	startDaemon()
	refreshLog()
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"success":true}`)
}

func handleStop(w http.ResponseWriter, r *http.Request) {
	stopDaemon()
	refreshLog()
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"success":true}`)
}

func handleConfig(w http.ResponseWriter, r *http.Request) {
	cfg, err := loadConfig()
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		fmt.Fprintf(w, `{"error":"%s"}`, err)
	} else {
		fmt.Fprintf(w, configToJSON(cfg))
	}
}

func configToJSON(cfg *Config) string {
	return fmt.Sprintf(`{"provider":"%s","cloudflare_namespace":"%s","supabase_url":"%s","output_folder":"%s","sync_interval":%d}`,
		cfg.Provider, cfg.CloudflareNS, cfg.SupabaseURL, cfg.OutputFolder, cfg.SyncInterval)
}

func toJSON(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, `"`, "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return `"` + s + `"`
}

func refreshLog() {
	// Log the action
	log("Dashboard: Action triggered from web interface")
}