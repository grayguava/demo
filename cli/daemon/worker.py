# cli/daemon/worker.py
# Entry point for detached sync process

import os
import sys
import os
import time
from pathlib import Path
from datetime import datetime

SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
PROJECT_DIR = os.path.dirname(os.path.dirname(SCRIPT_DIR))
sys.path.insert(0, PROJECT_DIR)

from cli.commands.config import load_config, get_token


PID_FILE = Path.home() / ".formsealdaemon" / "sync.pid"
LOG_FILE = Path.home() / ".formsealdaemon" / "sync.log"


def _log(msg):
    timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    LOG_FILE.parent.mkdir(parents=True, exist_ok=True)
    with open(LOG_FILE, "a") as f:
        f.write(f"[{timestamp}] {msg}\n")


def _run_sync_once(output_folder):
    cfg = load_config()
    provider = cfg.get("provider")

    if not provider:
        _log("No provider configured")
        return 0, 0

    pcfg = {k.replace(f"{provider}.", ""): v for k, v in cfg.items() if k.startswith(f"{provider}.")}
    token = get_token(provider)

    if not token:
        _log(f"No token for {provider}")
        return 0, 0

    try:
        if provider == "cloudflare":
            from cli.providers.cloudflare.storage import fetch as storage_fetch
            from cli.providers.cloudflare.account import get_account_id
        elif provider == "supabase":
            from cli.providers.supabase.storage import fetch as storage_fetch
            from cli.providers.supabase.account import validate_token
        else:
            _log(f"Unknown provider: {provider}")
            return 0, 0
    except (ImportError, AttributeError) as e:
        _log(f"Provider {provider} not configured: {e}")
        return 0, 0

    output_path = os.path.join(output_folder, "ciphertexts.jsonl")

    if provider == "cloudflare":
        namespace = pcfg.get("namespace")
        if not namespace:
            _log("No namespace set")
            return 0, 0

        try:
            account_id = get_account_id(token)
        except Exception as e:
            _log(f"Auth failed: {e}")
            return 0, 0

        written, skipped = storage_fetch(
            namespace=namespace,
            account_id=account_id,
            token=token,
            output_path=output_path,
        )
    else:
        pcfg["_token"] = token
        written, skipped = storage_fetch(pcfg, output_path)

    cfg["last_sync"] = datetime.now().isoformat()
    from cli.commands.config import save_config
    save_config(cfg)

    _log(f"Sync complete: {written} new, {skipped} skipped")
    return written, skipped


def run():
    _log("Background sync started")

    PID_FILE.write_text(str(os.getpid()))

    while True:
        try:
            cfg = load_config()
            output_folder = cfg.get("output_folder", "data")
            interval = cfg.get("sync_interval", 15)
            _run_sync_once(output_folder)
        except Exception as e:
            _log(f"Sync error: {e}")

        for _ in range(interval * 60):
            time.sleep(1)

    _log("Background sync stopped")


if __name__ == "__main__":
    run()