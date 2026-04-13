# cli/daemon/sync.py
# Background sync daemon

import os
import sys
import subprocess
from pathlib import Path

from cli.ui import br, fail, ok, info, warn, G, W, D, C, Y, O, R
from cli.commands.config import load_config


PID_FILE = Path.home() / ".formsealdaemon" / "sync.pid"
LOG_FILE = Path.home() / ".formsealdaemon" / "sync.log"


def run(args):
    if not args:
        fail("Usage: fsync sync <start|stop|status|run>")

    command = args[0]

    match command:
        case "start":
            _sync_start()
        case "stop":
            _sync_stop()
        case "status":
            _sync_status()
        case "run":
            _sync_run()
        case _:
            fail(f"Unknown sync command: {command}\nValid: start, stop, status, run")


def _sync_start():
    existing_pid = None
    
    if PID_FILE.exists():
        try:
            with open(PID_FILE) as f:
                existing_pid = int(f.read().strip())
            if os.name == "nt":
                result = os.popen(f'tasklist /FI "PID eq {existing_pid}" /NH').read()
                if str(existing_pid) in result:
                    fail(f"Sync already running (PID: {existing_pid}). Run: fsync sync stop first.")
            PID_FILE.unlink()
        except:
            PID_FILE.unlink()

    cfg = load_config()
    if not cfg.get("provider"):
        fail("No provider set. Run: fsync set provider <cloudflare|supabase>")

    if not cfg.get("output_folder"):
        fail("No output folder set. Run: fsync set output_folder <path>")

    interval = cfg.get("sync_interval", 15)

    worker_path = os.path.join(os.path.dirname(os.path.abspath(__file__)), "worker.py")

    if os.name == "nt":
        CREATE_NEW_PROCESS_GROUP = 0x00000200
        DETACHED_PROCESS = 0x00000008
        startupinfo = subprocess.STARTUPINFO()
        startupinfo.dwFlags |= subprocess.STARTF_USESHOWWINDOW
        startupinfo.wShowWindow = subprocess.SW_HIDE
        subprocess.Popen(
            [sys.executable, worker_path],
            creationflags=DETACHED_PROCESS | CREATE_NEW_PROCESS_GROUP,
            stdout=subprocess.DEVNULL,
            stderr=subprocess.DEVNULL,
            startupinfo=startupinfo,
        )
    else:
        subprocess.Popen(
            [sys.executable, worker_path],
            stdout=subprocess.DEVNULL,
            stderr=subprocess.DEVNULL,
            start_new_session=True,
        )

    br()
    ok(f"Background sync started (interval: {interval} min)")
    info(f"Log: {LOG_FILE}")
    br()


def _sync_stop():
    if not PID_FILE.exists():
        br()
        info("Sync not running")
        br()
        return

    try:
        with open(PID_FILE) as f:
            pid = int(f.read().strip())

        if os.name == "nt":
            os.system(f"taskkill /PID {pid} /F")

        PID_FILE.unlink()

        br()
        ok("Background sync stopped")
        br()
    except Exception as e:
        fail(f"Failed to stop sync: {e}")


def _sync_status():
    br()
    print(f"{C} \u250c\u2500 {R}{W}formseal-sync{R}  {G}sync status{R}")
    print(G + " " + "\u2500" * 52 + R)
    br()

    if PID_FILE.exists():
        try:
            with open(PID_FILE) as f:
                pid = int(f.read().strip())
            print(f"  {G}Running:{R}  {W}Yes{R} (PID: {pid})")
        except:
            print(f"  {G}Running:{R}  {W}Unknown{R}")
    else:
        print(f"  {G}Running:{R}  {D}No{R}")

    cfg = load_config()
    interval = cfg.get("sync_interval", 15)
    print(f"  {D}interval:{R}  {W}{interval} min{R}")

    last_sync = cfg.get("last_sync")
    if last_sync:
        print(f"  {D}last_sync:{R}  {W}{last_sync}{R}")
    else:
        print(f"  {D}last_sync:{R}  {D}Never{R}")

    if LOG_FILE.exists():
        br()
        info("Recent log:")
        with open(LOG_FILE) as f:
            lines = f.readlines()
            for line in lines[-5:]:
                print(f"  {D}{line.strip()}{R}")

    br()


def _run_sync_once(output_folder):
    cfg = load_config()
    provider = cfg.get("provider")

    if not provider:
        fail("No provider set. Run: fsync set provider <cloudflare|supabase>")

    pcfg = {k.replace(f"{provider}.", ""): v for k, v in cfg.items() if k.startswith(f"{provider}.")}
    from cli.commands.config import get_token
    token = get_token(provider)

    if not token:
        fail("No token. Set via env var FSYNC_CF_TOKEN or FSYNC_SU_KEY")

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

    if provider == "cloudflare":
        namespace = pcfg.get("namespace")
        if not namespace:
            fail("No namespace set. Run: fsync setup quick")

        try:
            account_id = get_account_id(token)
        except Exception as e:
            fail(f"Auth failed: {e}")

        written, skipped = storage_fetch(
            namespace=namespace,
            account_id=account_id,
            token=token,
            output_path=os.path.join(output_folder, "ciphertexts.jsonl"),
        )
    else:
        pcfg["_token"] = token
        written, skipped = storage_fetch(pcfg, os.path.join(output_folder, "ciphertexts.jsonl"))

    from cli.commands.config import save_config
    from datetime import datetime
    cfg["last_sync"] = datetime.now().isoformat()
    save_config(cfg)

    return written, skipped


def _sync_run():
    cfg = load_config()
    if not cfg.get("output_folder"):
        fail("No output folder set. Run: fsync set output_folder <path>")

    br()
    print(f"{C} \u250c\u2500 {R}{W}formseal-sync{R}  {G}sync run{R}")
    print(G + " " + "\u2500" * 52 + R)
    br()

    written, skipped = _run_sync_once(cfg.get("output_folder", "data"))

    br()
    ok(f"{written} new, {skipped} duplicates")
    br()