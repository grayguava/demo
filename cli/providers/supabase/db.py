"""Supabase DB provider"""

import os
import requests


def fetch(url, key, table, output_path):
    fetch_url = f'{url}/rest/v1/{table}?select=data'

    resp = requests.get(
        fetch_url,
        headers={
            'Authorization': f'Bearer {key}',
            'apikey': key
        }
    )

    if resp.status_code != 200:
        raise Exception(f'HTTP {resp.status_code}: {resp.text}')

    rows = resp.json()

    os.makedirs(os.path.dirname(output_path), exist_ok=True)

    seen = set()
    if os.path.exists(output_path):
        with open(output_path) as f:
            seen = set(line.strip() for line in f if line.strip())

    written = 0
    skipped = 0

    with open(output_path, 'a') as f:
        for row in rows:
            value = row.get('data', '').strip()
            if not value:
                continue

            if value in seen:
                skipped += 1
                continue

            f.write(value + '\n')
            seen.add(value)
            written += 1

    return written, skipped