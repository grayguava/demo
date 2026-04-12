// providers/cloudflare/account.rs
// Cloudflare account authentication

use reqwest::blocking::Client;
use serde::Deserialize;

#[derive(Debug, Deserialize)]
struct UserResponse {
    success: bool,
    result: UserResult,
    errors: Option<Vec<CfError>>,
}

#[derive(Debug, Deserialize)]
struct UserResult {
    accounts: Vec<Account>,
}

#[derive(Debug, Deserialize)]
struct Account {
    id: String,
}

#[derive(Debug, Deserialize)]
struct CfError {
    message: String,
}

pub fn get_account_id(token: &str) -> Result<String, String> {
    let client = Client::new();
    let url = "https://api.cloudflare.com/client/v4/user";

    let resp = client
        .get(url)
        .header("Authorization", format!("Bearer {}", token))
        .send()
        .map_err(|e| e.to_string())?;

    if !resp.status().is_success() {
        return Err(format!("HTTP {}", resp.status()));
    }

    let data: UserResponse = resp.json().map_err(|e| e.to_string())?;

    if !data.success {
        let msg = data
            .errors
            .and_then(|e| e.into_iter().next().map(|e| e.message))
            .unwrap_or_else(|| "unknown error".to_string());
        return Err(msg);
    }

    data.result
        .accounts
        .into_iter()
        .next()
        .map(|a| a.id)
        .ok_or_else(|| "No accounts found".to_string())
}
