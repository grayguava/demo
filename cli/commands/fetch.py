# cli/commands/fetch.py
# Fetch ciphertexts from configured provider

from cli.ui import br, fail, ok, info, G, W, D, C, R
from cli.commands.config import load_config, get_token


def run(args):
    import argparse
    parser = argparse.ArgumentParser(prog="fsync fetch")
    parser.add_argument("--output", default=None)
    parsed = parser.parse_args(args)

    cfg = load_config()
    provider = cfg.get("provider")
    output_folder = cfg.get("output_folder", "data")

    # Use --output if provided, otherwise use config
    if parsed.output:
        output_path = parsed.output
    else:
        output_path = f"{output_folder}/ciphertexts.jsonl"

    if not provider:
        fail("No provider set. Run: fsync set provider <cloudflare|supabase>")

    # Build provider config from stored keys (removes the prefix)
    pcfg = {k.replace(f"{provider}.", ""): v for k, v in cfg.items() if k.startswith(f"{provider}.")}

    # Get token from env var
    token = get_token(provider)

    # Import provider modules
    try:
        if provider == "cloudflare":
            from cli.providers.cloudflare.storage import fetch as storage_fetch
            from cli.providers.cloudflare.account import get_account_id
        elif provider == "supabase":
            from cli.providers.supabase.storage import fetch as storage_fetch
            from cli.providers.supabase.account import validate_token
        else:
            fail(f"Unknown provider: {provider}")
    except (ImportError, AttributeError) as e:
        fail(f"Provider {provider} not configured: {e}")

    br()
    print(f"{C} \u250c\u2500 {R}{W}formseal-sync{R}  {G}{provider} fetch{R}")
    print(G + " " + "\u2500" * 52 + R)
    br()

    if provider == "cloudflare":
        namespace = pcfg.get("namespace")

        if not token:
            fail("No token. Set via env var FSYNC_CF_TOKEN")
        if not namespace:
            fail("No namespace. Run: fsync setup quick to configure")

        # Auto-detect account_id
        try:
            account_id = get_account_id(token)
            info(f"Account: {account_id}")
        except Exception as e:
            fail(f"Auth failed: {e}")

        info(f"Namespace: {namespace}")
        br()

        written, skipped = storage_fetch(
            namespace    = namespace,
            account_id  = account_id,
            token       = token,
            output_path = output_path,
        )
    else:
        # Pass token as part of pcfg
        pcfg["_token"] = token
        written, skipped = storage_fetch(pcfg, output_path)

    br()
    ok(f"{written} new ciphertexts saved → {output_path}")
    if skipped:
        info(f"{skipped} duplicates skipped")
    br()