# Documentation

Welcome to the formseal-fetch documentation.

## Quick links

| Guide | Description |
|-------|-------------|
| [Getting Started](./getting-started.md) | Installation and first-time setup |
| [Security](./security.md) | How credentials are stored and protected |
| [Commands Reference](./reference/commands.md) | All available commands |
| [Configuration](./reference/configuration.md) | Config files and storage |
| [Cloudflare KV](./backends/cloudflare-kv.md) | Cloudflare backend setup |
| [Supabase](./backends/supabase.md) | Supabase backend (coming soon) |
| [Troubleshooting](./troubleshooting.md) | Common issues and solutions |

## What is formseal-fetch?

formseal-fetch is a CLI tool that downloads encrypted form submissions from your storage backend (currently we only support Cloudflare KV). Use it together with [formseal-embed](https://github.com/grayguava/formseal-embed) — the client-side form encryption library.

## Workflow

```
Browser (formseal-embed)
       │
       ▼ (encrypted submissions)
  Cloudflare KV
       │
       ▼ (fsf fetch)
  Your PC ──► ciphertexts.jsonl
       │
       ▼ (decrypt offline)
  Plaintext form data
```