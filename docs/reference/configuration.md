# Configuration

This document explains formseal-fetch's configuration files and storage locations.

## Configuration directory

All configuration is stored in:

```
~/.config/formseal-fetch/
```

| File | Contents | Format |
|------|----------|--------|
| `config.json` | Provider, output folder | Plain JSON |
| `secrets.json` | API token, namespace ID | Base64-encoded JSON (fallback only) |

## config.json

Created when you run `fsf connect`. Contains non-sensitive configuration:

```json
{
  "provider": "cloudflare",
  "output_folder": "data"
}
```

**Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `provider` | string | Storage backend name (e.g., "cloudflare") |
| `output_folder` | string | Path to save ciphertexts |

### Modifying configuration

Use `fsf set` to change configuration values:

```bash
fsf set output_folder my-folder
```

Or edit `config.json` directly (not recommended).

## secrets.json

**Only exists** if OS keychain is unavailable and fallback was triggered. Contains base64-encoded credentials:

```json
{
  "cloudflare:token": "dG9rZW4tZXhhbXBsZQ==",
  "cloudflare:namespace": "bmFtZXNwYWNlLWlk"
}
```

**Note**: This file is not encrypted — only base64 encoded. It exists as a last-resort fallback.

## OS keychain

Credentials stored in OS keychain are **encrypted** by the operating system:

- **Windows**: `cmdkey` or Credential Manager
- **macOS**: `security find-internet-password`
- **Linux**: `secret-tool` or DBus Secret Service

You can view stored credentials via OS utilities, but formseal-fetch manages them automatically.

## Runtime behavior

1. On `fsf connect`: Credentials saved to OS keychain (preferred) or secrets.json (fallback)
2. On `fsf fetch`: Credentials loaded from OS keychain first, then secrets.json
3. On `fsf disconnect`: Credentials deleted from both OS keychain and secrets.json

## Environment variables

formseal-fetch does **not** use environment variables for configuration. All settings are stored in the config directory.

## Removing configuration

```bash
fsf disconnect
```

This removes:
- `config.json`
- `secrets.json` (if exists)
- Credentials from OS keychain

The output folder and its contents are **not** affected.