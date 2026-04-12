const { invoke } = window.__TAURI__.core;
const { open } = window.__TAURI__.dialog;
const { getCurrentWindow } = window.__TAURI__.window;

const appWindow = getCurrentWindow();

// window controls
document.getElementById("btn-minimize").addEventListener("click", () => appWindow.minimize());
document.getElementById("btn-close").addEventListener("click", () => appWindow.close());

// tabs
document.querySelectorAll(".nav-tab").forEach(tab => {
  tab.addEventListener("click", () => {
    const target = tab.dataset.tab;
    document.querySelectorAll(".nav-tab").forEach(t => t.classList.remove("active"));
    document.querySelectorAll(".tab-content").forEach(c => c.classList.remove("active"));
    tab.classList.add("active");
    document.getElementById("tab-" + target).classList.add("active");
  });
});

// provider toggle
document.getElementById("provider").addEventListener("change", function () {
  document.getElementById("cf-fields").style.display = this.value === "cloudflare" ? "block" : "none";
  document.getElementById("sb-fields").style.display = this.value === "supabase"   ? "block" : "none";
});

// log
function log(msg, cls = "") {
  const el = document.getElementById("log");
  const line = document.createElement("div");
  line.className = "log-line " + cls;
  const ts = new Date().toLocaleTimeString("en-GB");
  line.textContent = ts + "  " + msg;
  el.appendChild(line);
  el.scrollTop = el.scrollHeight;
}

function clearLog() {
  document.getElementById("log").innerHTML = "";
}

function setBadge(state, label) {
  const badge = document.getElementById("status-badge");
  badge.className = "badge " + (state === "idle" ? "" : state);
  badge.textContent = label;
}

// stats
async function refreshStats() {
  try {
    const result = await invoke("get_stats");
    document.getElementById("stat-total").textContent = result.total.toLocaleString();
  } catch (_) {
    document.getElementById("stat-total").textContent = "—";
  }
}

// fetch
document.getElementById("btn-fetch").addEventListener("click", async () => {
  const btn = document.getElementById("btn-fetch");
  btn.disabled = true;
  clearLog();
  setBadge("idle", "fetching...");
  log("Connecting...", "dim");

  try {
    const result = await invoke("fetch_ciphertexts");
    log(result.written + " new · " + result.skipped + " duplicates skipped", "ok");
    log("Saved to output folder", "dim");
    setBadge("ok", "last fetch ok");
    document.getElementById("status-ts").textContent =
      new Date().toLocaleString("en-GB", { hour: "2-digit", minute: "2-digit", day: "2-digit", month: "short" });
    document.getElementById("stat-new").textContent = result.written.toLocaleString();
    await refreshStats();
  } catch (e) {
    log("Error: " + String(e), "err");
    setBadge("err", "fetch failed");
  } finally {
    btn.disabled = false;
  }
});

// settings load
async function loadSettings() {
  try {
    const cfg = await invoke("get_config");
    document.getElementById("provider").value        = cfg.provider || "cloudflare";
    document.getElementById("cf-token").value        = cfg.cloudflare?.token || "";
    document.getElementById("cf-namespace").value    = cfg.cloudflare?.namespace || "";
    document.getElementById("sb-url").value          = cfg.supabase?.url || "";
    document.getElementById("sb-key").value          = cfg.supabase?.key || "";
    document.getElementById("sb-table").value        = cfg.supabase?.table || "";
    document.getElementById("output-folder").value   = cfg.output_folder || "";
    document.getElementById("provider").dispatchEvent(new Event("change"));
  } catch (_) {}
}

// browse
document.getElementById("btn-browse").addEventListener("click", async () => {
  const selected = await open({ directory: true, multiple: false, title: "Select output folder" });
  if (selected && typeof selected === "string") {
    document.getElementById("output-folder").value = selected;
  }
});

// save
document.getElementById("btn-save").addEventListener("click", async () => {
  const msg = document.getElementById("save-msg");
  try {
    await invoke("save_config", {
      provider:           document.getElementById("provider").value,
      cloudflareToken:    document.getElementById("cf-token").value.trim(),
      cloudflareNamespace:document.getElementById("cf-namespace").value.trim(),
      supabaseUrl:        document.getElementById("sb-url").value.trim(),
      supabaseKey:        document.getElementById("sb-key").value.trim(),
      supabaseTable:      document.getElementById("sb-table").value.trim() || null,
      outputFolder:       document.getElementById("output-folder").value.trim(),
    });
    msg.style.color = "var(--ok)";
    msg.textContent = "Saved.";
    setTimeout(() => msg.textContent = "", 2000);
  } catch (e) {
    msg.style.color = "var(--err)";
    msg.textContent = "Error: " + String(e);
  }
});

// init
loadSettings();
refreshStats();