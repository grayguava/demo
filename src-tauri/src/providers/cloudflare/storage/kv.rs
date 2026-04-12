// providers/cloudflare/storage/kv.rs
// Cloudflare KV storage adapter

use anyhow::{anyhow, Result};
use reqwest::blocking::Client;
use serde::Deserialize;
use std::collections::HashSet;
use std::fs::{self, OpenOptions};
use std::io::Write;
use std::path::PathBuf;

use crate::config::ProviderConfig;

#[derive(Debug, Deserialize)]
struct KeysResponse {
    success: bool,
    result: Vec<KvKey>,
    result_info: Option<ResultInfo>,
    errors: Option<Vec<CfError>>,
}

#[derive(Debug, Deserialize)]
struct KvKey {
    name: String,
}

#[derive(Debug, Deserialize)]
struct ResultInfo {
    cursor: Option<String>,
}

#[derive(Debug, Deserialize)]
struct CfError {
    message: String,
}

pub struct FetchResult {
    pub written: usize,
    pub skipped: usize,
}

fn load_seen(path: &PathBuf) -> HashSet<String> {
    if !path.exists() {
        return HashSet::new();
    }
    fs::read_to_string(path)
        .unwrap_or_default()
        .lines()
        .filter(|l| !l.trim().is_empty())
        .map(|l| l.trim().to_string())
        .collect()
}

pub fn fetch(cfg: &ProviderConfig, account_id: &str, output_path: &PathBuf) -> Result<FetchResult> {
    if cfg.namespace.is_empty() {
        return Err(anyhow!("KV Namespace ID not set. Go to Settings."));
    }
    if cfg.token.is_empty() {
        return Err(anyhow!("API token not set. Go to Settings."));
    }

    let client = Client::new();
    let base = format!(
        "https://api.cloudflare.com/client/v4/accounts/{}/storage/kv/namespaces/{}",
        account_id, cfg.namespace
    );

    // List all keys (paginated)
    let mut all_keys: Vec<String> = Vec::new();
    let mut cursor: Option<String> = None;

    loop {
        let mut url = format!("{}/keys", base);
        if let Some(ref c) = cursor {
            url = format!("{}?cursor={}", url, c);
        }

        let resp = client
            .get(&url)
            .header("Authorization", format!("Bearer {}", cfg.token)) // Fixed: Token not Bearer
            .send()?;

        if !resp.status().is_success() {
            let status = resp.status();
            let body = resp.text().unwrap_or_default();
            return Err(anyhow!("HTTP {}: {}", status, body));
        }

        let data: KeysResponse = resp.json()?;

        if !data.success {
            let msg = data
                .errors
                .and_then(|e| e.into_iter().next().map(|e| e.message))
                .unwrap_or_else(|| "unknown error".into());
            return Err(anyhow!("Cloudflare API error: {}", msg));
        }

        all_keys.extend(data.result.into_iter().map(|k| k.name));

        cursor = data
            .result_info
            .and_then(|ri| ri.cursor)
            .filter(|c| !c.is_empty());

        if cursor.is_none() {
            break;
        }
    }

    if all_keys.is_empty() {
        return Ok(FetchResult {
            written: 0,
            skipped: 0,
        });
    }

    // Load existing ciphertexts for dedup
    fs::create_dir_all(output_path.parent().unwrap())?;
    let mut seen = load_seen(output_path);

    let mut file = OpenOptions::new()
        .create(true)
        .append(true)
        .open(output_path)?;

    let mut written = 0;
    let mut skipped = 0;

    // Fetch each value individually
    for key in &all_keys {
        let encoded = urlencoding::encode(key);
        let url = format!("{}/values/{}", base, encoded);

        let resp = client
            .get(&url)
            .header("Authorization", format!("Bearer {}", cfg.token)) // Fixed: Token not Bearer
            .send()?;

        if !resp.status().is_success() {
            continue;
        }

        let value = resp.text()?.trim().to_string();

        if value.is_empty() {
            continue;
        }

        if seen.contains(&value) {
            skipped += 1;
            continue;
        }

        writeln!(file, "{}", value)?;
        seen.insert(value);
        written += 1;
    }

    Ok(FetchResult { written, skipped })
}
