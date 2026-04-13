# cli/providers/supabase/storage/db.py
# Supabase database storage adapter

import json
import urllib.request
import urllib.error
from pathlib import Path

from cli.ui import fail, info


def _load_seen(output_path) -> set:
    p = Path(output_path)
    if not p.exists():
        return set()
    return {line.strip() for line in p.read_text(encoding="utf-8").splitlines() if line.strip()}


def fetch(url, key, table, output_path):
    """Fetch all ciphertexts from a Supabase table."""
    req = urllib.request.Request(
        f"{url}/rest/v1/{table}?select=data",
        headers={
            "Authorization": f"Bearer {key}",
            "apikey": key,
        }
    )

    try:
        with urllib.request.urlopen(req, timeout=30) as resp:
            rows = json.loads(resp.read().decode())
    except urllib.error.HTTPError as e:
        fail(f"HTTP {e.code}: {e.read().decode()}")

    seen = _load_seen(output_path)
    Path(output_path).parent.mkdir(parents=True, exist_ok=True)

    written = 0
    skipped = 0

    with open(output_path, "a", encoding="utf-8") as f:
        for row in rows:
            value = (row.get("data") or "").strip()
            if not value:
                continue
            if value in seen:
                skipped += 1
                continue
            f.write(value + "\n")
            seen.add(value)
            written += 1

    return written, skipped