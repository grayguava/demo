# Getting started

This guide will help you install formseal-fetch and connect to your storage backend within minutes.

## Prerequisites

- **Python 3.8 or later**
- **Operating System**: Windows, macOS, or Linux
- **Network access** to your storage backend

## Installation

### With pipx (recommended)

```bash
pipx install fsf
```

### With pip

```bash
pip install fsf
```

### From source

```bash
git clone https://github.com/grayguava/formseal-fetch.git
cd formseal-fetch
pip install -e .
```

### Verify installation

```bash
fsf --help
```

You should see the command help output.

## Quick start

### Step 1: Connect to your backend

```bash
fsf connect provider:<name>
```

You'll be prompted for provider-specific credentials:
- **Namespace ID**: Your storage namespace identifier
- **API Token**: Your API token with read permissions
- **Output Folder**: Where to save downloaded ciphertexts (default: `data`)

You can also provide these non-interactively:

```bash
fsf connect provider:<name> namespace:<id> token:<value> output:<path>
```

### Step 2: Fetch ciphertexts

```bash
fsf fetch
```

This downloads all encrypted form submissions from your storage namespace to `data/ciphertexts.jsonl`. Each line is a raw ciphertext string (no JSON wrapper), one per line.

### Step 3: Check connection and configuration status

```bash
fsf status
```

Shows your current provider, namespace ID, token location, and output folder.

## Next steps

- See [Commands reference](./reference/commands.md) for all available commands
- Read [Security](./security.md) to understand how credentials are stored
- Check [Troubleshooting](./troubleshooting.md) if you encounter issues