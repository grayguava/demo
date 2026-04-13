# cli/commands/setup.py
# Interactive setup commands

import os
import sys

from cli.ui import br, fail, ok, info, warn, G, W, D, C, Y, O, R
from cli.commands.config import load_config, save_config, ENV_KEYS


# Schema: what fields to ask for each provider/storage combo
SETUP_SCHEMA = {
    "cloudflare": {
        "kv": {
            "config_fields": [
                {"key": "namespace", "prompt": "KV Namespace ID", "required": True},
            ],
            "env_key": "FSYNC_CF_TOKEN"
        }
    },
    "supabase": {
        "db": {
            "config_fields": [
                {"key": "url", "prompt": "Project URL", "required": True},
                {"key": "table", "prompt": "Table name", "required": False, "default": "ciphertexts"},
            ],
            "env_key": "FSYNC_SU_KEY"
        }
    }
}


def run(args):
    if not args:
        fail("Usage: fsync setup <quick|reset|sync>")

    command = args[0]

    match command:
        case "quick":
            _setup_quick()
        case "reset":
            _setup_reset()
        case "sync":
            _setup_sync()
        case _:
            fail(f"Unknown setup command: {command}")


def _setup_quick():
    """Interactive quick setup."""
    print()
    print(f"{C} \u250c\u2500 {R}{W}formseal-sync{R}  {G}quick setup{R}")
    print(G + " " + "\u2500" * 52 + R)
    print()

    cfg = load_config()

    # 1. Provider
    print(f"  {D}See supported: fsync providers{R}")
    sys.stdout.write(f"  1. Provider: ")
    sys.stdout.flush()
    provider = input().strip().lower()
    if not provider:
        fail("Provider required. Run: fsync providers to see available.")
    if provider not in SETUP_SCHEMA:
        fail(f"Invalid provider. Available: {', '.join(SETUP_SCHEMA.keys())}")
    cfg["provider"] = provider
    print()

    # 2. Storage Type
    storage_options = list(SETUP_SCHEMA[provider].keys())
    print(f"  {D}Available: {', '.join(storage_options)}{R}")
    sys.stdout.write(f"  2. Storage Type: ")
    sys.stdout.flush()
    storage = input().strip().lower()
    if not storage:
        storage = storage_options[0]
    if storage not in SETUP_SCHEMA[provider]:
        fail(f"Invalid storage type for {provider}. Available: {', '.join(storage_options)}")
    
    # Save storage type
    cfg[f"{provider}.storage"] = storage
    print()

    # 3. Provider-specific config fields
    fields = SETUP_SCHEMA[provider][storage]["config_fields"]
    env_key = SETUP_SCHEMA[provider][storage]["env_key"]
    
    # Show env var reminder
    env_key = SETUP_SCHEMA[provider][storage]["env_key"]
    print(f"  {D}Sensitive data (NOT stored in config):{R}")
    print(f"    Token/Key - set via env var {env_key}")
    print(f"    Run: fsync help --vars")
    print()
    print(f"  {D}Auto-detected:{R}")
    if provider == "cloudflare":
        print(f"    Account ID - detected from token")
    print()

    for i, field in enumerate(fields):
        required = field.get("required", False)
        default = field.get("default", "")
        
        sys.stdout.write(f"  {field['prompt']}: ")
        sys.stdout.flush()
        val = input().strip()
        
        if not val and default:
            val = default
        
        if required and not val:
            fail(f"{field['prompt']} is required")
        
        if val:
            cfg[f"{provider}.{field['key']}"] = val
        print()

    # 4. Output folder
    sys.stdout.write(f"  Output Folder: ")
    sys.stdout.flush()
    output_folder = input().strip()
    if not output_folder:
        output_folder = "data"
    cfg["output_folder"] = output_folder
    print()

    save_config(cfg)

    print(f"{G} \u2713{R} Saved!")
    print()

    # Remind about env vars
    print(f"  Don't forget to set sensitive data:")
    print(f"    Run: fsync help --vars")
    print()


def _setup_reset():
    """Reset configuration."""
    from cli.commands.config import CONFIG_FILE
    
    if CONFIG_FILE.exists():
        CONFIG_FILE.unlink()
    
    print()
    ok("Configuration reset")
    print()
    print(f"  Run {W}fsync setup quick{R} to set up again")
    print()


def _setup_sync():
    """Configure sync settings."""
    print()
    print(f"{C} \u250c\u2500 {R}{W}formseal-sync{R}  {G}sync setup{R}")
    print(G + " " + "\u2500" * 52 + R)
    print()

    cfg = load_config()

    current_interval = cfg.get("sync_interval", 15)
    print(f"  {D}Current interval:{R}  {W}{current_interval} min{R}")
    print()

    sys.stdout.write(f"  Sync Interval (minutes): ")
    sys.stdout.flush()
    interval = input().strip()

    if interval:
        try:
            cfg["sync_interval"] = int(interval)
            save_config(cfg)
            print()
            ok(f"Sync interval set to {interval} minutes")
        except ValueError:
            fail("Invalid number")
    else:
        print()
        info("No change (press Enter to keep current)")

    print()


def _get_os():
    """Detect OS for help message."""
    if sys.platform == "win32":
        return "windows"
    elif sys.platform == "darwin":
        return "macos"
    else:
        return "linux"