# Troubleshooting

Solutions to common issues with formseal-fetch.

## 1. Installation issues

### "Command not found: fsf"

**Cause**: Package not installed or PATH not updated.

**Solution**:

```bash
# Verify installation
pip show fsf

# If installed but not found, add Python Scripts to PATH
# Windows: Add C:\Users\<you>\AppData\Local\Programs\Python\Python314\Scripts to PATH
# macOS/Linux: Typically added automatically via pip
```

### "ModuleNotFoundError: No module named 'keyring'"

**Cause**: Keyring package not installed.

**Solution**:

```bash
pip install keyring
```

---

## 2. Connection issues

### "No provider configured. Run: fsf connect provider:<name>"

**Cause**: Never connected or configuration deleted.

**Solution**:

```bash
fsf connect provider:<name>
```

### "No token. Run: fsf connect provider:<name> to set token"

**Cause**: Connected but no API token stored.

**Solution**:

```bash
# Reconnect to store a new token
fsf connect provider:<name> token:<value>
```

### "Invalid format: ... Use flag:value (e.g., provider:<name>)"

**Cause**: Incorrect argument syntax.

**Solution**: Use colon (`:`) to separate flag and value:

```bash
# Correct
fsf connect provider:<name> namespace:<id>

# Incorrect
fsf connect provider=<name>
```

---

## 3. Authentication errors

### "Auth failed: HTTP 401: Invalid token"

**Cause**: API token is invalid or expired.

**Solution**:
1. Verify token is correct (no extra spaces)
2. Create a new token in your provider dashboard
3. Reconnect with new token:

```bash
fsf disconnect
fsf connect provider:<name> token:<value>
```

### "No accounts found. Token needs 'Account Settings: Read' scope"

**Cause**: Token lacks required permissions.

**Solution**:
1. Go to Cloudflare Dashboard → Profile → API Tokens
2. Edit your token or create a new one
3. Ensure it has:
   - **Account**: Read
   - **Workers KV**: Read

---

## 4. Fetch issues

### "HTTP 403: Authorization header not provided"

**Cause**: Token was cleared or not saved properly.

**Solution**:

```bash
# Check status to see token location
fsf status

# If token shows as "not set" or "Config File", reconnect
fsf connect provider:<name> token:<value>
```

### "No namespace. Run: fsf connect provider:<name> to set namespace"

**Cause**: Namespace ID not stored.

**Solution**:

```bash
fsf connect provider:<name> namespace:<id>
```

### "Fetch failed: HTTP 403"

**Cause**: Token doesn't have KV read permission.

**Solution**: Update token permissions in Cloudflare Dashboard.

---

## 5. Credential storage issues

### "Token Location: Config file"

**Cause**: OS keychain unavailable; fell back to JSON file.

**This is not an error** — credentials still work. This happens when:
- Running in a container without keyring support
- Keyring service not available on your Linux distribution

To verify keyring works, see [Security](../security.md).

### Credentials not persisting after reboot

**Cause**: Likely using fallback JSON storage and file was deleted.

**Solution**:
1. Reconnect to restore credentials to keychain
2. Verify "Token Location" shows "OS Keychain" in `fsf status`

---

## 6. Output issues

### "Permission denied" when writing to output folder

**Cause**: No write permission to the output directory.

**Solution**:

```bash
# Check current output folder
fsf status

# Change to a folder you have write access to
fsf set output_folder ~/formseal-data
```

### File written but empty

**Cause**: KV namespace is empty — no keys to fetch.

**Solution**: Verify your KV namespace contains data via Cloudflare Dashboard.

---

## 7. Keyboard interrupt

### Ctrl+C during interactive prompt

**Behavior**: Prompt cancels gracefully, no changes made.

This is intentional — you can safely interrupt at any prompt.

---

## 8. Still stuck?

1. Run with verbose error output:
   ```bash
   fsf fetch 2>&1
   ```

2. Check [GitHub Issues](https://github.com/formseal/formseal-fetch/issues)
3. Open a new issue with:
   - Command you ran
   - Full error message
   - OS and Python version (`python --version`)