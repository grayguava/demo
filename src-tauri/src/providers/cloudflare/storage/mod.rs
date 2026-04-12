// providers/cloudflare/storage/mod.rs
pub mod kv;

pub use kv::{fetch as fetch_kv, FetchResult};
