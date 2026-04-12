"""config command - manage configuration"""

import json
import os
from pathlib import Path

CONFIG_PATH = Path.home() / '.formseal-sync' / 'config.json'

DEFAULT_CONFIG = {
    'provider': '',
    'cloudflare': {'token': '', 'namespace': ''},
    'supabase': {'url': '', 'key': '', 'table': ''},
    'output_folder': ''
}


def load():
    if not CONFIG_PATH.exists():
        return DEFAULT_CONFIG.copy()
    try:
        with open(CONFIG_PATH) as f:
            return json.load(f)
    except:
        return DEFAULT_CONFIG.copy()


def save(cfg):
    CONFIG_PATH.parent.mkdir(parents=True, exist_ok=True)
    with open(CONFIG_PATH, 'w') as f:
        json.dump(cfg, f, indent=2)


def run():
    cfg = load()
    print('Provider:', cfg.get('provider', ''))
    print('Cloudflare Token:', cfg.get('cloudflare', {}).get('token', ''))
    print('Cloudflare Namespace:', cfg.get('cloudflare', {}).get('namespace', ''))
    print('Supabase URL:', cfg.get('supabase', {}).get('url', ''))
    print('Supabase Key:', cfg.get('supabase', {}).get('key', ''))
    print('Supabase Table:', cfg.get('supabase', {}).get('table', ''))
    print('Output Folder:', cfg.get('output_folder', ''))


def set(args):
    if len(args) < 2:
        print('Usage: formseal-sync set <key> <value>')
        return

    key, value = args[0], args[1]
    cfg = load()

    mapping = {
        'provider': ('provider', None),
        'cf-token': ('cloudflare', 'token'),
        'cf-namespace': ('cloudflare', 'namespace'),
        'sb-url': ('supabase', 'url'),
        'sb-key': ('supabase', 'key'),
        'sb-table': ('supabase', 'table'),
        'output': ('output_folder', None),
    }

    if key not in mapping:
        print(f'Unknown key: {key}')
        return

    section, subkey = mapping[key]

    if section == 'output_folder':
        cfg['output_folder'] = value
    elif subkey:
        if section not in cfg:
            cfg[section] = {}
        cfg[section][subkey] = value
    else:
        cfg[section] = value

    save(cfg)
    print('Saved.')