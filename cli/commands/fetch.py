# Fetch ciphertexts

import argparse

from cli.ui import br, fail, ok, info, G, W, D, C, R
from cli.commands.config import load_config
from cli.security import tokens
from cli.providers import get_provider


def run(args):
    parser = argparse.ArgumentParser(prog="fsf fetch")
    parser.add_argument("--output", default=None)
    parsed = parser.parse_args(args)

    cfg = load_config()
    provider_name = cfg.get("provider")
    output_folder = cfg.get("output_folder", "data")

    if not provider_name:
        fail("No provider set. Run: fsf connect provider:<name>")

    if parsed.output:
        output_path = parsed.output
    else:
        output_path = f"{output_folder}/ciphertexts.jsonl"

    token = tokens.load_token(provider_name)
    if not token:
        fail("No token. Run: fsf connect provider:<name> to set token")

    provider = get_provider(provider_name)
    if not provider:
        fail(f"Unknown provider: {provider_name}")

    provider_config = {"namespace": tokens.load_namespace(provider_name)}

    account_id = None
    try:
        account_id = provider.authenticate(token)
    except Exception as e:
        fail(f"Auth failed: {e}")

    br()
    print(f"{C} \u250c\u2500 {R}{W}formseal-fetch{R}  {G}{provider.display_name}{R}")
    print(G + " " + "\u2500" * 52 + R)
    br()

    if account_id:
        account_trunc = account_id[:8] + "***" if len(account_id) > 8 else account_id
        print(f"{D}>> Account:{R}  {D}{account_trunc}{R}")

    for field in provider.get_config_fields():
        key = field["key"]
        value = provider_config.get(key)
        if value:
            trunc = value[:8] + "***" if len(value) > 8 else value
            print(f"{D}>> {field.get('prompt', key).replace(':','')}:{R}  {D}{trunc}{R}")

    br()

    try:
        written, skipped = provider.fetch(token, provider_config, output_path)
    except Exception as e:
        fail(f"Fetch failed: {e}")

    br()
    ok(f"{written} new ciphertexts saved → {output_path}")
    if skipped:
        print(f"  {D}({skipped} duplicates skipped){R}")
    br()