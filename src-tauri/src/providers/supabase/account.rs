// providers/supabase/account.rs
use reqwest::blocking::Client;

pub fn validate_token(url: &str, key: &str) -> Result<bool, String> {
    let client = Client::new();
    let fetch_url = format!("{}/rest/v1/", url);

    let resp = client
        .get(&fetch_url)
        .header("Authorization", format!("Bearer {}", key))
        .header("apikey", key)
        .send()
        .map_err(|e| e.to_string())?;

    Ok(resp.status().is_success())
}