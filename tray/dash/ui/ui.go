package ui

const IndexHTML = `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>FormSeal-Sync Dashboard</title>
	<style>
		* { box-sizing: border-box; margin: 0; padding: 0; }
		body { 
			font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; 
			background: #0f172a; 
			color: #e2e8f0; 
			min-height: 100vh;
			padding: 20px;
		}
		h1 { color: #f97316; margin-bottom: 20px; font-size: 24px; }
		.grid { 
			display: grid; 
			grid-template-columns: 1fr 1fr; 
			gap: 20px; 
			max-width: 1200px;
			margin: 0 auto;
		}
		.card { 
			background: #1e293b; 
			border-radius: 12px; 
			padding: 24px;
		}
		.card h2 {
			color: #64748b;
			font-size: 11px;
			text-transform: uppercase;
			letter-spacing: 1.5px;
			margin-bottom: 16px;
			border-bottom: 1px solid #334155;
			padding-bottom: 12px;
		}
		.status-value { font-size: 32px; font-weight: 700; margin: 8px 0; }
		.status-value.running { color: #22c55e; }
		.status-value.idle { color: #64748b; }
		.stat-row { display: flex; justify-content: space-between; padding: 8px 0; border-bottom: 1px solid #334155; }
		.stat-label { color: #94a3b8; }
		.stat-value { font-weight: 600; }
		.btn-group { margin-top: 16px; }
		.btn { padding: 10px 20px; border: none; border-radius: 8px; font-weight: 600; cursor: pointer; margin-right: 8px; }
		.btn-start { background: #22c55e; color: white; }
		.btn-stop { background: #ef4444; color: white; }
		.btn-edit { background: #f97316; color: white; margin-top: 16px; display: inline-block; padding: 10px 20px; border-radius: 8px; font-weight: 600; cursor: pointer; }
		.log-container { background: #0f172a; padding: 16px; border-radius: 8px; max-height: 280px; overflow-y: auto; font-family: monospace; font-size: 12px; color: #64748b; white-space: pre-wrap; line-height: 1.6; }
		.quick-stats { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; margin-bottom: 16px; }
		.quick-stat { background: #334155; padding: 12px; border-radius: 8px; text-align: center; }
		.quick-stat-value { font-size: 20px; font-weight: 700; color: #f97316; }
		.quick-stat-label { font-size: 11px; color: #94a3b8; margin-top: 4px; }
		.btn-small { padding: 8px 16px; background: #334155; color: #e2e8f0; border: none; border-radius: 6px; cursor: pointer; margin-right: 8px; }
		.modal { display: none; position: fixed; top: 0; left: 0; width: 100%; height: 100%; background: rgba(0,0,0,0.8); z-index: 100; }
		.modal.active { display: flex; align-items: center; justify-content: center; }
		.modal-content { background: #1e293b; padding: 32px; border-radius: 16px; max-width: 500px; width: 90%; }
		.modal h3 { color: #f97316; margin-bottom: 24px; }
		.form-group { margin-bottom: 16px; }
		.form-group label { display: block; color: #94a3b8; margin-bottom: 8px; }
		.form-group input, .form-group select { width: 100%; padding: 12px; background: #334155; border: 1px solid #475569; color: #e2e8f0; border-radius: 8px; }
		.modal-buttons { margin-top: 24px; display: flex; gap: 12px; }
	</style>
</head>
<body>
	<h1>FormSeal-Sync</h1>
	<div class="grid">
		<div class="card">
			<h2>Status</h2>
			<div id="status" class="status-value idle">Loading...</div>
			<div class="stat-row"><span class="stat-label">Messages</span><span class="stat-value" id="msgCount">0</span></div>
			<div class="stat-row"><span class="stat-label">Last Sync</span><span class="stat-value" id="lastSync">-</span></div>
			<div class="btn-group">
				<button class="btn btn-start" onclick="doAction('/api/start')">Start</button>
				<button class="btn btn-stop" onclick="doAction('/api/stop')">Stop</button>
			</div>
		</div>
		<div class="card">
			<h2>Configuration</h2>
			<div class="stat-row"><span class="stat-label">Provider</span><span class="stat-value" id="provider">-</span></div>
			<div class="stat-row"><span class="stat-label">Output Folder</span><span class="stat-value" id="output">-</span></div>
			<div class="stat-row"><span class="stat-label">Interval</span><span class="stat-value"><span id="interval">15</span> min</span></div>
			<button class="btn-edit" onclick="openModal()">Edit Config</button>
		</div>
		<div class="card">
			<h2>Live Logs</h2>
			<div class="log-container" id="logs">Loading...</div>
		</div>
		<div class="card">
			<h2>Quick Panel</h2>
			<div class="quick-stats">
				<div class="quick-stat"><div class="quick-stat-value" id="totalFetched">0</div><div class="quick-stat-label">Total</div></div>
				<div class="quick-stat"><div class="quick-stat-value" id="duplicates">0</div><div class="quick-stat-label">Dups</div></div>
			</div>
			<button class="btn-small" onclick="loadStatus()">Refresh</button>
		</div>
	</div>
	<div class="modal" id="configModal">
		<div class="modal-content">
			<h3>Edit Configuration</h3>
			<div class="form-group"><label>Provider</label><select id="editProvider"><option value="cloudflare">Cloudflare KV</option><option value="supabase">Supabase</option></select></div>
			<div class="form-group"><label>API Token</label><input type="password" id="editToken"></div>
			<div class="form-group"><label>Namespace / URL</label><input type="text" id="editNamespace"></div>
			<div class="form-group"><label>Table</label><input type="text" id="editTable"></div>
			<div class="form-group"><label>Output Folder</label><input type="text" id="editOutput"></div>
			<div class="form-group"><label>Sync Interval (minutes)</label><input type="number" id="editInterval"></div>
			<div class="modal-buttons">
				<button class="btn btn-start" onclick="saveConfig()">Save</button>
				<button class="btn btn-stop" onclick="closeModal()">Cancel</button>
			</div>
		</div>
	</div>
	<script>
		function loadStatus() {
			fetch('/api/status').then(r => r.json()).then(d => {
				var status = document.getElementById('status');
				status.textContent = d.running ? 'Running' : 'Idle';
				status.className = 'status-value ' + (d.running ? 'running' : 'idle');
				document.getElementById('msgCount').textContent = d.msgCount || 0;
				document.getElementById('lastSync').textContent = d.lastSync || 'Never';
				if (d.config) {
					document.getElementById('provider').textContent = d.config.provider || '-';
					document.getElementById('output').textContent = d.config.output_folder || '-';
					document.getElementById('interval').textContent = d.config.sync_interval || 15;
				}
				document.getElementById('logs').textContent = d.logs || 'No logs yet';
				document.getElementById('totalFetched').textContent = d.msgCount || 0;
			});
		}
		function doAction(url) { fetch(url, {method: 'POST'}).then(loadStatus); }
		function openModal() { 
			fetch('/api/status').then(r=>r.json()).then(d=>{
				if(d.config){
					document.getElementById('editProvider').value = d.config.provider || 'cloudflare';
					document.getElementById('editToken').value = d.config.api_token || '';
					document.getElementById('editNamespace').value = d.config.cloudflare_namespace || d.config.supabase_url || '';
					document.getElementById('editTable').value = d.config.supabase_table || '';
					document.getElementById('editOutput').value = d.config.output_folder || '';
					document.getElementById('editInterval').value = d.config.sync_interval || 15;
				}
			});
			document.getElementById('configModal').classList.add('active'); 
		}
		function closeModal() { document.getElementById('configModal').classList.remove('active'); }
		function saveConfig() {
			var data = {
				provider: document.getElementById('editProvider').value,
				api_token: document.getElementById('editToken').value,
				cloudflare_namespace: document.getElementById('editProvider').value === 'cloudflare' ? document.getElementById('editNamespace').value : '',
				supabase_url: document.getElementById('editProvider').value === 'supabase' ? document.getElementById('editNamespace').value : '',
				supabase_table: document.getElementById('editTable').value,
				output_folder: document.getElementById('editOutput').value,
				sync_interval: parseInt(document.getElementById('editInterval').value)
			};
			fetch('/api/save', {method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify(data)}).then(closeModal).then(loadStatus);
		}
		loadStatus();
		setInterval(loadStatus, 5000);
	</script>
</body>
</html>`