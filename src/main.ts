import { invoke } from "@tauri-apps/api/core";
import { open } from "@tauri-apps/plugin-dialog";
import { getCurrentWindow } from "@tauri-apps/api/window";

const appWindow = getCurrentWindow();

// --- window controls ---
document.getElementById("btn-minimize")!.addEventListener("click", () => appWindow.minimize());
document.getElementById("btn-close")!.addEventListener("click", () => appWindow.close());

// --- tabs ---
document.querySelectorAll(".nav-tab").forEach((tab) => {
  tab.addEventListener("click", () => {
    const target = (tab as HTMLElement).dataset.tab!;
    document.querySelectorAll(".nav-tab").forEach((t) => t.classList.remove("active"));
    document.querySelectorAll(".tab-content").forEach((c) => c.classList.remove("active"));
    tab.classList.add("active");
    document.getElementById(`tab-${target}`)!.classList.add("active");
  });
});

// --- provider selection ---
document.getElementById("provider")!.addEventListener("change", (e) => {
  const provider = (e.target as HTMLSelectElement).value;
  (document.getElementById("cf-fields") as HTMLElement).style.display = provider === "cloudflare" ? "block" : "none";
  (document.getElementById("su-fields") as HTMLElement).style.display = provider === "supabase" ? "block" : "none";
});

// --- log ---
function log(msg: string, cls: "ok" | "err" | "dim" | "accent" | "" = "") {
  const el = document.getElementById("log")!;
  const line = document.createElement("div");
  line.className = `log-line ${cls}`;
  const ts = new Date().toLocaleTimeString("en-GB");
  line.textContent = `${ts}  ${msg}`;
  el.appendChild(line);
  el.scrollTop = el.scrollHeight;
}

function clearLog() {
  document.getElementById("log")!.innerHTML = "";
}

// --- badge ---
function setBadge(state: "idle" | "ok" | "err", label: string) {
  const badge = document.getElementById("status-badge")!;
  badge.className = `badge ${state === "idle" ? "" : state}`;
  badge.textContent = label;
}

function setTs(text: string) {
  document.getElementById("status-ts")!.textContent = text;
}

// --- stats ---
async function refreshStats() {
  try {
    const result = await invoke<{ total: number }>("get_stats");
    document.getElementById("stat-total")!.textContent = result.total.toLocaleString();
  } catch (_) {
    document.getElementById("stat-total")!.textContent = "—";
  }
}

// --- fetch ---
let lastNewCount = 0;

document.getElementById("btn-fetch")!.addEventListener("click", async () => {
  const btn = document.getElementById("btn-fetch") as HTMLButtonElement;
  btn.disabled = true;
  clearLog();
  setBadge("idle", "fetching...");
  log("Connecting to provider...", "dim");

  try {
    const result = await invoke<{ written: number; skipped: number }>("fetch_ciphertexts");
    lastNewCount = result.written;

    log(`Fetched from provider`, "dim");
    log(`${result.written} new · ${result.skipped} duplicates skipped`, "ok");
    log(`Saved to output folder`, "dim");

    setBadge("ok", "last fetch ok");
    setTs(new Date().toLocaleString("en-GB", { hour: "2-digit", minute: "2-digit", day: "2-digit", month: "short" }));

    document.getElementById("stat-new")!.textContent = result.written.toLocaleString();
    await refreshStats();
  } catch (e: unknown) {
    const msg = e instanceof Error ? e.message : String(e);
    log(`Error: ${msg}`, "err");
    setBadge("err", "fetch failed");
  } finally {
    btn.disabled = false;
  }
});

// --- settings: load ---
interface Config {
  provider: string;
  cloudflare: { token: string; namespace: string };
  supabase: { url: string; key: string; table: string | null };
  output_folder: string;
}

async function loadSettings() {
  try {
    const cfg = await invoke<Config>("get_config");

    // Set provider
    (document.getElementById("provider") as HTMLSelectElement).value = cfg.provider || "";

    // Show/hide fields based on provider
    if (cfg.provider === "cloudflare") {
      (document.getElementById("cf-fields") as HTMLElement).style.display = "block";
      (document.getElementById("su-fields") as HTMLElement).style.display = "none";
      (document.getElementById("cf-token") as HTMLInputElement).value = cfg.cloudflare?.token || "";
      (document.getElementById("cf-namespace") as HTMLInputElement).value = cfg.cloudflare?.namespace || "";
    } else if (cfg.provider === "supabase") {
      (document.getElementById("cf-fields") as HTMLElement).style.display = "none";
      (document.getElementById("su-fields") as HTMLElement).style.display = "block";
      (document.getElementById("su-url") as HTMLInputElement).value = cfg.supabase?.url || "";
      (document.getElementById("su-key") as HTMLInputElement).value = cfg.supabase?.key || "";
      (document.getElementById("su-table") as HTMLInputElement).value = cfg.supabase?.table || "";
    }

    (document.getElementById("output-folder") as HTMLInputElement).value = cfg.output_folder || "";
  } catch (_) {}
}

// --- settings: browse ---
document.getElementById("btn-browse")!.addEventListener("click", async () => {
  const selected = await open({ directory: true, multiple: false, title: "Select output folder" });
  if (selected && typeof selected === "string") {
    (document.getElementById("output-folder") as HTMLInputElement).value = selected;
  }
});

// --- settings: save ---
document.getElementById("btn-save")!.addEventListener("click", async () => {
  const provider = (document.getElementById("provider") as HTMLSelectElement).value;
  const cloudflare_token = (document.getElementById("cf-token") as HTMLInputElement).value.trim();
  const cloudflare_namespace = (document.getElementById("cf-namespace") as HTMLInputElement).value.trim();
  const supabase_url = (document.getElementById("su-url") as HTMLInputElement).value.trim();
  const supabase_key = (document.getElementById("su-key") as HTMLInputElement).value.trim();
  const supabase_table = (document.getElementById("su-table") as HTMLInputElement).value.trim() || null;
  const output_folder = (document.getElementById("output-folder") as HTMLInputElement).value.trim();

  const msg = document.getElementById("save-msg")!;

  if (!provider) {
    msg.textContent = "Select a provider first.";
    msg.style.color = "var(--err)";
    return;
  }

  try {
    await invoke("save_config", { 
      provider,
      cloudflare_token: cloudflare_token,
      cloudflare_namespace: cloudflare_namespace,
      supabase_url: supabase_url,
      supabase_key: supabase_key,
      supabase_table: supabase_table,
      output_folder: output_folder 
    });
    msg.textContent = "Saved.";
    msg.style.color = "var(--ok)";
    setTimeout(() => (msg.textContent = ""), 2000);
  } catch (e: unknown) {
    const errMsg = e instanceof Error ? e.message : String(e);
    msg.textContent = `Error: ${errMsg}`;
    msg.style.color = "var(--err)";
  }
});

// --- init ---
loadSettings();
refreshStats();