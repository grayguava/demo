"""fetch command - fetch ciphertexts from provider"""

import sys
import os

sys.path.insert(0, os.path.dirname(os.path.dirname(__file__)))

from commands.config import load


def run():
    cfg = load()

    if not cfg.get('output_folder'):
        print('Error: output folder not set')
        print('Run: formseal-sync set output <path>')
        return

    output_path = os.path.join(cfg['output_folder'], 'ciphertexts.jsonl')

    provider = cfg.get('provider', '')
    print(f'Fetching from {provider}...')

    if provider == 'cloudflare':
        fetch_cloudflare(cfg, output_path)
    elif provider == 'supabase':
        fetch_supabase(cfg, output_path)
    else:
        print('Error: provider not set')
        print('Run: formseal-sync set provider <cloudflare|supabase>')


def fetch_cloudflare(cfg, output_path):
    token = cfg.get('cloudflare', {}).get('token', '')
    namespace = cfg.get('cloudflare', {}).get('namespace', '')

    if not token or not namespace:
        print('Error: Cloudflare token or namespace not set')
        return

    try:
        from providers.cloudflare import kv
        account_id = kv.get_account_id(token)
        written, skipped = kv.fetch(token, namespace, account_id, output_path)
        print(f'Done! {written} new, {skipped} duplicates')
    except Exception as e:
        print(f'Error: {e}')


def fetch_supabase(cfg, output_path):
    url = cfg.get('supabase', {}).get('url', '')
    key = cfg.get('supabase', {}).get('key', '')
    table = cfg.get('supabase', {}).get('table', 'ciphertexts')

    if not url or not key:
        print('Error: Supabase URL or key not set')
        return

    try:
        from providers.supabase import db
        written, skipped = db.fetch(url, key, table, output_path)
        print(f'Done! {written} new, {skipped} duplicates')
    except Exception as e:
        print(f'Error: {e}')