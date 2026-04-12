mod config;
mod providers;

use serde::Serialize;
use std::fs;
use std::path::PathBuf;

#[derive(Serialize)]
struct FetchResult {
    written: usize,
    skipped: usize,
}

#[derive(Serialize)]
struct StatsResult {
    total: usize,
}

#[tauri::command]
fn get_config() -> config::Config {
    config::load()
}

#[tauri::command]
fn save_config(
    provider: String,
    cloudflare_token: String,
    cloudflare_namespace: String,
    supabase_url: String,
    supabase_key: String,
    supabase_table: Option<String>,
    output_folder: String,
) -> Result<(), String> {
    let cfg = config::Config {
        provider,
        cloudflare: config::ProviderConfig {
            token: cloudflare_token,
            namespace: cloudflare_namespace,
            ..Default::default()
        },
        supabase: config::ProviderConfig {
            url: supabase_url,
            key: supabase_key,
            table: supabase_table,
            ..Default::default()
        },
        output_folder,
    };
    config::save(&cfg).map_err(|e: anyhow::Error| e.to_string())
}

#[tauri::command]
fn fetch_ciphertexts() -> Result<FetchResult, String> {
    let cfg = config::load();

    if cfg.output_folder.is_empty() {
        return Err("Output folder not set. Go to Settings.".to_string());
    }

    let output_path = PathBuf::from(&cfg.output_folder).join("ciphertexts.jsonl");

    let (written, skipped) = match cfg.provider.as_str() {
        "cloudflare" => {
            if cfg.cloudflare.token.is_empty() {
                return Err("Cloudflare token not set. Go to Settings.".to_string());
            }
            if cfg.cloudflare.namespace.is_empty() {
                return Err("KV Namespace ID not set. Go to Settings.".to_string());
            }
            let account_id = providers::cloudflare::get_account_id(&cfg.cloudflare.token)
                .map_err(|e| e.to_string())?;
            let r = providers::cloudflare::fetch(&cfg.cloudflare, &account_id, &output_path)
                .map_err(|e: anyhow::Error| e.to_string())?;
            (r.written, r.skipped)
        }
        "supabase" => {
            if cfg.supabase.url.is_empty() {
                return Err("Supabase URL not set. Go to Settings.".to_string());
            }
            if cfg.supabase.key.is_empty() {
                return Err("Supabase key not set. Go to Settings.".to_string());
            }
            let r = providers::supabase::fetch(&cfg.supabase, &output_path)
                .map_err(|e: anyhow::Error| e.to_string())?;
            (r.written, r.skipped)
        }
        _ => return Err("No provider set. Go to Settings.".to_string()),
    };

    Ok(FetchResult { written, skipped })
}

#[tauri::command]
fn get_stats() -> StatsResult {
    let cfg = config::load();
    if cfg.output_folder.is_empty() {
        return StatsResult { total: 0 };
    }
    let path = PathBuf::from(&cfg.output_folder).join("ciphertexts.jsonl");
    if !path.exists() {
        return StatsResult { total: 0 };
    }
    let total = fs::read_to_string(&path)
        .unwrap_or_default()
        .lines()
        .filter(|l| !l.trim().is_empty())
        .count();
    StatsResult { total }
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_dialog::init())
        .plugin(tauri_plugin_shell::init())
        .invoke_handler(tauri::generate_handler![
            get_config,
            save_config,
            fetch_ciphertexts,
            get_stats,
        ])
        .run(tauri::generate_context!())
        .expect("error running formseal-sync");
}