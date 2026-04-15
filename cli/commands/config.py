# cli/commands/config.py
# Config management commands

import json
import sys
from pathlib import Path

from pathlib import Path

from cli.ui import br, fail, ok, info, warn, G, W, D, C, Y, R
from cli.security import tokens


def _load_version():
    p = Path(__file__).parent.parent.parent / "version.txt"
    if p.exists():
        return p.read_text().strip()
    return "dev"


VERSION = _load_version()


CONFIG_DIR = Path.home() / ".config" / "formseal-fetch"
CONFIG_FILE = CONFIG_DIR / "config.json"


VALID_PROVIDERS = {
    "cloudflare": {},
}

VALID_GLOBAL = {
    "output_folder": "Output folder path",
}


def load_config():
    CONFIG_DIR.mkdir(parents=True, exist_ok=True)
    if CONFIG_FILE.exists():
        return json.loads(CONFIG_FILE.read_text())
    return {}


def save_config(cfg):
    CONFIG_DIR.mkdir(parents=True, exist_ok=True)
    CONFIG_FILE.write_text(json.dumps(cfg, indent=2))


def get_token(provider: str):
    return tokens.load_token(provider)


def get_namespace(provider: str):
    return tokens.load_namespace(provider)


def run_set(args):
    if len(args) < 2:
        fail("Usage: fsf set <key> <value>")

    key = args[0]
    value = " ".join(args[1:])

    cfg = load_config()

    if key == "provider":
        if value not in VALID_PROVIDERS:
            fail(f"Unknown provider: {value}\nValid: {', '.join(VALID_PROVIDERS.keys())}")
        cfg["provider"] = value
        save_config(cfg)
        br()
        ok(f"Provider set to {value}")
        br()
        return

    if key in VALID_GLOBAL:
        cfg[key] = value
        save_config(cfg)
        br()
        ok(f"Set {key} = {value}")
        br()
        return

    fail(f"Unknown key: {key}")


def run_status():
    cfg = load_config()

    br()
    print(f"{C} \u250c\u2500 {R}{W}formseal-fetch{R}  {Y}v{VERSION}{R}")
    print(G + " " + "\u2500" * 52 + R)
    br()

    print(f"  {D}Configuration Status:{R}")
    br()

    provider = cfg.get("provider")
    if not provider:
        warn("No provider configured. Run: fsf connect provider:cloudflare")
        br()
        return

    def row(label, value, color=W):
        print(f"  {D}{label:<26}{R}{color}{value}{R}")

    row("Provider:", "Cloudflare")

    if provider == "cloudflare":
        namespace = get_namespace(provider)
        token = get_token(provider)

        row("KV Namespace ID:", namespace or "(not set)", W if namespace else D)
        if namespace:
            row("NS-ID Location:", tokens.namespace_location(provider), G)

        if token:
            try:
                from cli.providers.cloudflare.account import get_account_id
                account_id = get_account_id(token)
                partial = account_id[:12] + "..." if len(account_id) > 12 else account_id
                row("Account ID:", partial, G)
            except Exception:
                row("Account ID:", "(auth error)", R)
            row("API Token:", "****")
            row("Token Location:", tokens.token_location(provider), G)
        else:
            row("API Token:", "(not set)", D)
            row("Token Location:", "Not set", D)

        br()
        row("Server Storage Type:", "Key-Value")

    else:
        fail(f"Unknown provider: {provider}")

    output_folder = cfg.get("output_folder")
    row("Output Folder:", output_folder or "(not set)", W if output_folder else D)

    br()

    br()


def run_config_show():
    cfg = load_config()

    br()
    print(f"{C} \u250c\u2500 {R}{W}formseal-fetch{R}  {G}config{R}")
    print(G + " " + "\u2500" * 52 + R)
    br()

    provider = cfg.get("provider")

    if not cfg:
        info("No config. Run: fsf connect provider:cloudflare")
    else:
        for k, v in cfg.items():
            print(f"  {D}{k}:{R}  {W}{v}{R}")

    if provider:
        token = get_token(provider)
        print(f"  {D}token:{R}  {'****' if token else 'not set'}")

    br()


def run_disconnect():
    br()
    print(f"  {Y}This will delete all config and credentials.{R}")
    print(f"  Downloaded ciphertexts will NOT be affected.")
    br()
    sys.stdout.write(f"  Continue? [y/N]: ")
    sys.stdout.flush()
    confirm = input().strip().lower()
    
    if confirm != "y":
        br()
        info("Cancelled.")
        br()
        return
    
    cfg = load_config()
    provider = cfg.get("provider")
    
    if CONFIG_FILE.exists():
        CONFIG_FILE.unlink()
    
    if provider:
        tokens.clear_all(provider)
    
    br()
    ok("Disconnected. All config and credentials cleared.")
    br()


# Alias for backward compatibility
run_logout = run_disconnect