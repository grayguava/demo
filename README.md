# formseal-fetch

<p align="center">
  <img src="https://img.shields.io/badge/python-3.8+-3776ab?style=flat&labelColor=1e293b">
  <img src="https://img.shields.io/badge/license-MIT-fc8181?style=flat&labelColor=1e293b">
  <img src="https://img.shields.io/badge/formseal-ecosystem-10b981?style=flat&labelColor=1e293b">
</p>

Download encrypted form submissions from your storage backend for offline decryption.

## What it does

```
Browser (formseal-embed)
       │
       ▼ (encrypted submissions)
  Storage (Cloudflare KV / Supabase / ...)
       │
       ▼ (fsf fetch)
  Your PC ──► ciphertexts.jsonl
       │
       ▼ (decrypt offline)
  Plaintext form data
```

## Install

```bash
pipx install fsf
```

Or with pip:

```bash
pip install fsf
```

## Quick start

```bash
fsf connect provider:<name>
fsf fetch
fsf status
```

## Features

- **Secure storage** : Credentials stored in OS keychain (Windows Credential Manager / macOS Keychain / Linux Secret Service)
- **Deduplication** : Skips already-downloaded ciphertexts automatically
- **Offline-capable** ; Decrypt downloaded data anytime without network access

## Commands

| Command | Description |
|---------|-------------|
| `fsf connect` | Connect to a storage provider |
| `fsf fetch` | Download ciphertexts |
| `fsf status` | Show connection info |
| `fsf disconnect` | Clear all credentials |

Run `fsf --help` for all commands.

## Security

Your API tokens never leave your machine.formseal-fetch:
- Stores credentials in your OS keychain (encrypted at rest)
- Makes direct API calls to your storage backend only
- Sends no telemetry, has no analytics

## Documentation

Detailed guides: [docs/](./docs/)

- [Getting Started](./docs/getting-started.md)
- [Security](./docs/security.md)
- [Commands Reference](./docs/reference/commands.md)
- [Troubleshooting](./docs/troubleshooting.md)

## License

MIT