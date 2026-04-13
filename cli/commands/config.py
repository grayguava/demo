# cli/commands/config.py
# Config management commands
# Sensitive data: user must set via env vars manually

import json
from pathlib import Path

from cli.ui import br, fail, ok, info, warn, G, W, D, C, Y, R


CONFIG_DIR  = Path.home() / ".formsealdaemon"
CONFIG_FILE = CONFIG_DIR / "config.json"

# Env var names for sensitive data
ENV_KEYS = {
    "cloudflare": {
        "token": "FSYNC_CF_TOKEN",
    },
    "supabase": {
        "key": "FSYNC_SU_KEY",
    },
}

# Config keys (non-sensitive)
VALID_PROVIDERS = {
    "cloudflare": {
        "namespace": "KV namespace ID",
    },
    "supabase": {
        "url": "Supabase project URL",
    },
}

# Global config keys
VALID_GLOBAL = {
    "output_folder": "Output folder path",
    "sync_interval": "Sync interval in minutes (default: 15)",
}


def load_config():
    if CONFIG_FILE.exists():
        return json.loads(CONFIG_FILE.read_text())
    return {}


def save_config(cfg):
    CONFIG_DIR.mkdir(parents=True, exist_ok=True)
    CONFIG_FILE.write_text(json.dumps(cfg, indent=2))


def get_token(provider):
    """Get token from env var."""
    import os
    env_key = ENV_KEYS.get(provider, {}).get("token") or ENV_KEYS.get(provider, {}).get("key")
    if env_key:
        val = os.environ.get(env_key)
        return val
    return None


def run_set(args):
    if len(args) < 2:
        fail("Usage: fsync set <key> <value>\n\nExample:\n  fsync set provider cloudflare\n  fsync set namespace <id>\n  fsync set output_folder <path>\n  fsync set sync_interval 15\n\nOr set env vars manually:\n  $env:FSYNC_CF_TOKEN=\"cfut_...\"")

    key = args[0]
    value = " ".join(args[1:])

    cfg = load_config()

    # Handle provider switch
    if key == "provider":
        if value not in VALID_PROVIDERS:
            fail(f"Unknown provider: {value}\nValid: {', '.join(VALID_PROVIDERS.keys())}")

        old_provider = cfg.get("provider")
        if old_provider and old_provider != value:
            for k in list(cfg.keys()):
                if k.startswith(old_provider + "."):
                    del cfg[k]

        cfg["provider"] = value
        save_config(cfg)

        br()
        ok(f"Provider set to {value}")
        br()
        return

    # Handle global keys (not provider-specific)
    if key in VALID_GLOBAL:
        if key == "output_folder":
            cfg["output_folder"] = value
        elif key == "sync_interval":
            try:
                cfg["sync_interval"] = int(value)
            except ValueError:
                fail("sync_interval must be a number")
        save_config(cfg)
        br()
        ok(f"Set {key} = {value}")
        br()
        return

    # Validate key - check if it's a sensitive token/key
    provider = cfg.get("provider")
    if not provider:
        fail("No provider set. Run: fsync set provider <cloudflare|supabase>")

    # Check if this key is sensitive (should be env var)
    env_key = ENV_KEYS.get(provider, {}).get(key)
    if env_key:
        br()
        info(f"Set {key} via env var instead:")
        print(f"  {W}Windows:{R} $env:{env_key}=\"<token>\"")
        print(f"  {W}Linux/Mac:{R} export {env_key}=\"<token>\"")
        br()
        warn("Tokens are not stored in config for security.")
        br()
        return

    # Non-sensitive key - save to config
    valid_key = f"{provider}.{key}"
    if key not in VALID_PROVIDERS.get(provider, {}):
        fail(f"Unknown key: {key}\n\nValid keys for {provider}:\n" +
             "\n".join(f"  {k}: {v}" for k, v in VALID_PROVIDERS[provider].items()) +
             "\n\nGlobal keys:\n" +
             "\n".join(f"  {k}: {v}" for k, v in VALID_GLOBAL.items()))

    cfg[valid_key] = value
    save_config(cfg)

    br()
    ok(f"Set {key} = {value}")
    br()


def run_status():
    import os

    cfg = load_config()

    br()
    print(f"{C} \u250c\u2500 {R}{W}formseal-sync{R}  {G}status{R}")
    print(G + " " + "\u2500" * 52 + R)
    br()

    provider = cfg.get("provider")
    if not provider:
        warn("No provider configured. Run: fsync configure quick")
        br()
        return

    print(f"  {W}provider:{R}   {W}{provider}{R}")
    br()

    # Provider-specific config (account/storage details)
    pcfg = {k.replace(f"{provider}.", ""): v for k, v in cfg.items() if k.startswith(f"{provider}.")}

    if provider == "cloudflare":
        namespace = pcfg.get("namespace")
        token = get_token(provider)
        
        print(f"  namespace:  {W if namespace else D}{namespace or '(not set)'}{R}")
        
        import os
        env_key = ENV_KEYS.get(provider, {}).get("token")
        env_set = env_key in os.environ and os.environ.get(env_key)
        
        if env_set and token:
            # Try to get and show partial account ID
            try:
                from cli.providers.cloudflare.account import get_account_id
                account_id = get_account_id(token)
                partial = account_id[:12] + "..." if len(account_id) > 12 else account_id
                print(f"  account:    {G}{partial}{R}")
            except Exception as e:
                print(f"  account:    {R}(auth error){R}")
                print(f"  {D}Token may need 'account:read' scope{R}")
            print(f"  token:      {W}****{R}")
        elif env_set:
            print(f"  token:      {D}(value empty){R}")
            print(f"  account:    {D}(set token to detect){R}")
        else:
            print(f"  token:      {D}(not set - env var {env_key}){R}")
            print(f"  account:    {D}(set token to detect){R}")
        
        print(f"  storage:    kv")
    else:  # supabase
        url = pcfg.get("url")
        token = get_token(provider)
        
        print(f"  url:        {W if url else D}{url or '(not set)'}{R}")
        print(f"  token:      {W if token else D}{'****' if token else '(not set)'}{R}")
        print(f"  storage:    db")

    br()

    # Global config
    output_folder = cfg.get("output_folder")
    print(f"  output_folder:  {W if output_folder else D}{output_folder or '(not set)'}{R}")

    sync_interval = cfg.get("sync_interval")
    print(f"  sync_interval:   {W if sync_interval else D}{sync_interval or '(not set)'}{R}")

    if cfg.get("last_sync"):
        print(f"  last_sync:      {W}{cfg['last_sync']}{R}")

    br()


def run_config_show():
    import os

    cfg = load_config()

    br()
    print(f"{C} \u250c\u2500 {R}{W}formseal-sync{R}  {G}config{R}")
    print(G + " " + "\u2500" * 52 + R)
    br()

    provider = cfg.get("provider")

    if not cfg:
        info("No config. Run: fsync set provider <cloudflare|supabase>")
    else:
        for k, v in cfg.items():
            print(f"  {D}{k}:{R}  {W}{v}{R}")

    # Show env var status
    if provider:
        token = get_token(provider)
        env_key = ENV_KEYS.get(provider, {}).get("token") or ENV_KEYS.get(provider, {}).get("key")
        print(f"  {D}{env_key}:{R}  {'****' if token else 'not set'}")

    br()


def run_logout():
    if CONFIG_FILE.exists():
        CONFIG_FILE.unlink()
        br()
        ok("Config cleared. Manually unset env vars:")
        print("  $env:FSYNC_CF_TOKEN=$null")
        br()
    else:
        br()
        info("No config to clear.")
        br()