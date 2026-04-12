"""status command - show stats"""

import os
import sys

sys.path.insert(0, os.path.dirname(os.path.dirname(__file__)))

from commands.config import load


def run():
    cfg = load()
    folder = cfg.get('output_folder', '')

    if not folder:
        print('Output folder not set')
        return

    path = os.path.join(folder, 'ciphertexts.jsonl')

    if not os.path.exists(path):
        print('0 ciphertexts')
        return

    with open(path) as f:
        count = sum(1 for line in f if line.strip())

    print(f'{count} ciphertexts stored')