# Supabase

PostgreSQL database backend for formseal-fetch.

## Status: coming soon

Supabase support is planned but not yet implemented. This document describes the expected functionality.

## Expected features

- Connect to a Supabase PostgreSQL database
- Fetch encrypted form submissions stored in a table
- Support for custom table names and column mappings

## Planned setup

When available, expected connection syntax:

```bash
fsf connect provider:supabase host:<host> database:<name> user:<username> password:<pass>
```

Or using Supabase connection string:

```bash
fsf connect provider:supabase connection:<connection_string>
```

## Expected requirements

- Supabase project (PostgreSQL)
- Table containing ciphertext data
- Read-only database user (recommended)

## Timeline

Supabase backend will be added in a future release. Check the [GitHub repository](https://github.com/formseal/formseal-fetch) for updates.

## Alternative

For now, only [Cloudflare KV](./cloudflare-kv.md) is supported.