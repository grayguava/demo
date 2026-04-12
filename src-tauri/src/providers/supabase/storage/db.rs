// providers/supabase/storage/db.rs
use anyhow::{anyhow, Result};
use reqwest::blocking::Client;
use serde::Deserialize;
use std::collections::HashSet;
use std::fs::{self, OpenOptions};
use std::io::Write;
use std::path::PathBuf;

use crate::config::ProviderConfig;

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

pub fn fetch(cfg: &ProviderConfig, output_path: &PathBuf) -> Result<FetchResult> {
    if cfg.url.is_empty() {
        return Err(anyhow!("Supabase URL not set. Go to Settings."));
    }
    if cfg.key.is_empty() {
        return Err(anyhow!("Supabase key not set. Go to Settings."));
    }

    let table = cfg.table.as_deref().unwrap_or("ciphertexts");
    let client = Client::new();
    let url = format!("{}/rest/v1/{}?select=data", cfg.url, table);

    let resp = client
        .get(&url)
        .header("Authorization", format!("Bearer {}", cfg.key))
        .header("apikey", cfg.key.clone())
        .send()?;

    if !resp.status().is_success() {
        return Err(anyhow!("HTTP {}: {}", resp.status(), resp.text().unwrap_or_default()));
    }

    #[derive(Deserialize)]
    struct Row {
        data: Option<String>,
    }

    let rows: Vec<Row> = resp.json()?;

    fs::create_dir_all(output_path.parent().unwrap())?;
    let mut seen = load_seen(output_path);
    let mut file = OpenOptions::new().create(true).append(true).open(output_path)?;

    let mut written = 0;
    let mut skipped = 0;

    for row in rows {
        if let Some(value) = row.data {
            let value = value.trim().to_string();
            if value.is_empty() { continue; }
            if seen.contains(&value) { skipped += 1; continue; }
            writeln!(file, "{}", value)?;
            seen.insert(value);
            written += 1;
        }
    }

    Ok(FetchResult { written, skipped })
}