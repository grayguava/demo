package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	<title>FormSeal-Sync Configuration</title>
	<style>
		body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; max-width: 500px; margin: 40px auto; padding: 20px; background: #1e293b; color: #e2e8f0; }
		h1 { color: #f97316; margin-bottom: 20px; }
		label { display: block; margin: 15px 0 5px; font-weight: 500; }
		input, select { width: 100%; padding: 10px; margin-top: 5px; background: #334155; border: 1px solid #475569; color: #e2e8f0; border-radius: 6px; box-sizing: border-box; }
		button { background: #f97316; color: white; border: none; padding: 12px 24px; margin-top: 20px; margin-right: 10px; cursor: pointer; border-radius: 6px; font-weight: bold; }
		button:hover { background: #ea580c; }
		button.secondary { background: #475569; }
		button.secondary:hover { background: #64748b; }
		.note { font-size: 12px; color: #94a3b8; margin-top: 5px; }
		.hidden { display: none; }
		.error { color: #ef4444; font-size: 12px; margin-top: 5px; }
	</style>
</head>
<body>
	<h1>FormSeal-Sync Configuration</h1>
	
	<label>Provider</label>
	<select id="provider" onchange="toggleFields()">
		<option value="">Select provider...</option>
		<option value="cloudflare"` + selected(cfg.Provider, "cloudflare") + `>Cloudflare KV</option>
		<option value="supabase"` + selected(cfg.Provider, "supabase") + `>Supabase</option>
	</select>

	<div id="cloudflare-fields" class="` + visible(cfg.Provider, "cloudflare") + `">
		<label>Namespace ID</label>
		<input type="text" id="cloudflare_namespace" value="` + cfg.CloudflareNS + `" placeholder="KV Namespace ID">
		<div class="note">Set FSYNC_CF_TOKEN environment variable</div>
	</div>

	<div id="supabase-fields" class="` + visible(cfg.Provider, "supabase") + `">
		<label>Project URL</label>
		<input type="text" id="supabase_url" value="` + cfg.SupabaseURL + `" placeholder="https://xxx.supabase.co">
		<label>Table Name</label>
		<input type="text" id="supabase_table" value="` + cfg.SupabaseTable + `" placeholder="ciphertexts">
		<div class="note">Set FSYNC_SU_KEY environment variable</div>
	</div>

	<div class="error" id="error"></div>

	<label>Output Folder</label>
	<input type="text" id="output_folder" value="` + cfg.OutputFolder + `" placeholder="D:/Documents/FormData">

	<label>Sync Interval (minutes)</label>
	<input type="number" id="sync_interval" value="` + fmt.Sprintf("%d", interval) + `" min="1" max="1440">

	<button onclick="save()">Save</button>
	<button class="secondary" onclick="openCLI()">Open CLI Setup</button>

	<script>
		function toggleFields() {
			var p = document.getElementById('provider').value;
			document.getElementById('cloudflare-fields').style.display = p === 'cloudflare' ? 'block' : 'none';
			document.getElementById('supabase-fields').style.display = p === 'supabase' ? 'block' : 'none';
		}
		toggleFields();

		function selected(a, b) { return a === b ? ' selected' : ''; }
		function visible(a, b) { return a === b ? '' : 'hidden'; }
		function itoa(n) { return '' + n; }

		function save() {
			var data = {
				provider: document.getElementById('provider').value,
				"cloudflare.namespace": document.getElementById('cloudflare_namespace').value,
				"supabase.url": document.getElementById('supabase_url').value,
				"supabase.table": document.getElementById('supabase_table').value,
				output_folder: document.getElementById('output_folder').value,
				sync_interval: parseInt(document.getElementById('sync_interval').value) || 15
			};

			if (!data.provider) {
				document.getElementById('error').textContent = 'Please select a provider';
				return;
			}

			var blob = new Blob([JSON.stringify(data, null, 2)], {type: 'application/json'});
			var url = URL.createObjectURL(blob);
			var a = document.createElement('a');
			a.href = url;
			a.download = 'config.json';
			a.click();

			document.getElementById('error').textContent = 'Config downloaded. Move to %APPDATA%\\formseal-sync\\config.json and restart.';
		}

		function openCLI() {
			window.location.href = 'formseal-sync:setup';
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