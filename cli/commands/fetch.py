# cli/commands/fetch.py
# Fetch ciphertexts from configured provider

import argparse

from cli.ui import br, fail, ok, info, G, W, D, C, R
from cli.commands.config import load_config
from cli.security import tokens


def run(args):
    parser = argparse.ArgumentParser(prog="fsf fetch")
    parser.add_argument("--output", default=None)
    parsed = parser.parse_args(args)

    cfg = load_config()
    provider = cfg.get("provider")
    output_folder = cfg.get("output_folder", "data")

    if not provider:
        fail("No provider set. Run: fsf connect provider:cloudflare")

    if parsed.output:
        output_path = parsed.output
    else:
        output_path = f"{output_folder}/ciphertexts.jsonl"

    token = tokens.load_token(provider)
    if not token:
        fail("No token. Run: fsf connect provider:cloudflare to set token")

    if provider == "cloudflare":
        namespace = tokens.load_namespace(provider)
        if not namespace:
            fail("No namespace. Run: fsf connect provider:cloudflare to set namespace")

        try:
            from cli.providers.cloudflare.account import get_account_id
            account_id = get_account_id(token)
        except Exception as e:
            fail(f"Auth failed: {e}")

        br()
        provider_display = f"{provider.capitalize()} KV"
        print(f"{C} \u250c\u2500 {R}{W}formseal-fetch{R}  {G}{provider_display}{R}")
        print(G + " " + "\u2500" * 52 + R)
        br()

        account_trunc = account_id[:8] + "***" if len(account_id) > 8 else account_id
        namespace_trunc = namespace[:8] + "***" if len(namespace) > 8 else namespace

        print(f"{D}>> Account:{R}  {D}{account_trunc}{R}")
        print(f"{D}>> Namespace:{R} {D}{namespace_trunc}{R}")
        br()

        try:
            from cli.providers.cloudflare.storage import fetch as storage_fetch
            written, skipped = storage_fetch(
                namespace=namespace,
                account_id=account_id,
                token=token,
                output_path=output_path,
            )
        except Exception as e:
            fail(f"Fetch failed: {e}")

        br()
        ok(f"{written} new ciphertexts saved → {output_path}")
        if skipped:
            print(f"  {D}({skipped} duplicates skipped){R}")
        br()
    else:
        fail(f"Unknown provider: {provider}")