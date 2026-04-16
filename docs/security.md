# Security

This document explains how formseal-fetch handles your credentials and sensitive data.

## Credential storage

formseal-fetch stores sensitive data (API tokens, namespace IDs) in your operating system's secure credential storage:

| OS | Storage location |
|---|------------------|
| Windows | Credential Manager |
| macOS | Keychain |
| Linux | Secret Service API (libsecret) |

### Why OS keychain?

- **Encrypted at rest** : Operating systems encrypt credential storage
- **Access controlled** : Requires your user account to access
- **Managed by OS** : Leverages built-in security features

## Fallback behavior

If the OS keychain is unavailable, formseal-fetch falls back to storing credentials in an encrypted JSON file at:

```
~/.config/formseal-fetch/secrets.json
```

Data in this file is **base64-encoded only** — not encrypted. This fallback exists for environments where keyring is not available (some Linux containers, minimal environments).

### Detecting where credentials are stored

Run `fsf status` — it shows "OS Keychain" or "Config File" next to token/namespace fields.

## What gets stored

| Data | Stored As | Location |
|------|-----------|----------|
| API Token | Encrypted | OS Keychain (preferred) or secrets.json |
| KV Namespace ID | Encrypted | OS Keychain (preferred) or secrets.json |
| Provider name | Plaintext | ~/.config/formseal-fetch/config.json |
| Output folder path | Plaintext | ~/.config/formseal-fetch/config.json |

## Clearing credentials

To remove all stored credentials:

```bash
fsf disconnect
```

This deletes:
- API token from OS Keychain (or secrets.json)
- KV namespace ID from OS Keychain (or secrets.json)
- Configuration file (`config.json`)

**Note**: Downloaded ciphertexts are **not** affected.

## Best practices

1. **Use minimum-required permissions** for your API token
2. **Rotate tokens periodically** — disconnect and reconnect with a new token
3. **Never share your output folder** — it contains encrypted form data
4. **Use `fsf disconnect`** when done, especially on shared machines

## Security considerations

- **Token visibility**: `fsf status` shows only `****` for the token, never the full value
- **Network traffic**: All API calls go directly to your storage backend — no proxy, no telemetry
- **Local-only**: formseal-fetch never sends data anywhere except your configured backend