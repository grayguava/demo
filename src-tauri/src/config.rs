use anyhow::Result;
use serde::{Deserialize, Serialize};
use std::fs;
use std::path::PathBuf;

#[derive(Debug, Clone, Serialize, Deserialize, Default)]
pub struct ProviderConfig {
    pub token: String,
    pub namespace: String,
    pub url: String,
    pub key: String,
    pub table: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize, Default)]
pub struct Config {
    pub provider: String,
    #[serde(default)]
    pub cloudflare: ProviderConfig,
    #[serde(default)]
    pub supabase: ProviderConfig,
    pub output_folder: String,
}

fn config_path() -> PathBuf {
    let mut path = dirs::home_dir().expect("no home dir");
    path.push(".formseal-sync");
    path.push("config.json");
    path
}

pub fn load() -> Config {
    let path = config_path();
    if !path.exists() {
        return Config::default();
    }
    let text = fs::read_to_string(&path).unwrap_or_default();
    serde_json::from_str(&text).unwrap_or_default()
}

pub fn save(cfg: &Config) -> Result<()> {
    let path = config_path();
    fs::create_dir_all(path.parent().unwrap())?;
    let text = serde_json::to_string_pretty(cfg)?;
    fs::write(&path, text)?;
    Ok(())
}
