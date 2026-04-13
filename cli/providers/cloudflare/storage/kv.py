# cli/providers/cloudflare/storage/kv.py
# Cloudflare KV storage adapter

import json
import urllib.request
import urllib.error
import urllib.parse
from pathlib import Path

from cli.ui import fail, info


def _get(url, token):
    req = urllib.request.Request(url, headers={"Authorization": f"Bearer {token}"})
    try:
        with urllib.request.urlopen(req, timeout=30) as resp:
            return json.loads(resp.read().decode())
    except urllib.error.HTTPError as e:
        fail(f"HTTP {e.code}: {e.read().decode()}")


def _get_raw(url, token):
    req = urllib.request.Request(url, headers={"Authorization": f"Bearer {token}"})
    try:
        with urllib.request.urlopen(req, timeout=30) as resp:
            return resp.read().decode("utf-8").strip()
    except urllib.error.HTTPError as e:
        fail(f"HTTP {e.code}: {e.read().decode()}")


def _load_seen(output_path) -> set:
    p = Path(output_path)
    if not p.exists():
        return set()
    return {line.strip() for line in p.read_text(encoding="utf-8").splitlines() if line.strip()}


def fetch(namespace, account_id, token, output_path):
    """Fetch all values from a KV namespace."""
    base = f"https://api.cloudflare.com/client/v4/accounts/{account_id}/storage/kv/namespaces/{namespace}"

    # List all keys (paginated)
    all_keys = []
    cursor = None

    while True:
        url = f"{base}/keys" + (f"?cursor={cursor}" if cursor else "")
        data = _get(url, token)
        if not data.get("success"):
            fail(f"API error: {data.get('errors')}")
        all_keys.extend(k["name"] for k in data.get("result", []))
        cursor = data.get("result_info", {}).get("cursor")
        if not cursor:
            break

    info(f"Found {len(all_keys)} keys")

    if not all_keys:
        info("No data to fetch.")
        return 0, 0

    # Deduplicate against existing
    seen = _load_seen(output_path)
    Path(output_path).parent.mkdir(parents=True, exist_ok=True)

    written = 0
    skipped = 0

    with open(output_path, "a", encoding="utf-8") as f:
        for key in all_keys:
            value = _get_raw(f"{base}/values/{urllib.parse.quote(key, safe='')}", token)
            if not value:
                continue
            if value in seen:
                skipped += 1
                continue
            f.write(value + "\n")
            seen.add(value)
            written += 1

    return written, skipped