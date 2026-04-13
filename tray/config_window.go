package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func openConfigWindow() {
	cfg, err := loadConfig()
	if err != nil {
		showError("Failed to load config: " + err.Error())
		return
	}

	interval := cfg.SyncInterval
	if interval == 0 {
		interval = 15
	}

	html := `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>FormSeal-Sync Settings</title>
	<style>
		* { box-sizing: border-box; margin: 0; padding: 0; }
		body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background: #1e293b; color: #e2e8f0; min-height: 100vh; }
		.container { max-width: 500px; margin: 0 auto; padding: 30px 20px; }
		h1 { color: #f97316; margin-bottom: 25px; font-size: 24px; }
		h2 { color: #94a3b8; margin: 20px 0 10px; font-size: 14px; text-transform: uppercase; letter-spacing: 1px; }
		label { display: block; margin: 12px 0 6px; font-weight: 500; font-size: 14px; }
		input, select { width: 100%; padding: 12px; margin-top: 4px; background: #334155; border: 1px solid #475569; color: #e2e8f0; border-radius: 6px; font-size: 14px; }
		input:focus, select:focus { outline: none; border-color: #f97316; }
		button { background: #f97316; color: white; border: none; padding: 14px 24px; margin: 10px 10px 0 0; cursor: pointer; border-radius: 6px; font-weight: bold; font-size: 14px; }
		button:hover { background: #ea580c; }
		button.secondary { background: #475569; }
		button.secondary:hover { background: #64748b; }
		button.green { background: #22c55e; }
		button.green:hover { background: #16a34a; }
		.note { font-size: 12px; color: #94a3b8; margin-top: 6px; }
		.error { color: #ef4444; font-size: 13px; margin-top: 12px; padding: 10px; background: rgba(239,68,68,0.1); border-radius: 6px; }
		.success { color: #22c55e; font-size: 13px; margin-top: 12px; padding: 10px; background: rgba(34,197,94,0.1); border-radius: 6px; }
		.section { background: #334155; padding: 20px; border-radius: 8px; margin-bottom: 20px; }
		.status { background: #0f172a; padding: 15px; border-radius: 6px; font-family: monospace; font-size: 12px; white-space: pre-wrap; max-height: 200px; overflow-y: auto; }
	</style>
</head>
<body>
	<div class="container">
		<h1>FormSeal-Sync Settings</h1>
		
		<div class="section">
			<h2>Quick Actions</h2>
			<button class="green" onclick="runCmd('start')">Start Sync</button>
			<button class="secondary" onclick="runCmd('stop')">Stop Sync</button>
			<button class="secondary" onclick="runCmd('status')">Check Status</button>
		</div>

		<div class="section">
			<h2>Configuration</h2>
			<button onclick="runCmd('setup')">Run Setup Wizard</button>
			<button class="secondary" onclick="runCmd('help')">View Help</button>
		</div>

		<div class="section">
			<h2>Current Config</h2>
			<div class="status" id="config-display">` + formatConfig(cfg) + `</div>
		</div>

		<div id="message"></div>

		<div style="margin-top: 30px; padding-top: 20px; border-top: 1px solid #475569;">
			<div class="note">
				<strong>Environment Variables:</strong><br>
				FSYNC_CF_TOKEN - Cloudflare API token<br>
				FSYNC_SU_KEY - Supabase service key<br>
				<br>
				<strong>Config Location:</strong><br>
				` + configFile + `
			</div>
		</div>
	</div>

	<script>
		function showMessage(msg, isError) {
			var el = document.getElementById('message');
			el.className = isError ? 'error' : 'success';
			el.textContent = msg;
			setTimeout(function() { el.textContent = ''; }, 5000);
		}

		function runCmd(action) {
			var cmd, args;
			
			switch(action) {
				case 'start':
					cmd = 'fsync';
					args = ['sync', 'start'];
					break;
				case 'stop':
					cmd = 'fsync';
					args = ['sync', 'stop'];
					break;
				case 'status':
					cmd = 'fsync';
					args = ['sync', 'status'];
					break;
				case 'setup':
					cmd = 'fsync';
					args = ['setup', 'quick'];
					break;
				case 'help':
					cmd = 'fsync';
					args = ['--help'];
					break;
			}

			showMessage('Running: ' + cmd + ' ' + args.join(' '), false);

			// Open a command prompt with the command pre-filled
			var command = cmd + ' ' + args.join(' ');
			window.location.href = 'cmd:' + command;
		}
	</script>
</body>
</html>`

	tmpDir := filepath.Join(os.TempDir(), "formseal-sync")
	os.MkdirAll(tmpDir, 0755)
	htmlPath := filepath.Join(tmpDir, "config.html")
	os.WriteFile(htmlPath, []byte(html), 0644)

	exec.Command("cmd", "/c", "start", htmlPath).Start()
}

func formatConfig(cfg *Config) string {
	var lines []string
	lines = append(lines, "Provider: "+cfg.Provider)
	if cfg.CloudflareNS != "" {
		lines = append(lines, "Cloudflare Namespace: "+cfg.CloudflareNS)
	}
	if cfg.SupabaseURL != "" {
		lines = append(lines, "Supabase URL: "+cfg.SupabaseURL)
	}
	lines = append(lines, "Output Folder: "+cfg.OutputFolder)
	lines = append(lines, "Sync Interval: "+fmt.Sprintf("%d", cfg.SyncInterval)+" min")
	return strings.Join(lines, "\n")
}

func selected(a, b string) string {
	if a == b {
		return " selected"
	}
	return ""
}

func visible(a, b string) string {
	if a == b {
		return ""
	}
	return "hidden"
}