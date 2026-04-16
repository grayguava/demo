# Commands reference

Complete reference for all formseal-fetch commands.

## Usage syntax

```bash
fsf <command> [options] [arguments]
```

## Commands

### connect

Connect to a storage backend.

```bash
fsf connect provider:<name> [namespace:<id>] [token:<value>] [output:<path>]
```

**Arguments:**

| Argument | Description |
|----------|-------------|
| `provider:<name>` | Storage provider (required) — available: `cloudflare` |
| `namespace:<id>` | Storage-specific namespace ID (optional) |
| `token:<value>` | API token (optional — you'll be prompted if not provided) |
| `output:<path>` | Output folder for ciphertexts (optional, default: `data`) |

**Examples:**

```bash
# Interactive mode — you'll be prompted for all required values
fsf connect provider:<name>

# Non-interactive — all values provided via arguments
fsf connect provider:<name> namespace:<id> token:<value> output:<path>

# Specify only namespace, enter token interactively
fsf connect provider:<name> namespace:<id>
```

**Interactive prompts:**

When running without all arguments, you'll be prompted for:
1. KV Namespace ID (required for Cloudflare)
2. Account API Token (required)
3. Output Folder (optional, default: `data`)

Press `Ctrl+C` at any prompt to cancel.

---

### fetch

Download ciphertexts from your connected backend.

```bash
fsf fetch [--output <path>]
```

**Options:**

| Option | Description |
|--------|-------------|
| `--output` | Custom output file path (default: `<output_folder>/ciphertexts.jsonl`) |

**Examples:**

```bash
# Use default output folder
fsf fetch

# Custom output file
fsf fetch --output my-data.jsonl
```

**Output format:**

Ciphertexts are saved as plain text — one raw ciphertext per line, no JSON formatting:

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
eyJraWQiOiIxMjM0NTY3ODkwIn0...
```

**Deduplication:**

If you run `fsf fetch` multiple times, duplicates are automatically skipped based on ciphertext content. The output shows:
- Number of new ciphertexts saved
- Number of duplicates skipped

---

### status

Show current connection status and configuration.

```bash
fsf status
```

**Output includes:**

- Provider name (e.g., "Cloudflare")
- KV Namespace ID (truncated)
- Storage location (OS Keychain or Config File)
- Account ID (truncated from API)
- API Token status (**** if set)
- Output folder path

---

### set

Set a configuration value.

```bash
fsf set <key> <value>
```

**Arguments:**

| Key | Description |
|-----|-------------|
| `output_folder` | Path where ciphertexts are saved |

**Examples:**

```bash
fsf set output_folder my-data
```

---

### disconnect

Clear all credentials and configuration.

```bash
fsf disconnect
```

**What it removes:**

- Provider configuration
- API token (from OS Keychain or secrets.json)
- KV namespace ID (from OS Keychain or secrets.json)
- Configuration file

**What it does NOT remove:**

- Downloaded ciphertexts in your output folder

**Confirmation:**

You'll be prompted to confirm with `y` or `n`. Press `Enter` to cancel.

---

### providers

List available storage backends.

```bash
fsf providers
```

**Currently supported:**

| Provider | Status |
|----------|--------|
| Cloudflare KV | Available |
| Supabase | Coming soon |

---

### --help

Show help information.

```bash
fsf --help
```

Displays all available commands grouped by category (Connect, Fetch, Config, Info).

---

### --about

Show project information.

```bash
fsf --about
```

Displays version and repository URL.