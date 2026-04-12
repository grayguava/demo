"""Cloudflare KV provider"""

import os
import requests


def get_account_id(token):
    url = 'https://api.cloudflare.com/client/v4/user'
    resp = requests.get(url, headers={'Authorization': f'Bearer {token}'})
    data = resp.json()
    if not data.get('success'):
        raise Exception(data['errors'][0]['message'])
    accounts = data['result']['accounts']
    if not accounts:
        raise Exception('No accounts found')
    return accounts[0]['id']


def fetch(token, namespace_id, account_id, output_path):
    base = f'https://api.cloudflare.com/client/v4/accounts/{account_id}/storage/kv/namespaces/{namespace_id}'

    all_keys = []
    cursor = None

    while True:
        url = f'{base}/keys'
        if cursor:
            url += f'?cursor={cursor}'

        resp = requests.get(url, headers={'Authorization': f'Bearer {token}'})
        data = resp.json()

        if not data.get('success'):
            raise Exception(data['errors'][0]['message'])

        all_keys.extend(k['name'] for k in data['result'])

        cursor = data.get('result_info', {}).get('cursor')
        if not cursor:
            break

    if not all_keys:
        return 0, 0

    os.makedirs(os.path.dirname(output_path), exist_ok=True)

    seen = set()
    if os.path.exists(output_path):
        with open(output_path) as f:
            seen = set(line.strip() for line in f if line.strip())

    written = 0
    skipped = 0

    with open(output_path, 'a') as f:
        for key in all_keys:
            url = f'{base}/values/{requests.utils.quote(key)}'
            resp = requests.get(url, headers={'Authorization': f'Bearer {token}'})

            if resp.status_code != 200:
                continue

            value = resp.text.strip()
            if not value:
                continue

            if value in seen:
                skipped += 1
                continue

            f.write(value + '\n')
            seen.add(value)
            written += 1

    return written, skipped