# Cloudflare KV

Connecting formseal-fetch to your Cloudflare KV namespace.

## What you need

To connect formseal-fetch to Cloudflare KV, you'll need:

| Item | Where to find it |
|------|------------------|
| KV Namespace ID | Cloudflare Dashboard → Workers → KV |
| API Token | Cloudflare Dashboard → Profile → API Tokens |

## Getting your credentials

### KV Namespace ID

From Cloudflare Dashboard:
1. Go to **Workers** → **KV**
2. Find your namespace
3. Copy the **ID** column (32-character string like `0f2b2dcf9dd94e4285d476043af3c26d`)

### API Token

Your token must have these permissions:

| Permission | Required for |
|------------|---------------|
| Account Settings: Read | Listing your account |
| Workers KV: Read | Reading from KV namespace |

**To create or check your token:**

1. Go to Cloudflare Dashboard → Profile → API Tokens
2. Ensure you have a token with the above permissions
3. If creating new: use "Create Custom Token" with:
   - Account: Read
   - Workers KV: Read

## Connecting

### Interactive mode

```bash
fsf connect provider:cloudflare
```

You'll be prompted for:
- KV Namespace ID
- Account API Token
- Output folder (default: `data`)

### Non-interactive mode

```bash
fsf connect provider:cloudflare namespace:<id> token:<value> output:<path>
```

**Example:**

```bash
fsf connect provider:cloudflare namespace:<id> token:<value> output:data
```

## Fetching ciphertexts

```bash
fsf fetch
```

Ciphertexts are saved to `data/ciphertexts.jsonl` — one raw ciphertext per line.

## Verifying connection

```bash
fsf status
```

Shows your namespace ID (truncated), account ID (truncated), and token location (OS Keychain or Config File).

## Common issues

| Error | Cause | Solution |
|-------|-------|----------|
| "No accounts found" | Token missing Account: Read | Update token permissions |
| "Auth failed" | Invalid or expired token | Create new token |
| "HTTP 403" | Token missing KV: Read | Add Workers KV permission |
| "No namespace" | Namespace ID not stored | Reconnect with namespace ID |