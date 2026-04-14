package main

var indexHTML = []byte(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>formseal-sync</title>
<style>
* { box-sizing: border-box; margin: 0; padding: 0; }
body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
  background: #0d0d0d;
  color: #d4d4d4;
  min-height: 100vh;
  padding: 40px 32px;
  font-size: 13px;
}
header {
  display: flex;
  align-items: baseline;
  gap: 10px;
  margin-bottom: 36px;
}
header h1 { font-size: 17px; font-weight: 600; color: #fff; }
header span { font-size: 12px; color: #444; }
.grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; max-width: 860px; }
.card {
  background: #161616;
  border: 1px solid #222;
  border-radius: 8px;
  padding: 20px 22px;
}
.card-title {
  font-size: 10px;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: #555;
  margin-bottom: 14px;
}
.stat { font-size: 28px; font-weight: 700; color: #fff; }
.stat.ok { color: #4ade80; }
.stat.err { color: #f87171; }
.stat.pending { color: #facc15; }
.sub { font-size: 11px; color: #444; margin-top: 4px; }
.rows { display: flex; flex-direction: column; gap: 10px; }
.row { display: flex; justify-content: space-between; align-items: center; }
.row label { color: #555; font-size: 12px; }
.row value { color: #ccc; font-weight: 500; }
hr { border: none; border-top: 1px solid #1e1e1e; margin: 14px 0; }
.btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  border: 1px solid #2a2a2a;
  border-radius: 6px;
  background: #1a1a1a;
  color: #ccc;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.1s, border-color 0.1s;
}
.btn:hover { background: #222; border-color: #333; }
.btn.primary { background: #1c2a1c; border-color: #2d4a2d; color: #4ade80; }
.btn.primary:hover { background: #223222; }
.btn:disabled { opacity: 0.4; cursor: not-allowed; }
.actions { display: flex; gap: 8px; margin-top: 16px; }
.full { grid-column: 1 / -1; }
input, select {
  width: 100%;
  padding: 8px 10px;
  background: #111;
  border: 1px solid #252525;
  border-radius: 5px;
  color: #ccc;
  font-size: 12px;
  outline: none;
  margin-top: 5px;
}
input:focus, select:focus { border-color: #3a3a3a; }
.field { margin-bottom: 14px; }
.field label { color: #555; font-size: 11px; text-transform: uppercase; letter-spacing: 0.06em; }
.note { font-size: 11px; color: #3a3a3a; margin-top: 4px; }
.token-badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 600;
}
.token-badge.set { background: #1c2a1c; color: #4ade80; }
.token-badge.unset { background: #2a1c1c; color: #f87171; }
</style>
</head>
<body>

<header>
  <h1>formseal-sync</h1>
  <span id="providerBadge">—</span>
</header>

<div class="grid">

  <!-- Status card -->
  <div class="card">
    <div class="card-title">Last Sync</div>
    <div class="stat" id="syncStat">—</div>
    <div class="sub" id="syncSub">Not yet run</div>
    <div class="actions">
      <button class="btn primary" id="syncBtn" onclick="triggerSync()">▶ Sync Now</button>
      <button class="btn" onclick="loadStatus()">↻ Refresh</button>
    </div>
  </div>

  <!-- Stats card -->
  <div class="card">
    <div class="card-title">Data</div>
    <div class="rows">
      <div class="row">
        <label>Total records</label>
        <value id="msgCount">—</value>
      </div>
      <div class="row">
        <label>Last written</label>
        <value id="lastWritten">—</value>
      </div>
      <div class="row">
        <label>Last skipped</label>
        <value id="lastSkipped">—</value>
      </div>
      <div class="row">
        <label>Token</label>
        <value><span class="token-badge unset" id="tokenBadge">not set</span></value>
      </div>
    </div>
  </div>

  <!-- Config card -->
  <div class="card full">
    <div class="card-title">Configuration</div>

    <div style="display:grid;grid-template-columns:1fr 1fr;gap:16px;">
      <div>
        <div class="field">
          <label>Provider</label>
          <select id="cfgProvider" onchange="onProviderChange()">
            <option value="cloudflare">Cloudflare KV</option>
            <option value="supabase">Supabase</option>
          </select>
        </div>
        <div class="field" id="cfNsField">
          <label>KV Namespace ID</label>
          <input type="text" id="cfgNamespace" placeholder="0f2b2dcf9dd94e4285d476043af3c26d">
        </div>
        <div class="field hidden" id="sbUrlField">
          <label>Supabase URL</label>
          <input type="text" id="cfgSupabaseURL" placeholder="https://xxx.supabase.co">
        </div>
        <div class="field hidden" id="sbTableField">
          <label>Supabase Table</label>
          <input type="text" id="cfgSupabaseTable" placeholder="submissions">
        </div>
      </div>
      <div>
        <div class="field">
          <label>Output Folder</label>
          <input type="text" id="cfgOutputFolder" placeholder="/path/to/data">
          <div class="note">ciphertexts.jsonl will be written here</div>
        </div>
        <div class="field">
          <label>Token Env Var</label>
          <input type="text" id="cfgTokenEnv" placeholder="FSYNC_CF_TOKEN">
          <div class="note">Name of the env var holding your API token. Never stored in config.</div>
        </div>
      </div>
    </div>

    <hr>
    <div style="display:flex;gap:8px;align-items:center;">
      <button class="btn primary" onclick="saveConfig()">Save Configuration</button>
      <span id="saveStatus" style="font-size:11px;color:#444;"></span>
    </div>
  </div>

</div>

<script>
function onProviderChange() {
  var p = document.getElementById('cfgProvider').value;
  document.getElementById('cfNsField').style.display = p === 'cloudflare' ? '' : 'none';
  document.getElementById('sbUrlField').style.display = p === 'supabase' ? '' : 'none';
  document.getElementById('sbTableField').style.display = p === 'supabase' ? '' : 'none';
}

function loadStatus() {
  fetch('/api/status').then(r => r.json()).then(d => {
    document.getElementById('providerBadge').textContent = d.provider || '—';
    document.getElementById('msgCount').textContent = d.msgCount ?? '—';

    var badge = document.getElementById('tokenBadge');
    badge.textContent = d.tokenSet ? 'set' : 'not set';
    badge.className = 'token-badge ' + (d.tokenSet ? 'set' : 'unset');

    var r = d.result;
    var stat = document.getElementById('syncStat');
    var sub = document.getElementById('syncSub');

    if (!r || !r.run_at || r.run_at === '0001-01-01T00:00:00Z') {
      stat.textContent = '—';
      stat.className = 'stat';
      sub.textContent = 'Not yet run';
    } else if (!r.done) {
      stat.textContent = 'Running…';
      stat.className = 'stat pending';
      sub.textContent = 'Sync in progress';
      setTimeout(loadStatus, 1000);
    } else if (r.error) {
      stat.textContent = 'Error';
      stat.className = 'stat err';
      sub.textContent = r.error;
    } else {
      stat.textContent = r.written + ' new';
      stat.className = 'stat ok';
      sub.textContent = r.skipped + ' duplicates skipped · ' + new Date(r.run_at).toLocaleTimeString();
    }

    document.getElementById('lastWritten').textContent = r && r.done ? r.written : '—';
    document.getElementById('lastSkipped').textContent = r && r.done ? r.skipped : '—';
  });
}

function loadConfig() {
  fetch('/api/config').then(r => r.json()).then(cfg => {
    document.getElementById('cfgProvider').value = cfg.provider || 'cloudflare';
    document.getElementById('cfgNamespace').value = cfg.cloudflare_namespace || '';
    document.getElementById('cfgSupabaseURL').value = cfg.supabase_url || '';
    document.getElementById('cfgSupabaseTable').value = cfg.supabase_table || '';
    document.getElementById('cfgOutputFolder').value = cfg.output_folder || '';
    document.getElementById('cfgTokenEnv').value = cfg.token_env || 'FSYNC_CF_TOKEN';
    onProviderChange();
  });
}

function saveConfig() {
  var data = {
    provider: document.getElementById('cfgProvider').value,
    cloudflare_namespace: document.getElementById('cfgNamespace').value,
    supabase_url: document.getElementById('cfgSupabaseURL').value,
    supabase_table: document.getElementById('cfgSupabaseTable').value,
    output_folder: document.getElementById('cfgOutputFolder').value,
    token_env: document.getElementById('cfgTokenEnv').value,
  };
  fetch('/api/config', {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(data)
  }).then(r => r.json()).then(() => {
    var s = document.getElementById('saveStatus');
    s.textContent = 'Saved.';
    s.style.color = '#4ade80';
    setTimeout(() => { s.textContent = ''; }, 2000);
    loadStatus();
  });
}

function triggerSync() {
  document.getElementById('syncBtn').disabled = true;
  fetch('/api/sync', {method: 'POST'}).then(() => {
    setTimeout(function poll() {
      fetch('/api/status').then(r => r.json()).then(d => {
        if (d.result && !d.result.done) {
          setTimeout(poll, 800);
        } else {
          loadStatus();
          document.getElementById('syncBtn').disabled = false;
        }
      });
    }, 500);
  });
}

loadConfig();
loadStatus();
</script>
</body>
</html>`)