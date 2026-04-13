#!/usr/bin/env python3
# cli/fsync.py
# fsync - formseal sync engine

import sys
import os

SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
PROJECT_DIR = os.path.dirname(SCRIPT_DIR)
sys.path.insert(0, PROJECT_DIR)

from cli.ui import br, rule, cmd_line, link, fail, ok, info, warn, C, G, Y, O, W, D, R
from cli.commands import fetch, config, providers, setup
from cli.commands.config import load_config
from cli.daemon.sync import run as sync_run


def intro():
    br()
    print(f"{C} \u250c\u2500 {R}{W}formseal-sync{R}")
    print(G + " " + "\u2500" * 52 + R)
    br()
    print(f"  {D}Local CLI tool for syncing encrypted form submissions{R}")
    print(f"  {D}from storage backends to your PC.{R}")
    br()
    print(f"  {G}Quick Start:{R}")
    br()
    print(f"  {W}fsync setup quick{R}")
    print(f"  {D}configure storage backend{R}")
    br()
    print(f"  {W}fsync sync start{R}")
    print(f"  {D}start background sync{R}")
    br()
    print(f"  {G}Help:{R}")
    br()
    print(f"  {W}fsync --help{R}")
    link("https://github.com/formseal/formseal-sync")
    br()


def main():
    args = sys.argv[1:]
    command = args[0] if args else None

    match command:
        case "setup":
            setup.run(args[1:])

        case "help":
            _help(args[1:])

        case "set":
            config.run_set(args[1:])

        case "fetch":
            fetch.run(args[1:])

        case "providers":
            providers.run(args[1:])

        case "status":
            config.run_status()

        case "config":
            config.run_config_show()

        case "logout":
            config.run_logout()

        case "sync":
            sync_run(args[1:])

        case None | "fsync":
            intro()

        case "--help" | "-h":
            _help()

        case "--about":
            _about()

        case _:
            fail(f"Unknown command: {command}\n           Run {W}fsync --help{R} for usage.")


def _help(args=None):
    if args and args[0].startswith("--vars"):
        os_name = args[0].replace("--vars-", "").replace("--vars", "")
        if os_name == "":
            os_name = None  # auto-detect
        _help_vars(os_name)
        return

    br()
    print(f"{C} \u250c\u2500 {R}{W}formseal-sync{R}  {G}help{R}")
    print(G + " " + "\u2500" * 52 + R)
    br()

    print(f"  {G}>>{R} {Y}Setup{R}")
    print(G + " " + "\u2500" * 52 + R)
    cmd_line("fsync setup quick", "interactive setup wizard")
    cmd_line("fsync setup sync", "configure sync settings")
    cmd_line("fsync setup reset", "clear configuration")
    br()

    print(f"  {G}>>{R} {Y}Sync{R}")
    print(G + " " + "\u2500" * 52 + R)
    cmd_line("fsync sync start", "start background sync")
    cmd_line("fsync sync stop", "stop background sync")
    cmd_line("fsync sync status", "show sync status")
    cmd_line("fsync sync run", "run sync once")
    br()

    print(f"  {G}>>{R} {Y}Fetch{R}")
    print(G + " " + "\u2500" * 52 + R)
    cmd_line("fsync fetch", "download ciphertexts")
    cmd_line("fsync fetch --output <file>", "custom output path")
    br()

    print(f"  {G}>>{R} {Y}Status{R}")
    print(G + " " + "\u2500" * 52 + R)
    cmd_line("fsync status", "show configuration")
    cmd_line("fsync logout", "clear credentials")
    br()

    print(f"  {G}>>{R} {Y}Env Vars{R}")
    print(G + " " + "\u2500" * 52 + R)
    cmd_line("fsync help --vars", "env var setup for your OS")
    br()

    print(f"  {G}>>{R} {Y}Docs{R}")
    print(G + " " + "\u2500" * 52 + R)
    link("https://github.com/formseal/formseal-sync")
    br()


def _help_vars(os_name=None):
    import platform
    
    # Auto-detect if not specified
    if not os_name:
        system = platform.system().lower()
        if system == "windows":
            os_name = "windows"
        elif system == "darwin":
            os_name = "macos"
        else:
            os_name = "linux"

    cfg = load_config()
    provider = cfg.get("provider", "cloudflare")

    br()
    print(f"{C} \u250c\u2500 {R}{W}formseal-sync{R}  {G}env vars{R}")
    print(G + " " + "\u2500" * 52 + R)
    br()

    # Show what each env var is for
    print(f"  {Y}What to set:{R}")
    if provider == "cloudflare":
        print(f"    FSYNC_CF_TOKEN  - Cloudflare API token (cfut_...)")
    else:
        print(f"    FSYNC_SU_KEY    - Supabase service key (eyJ...)")
    br()

    if os_name in ["windows", "win", "win32"]:
        print(f"  {Y}Windows (CMD){R}")
        print(f"    setx FSYNC_CF_TOKEN \"cfut_...\"")
        print(f"    setx FSYNC_SU_KEY \"eyJ...\"")
        br()
        print(f"  {Y}Windows (PowerShell){R}")
        print(f"    [System.Environment]::SetEnvironmentVariable('FSYNC_CF_TOKEN', 'cfut_...', 'User')")
        print(f"    [System.Environment]::SetEnvironmentVariable('FSYNC_SU_KEY', 'eyJ...', 'User')")
        br()
        print(f"  {D}Open a new terminal window for changes to take effect.{R}")
    elif os_name in ["macos", "darwin", "mac"]:
        print(f"  {Y}macOS (Terminal){R}")
        print(f"    export FSYNC_CF_TOKEN=\"cfut_...\"")
        print(f"    export FSYNC_SU_KEY=\"eyJ...\"")
        br()
        print(f"  Add to ~/.zshrc or ~/.bash_profile for persistence:")
        print(f"    echo 'export FSYNC_CF_TOKEN=\"cfut_...\"' >> ~/.zshrc")
    elif os_name in ["linux", "linux2", "unix"]:
        print(f"  {Y}Linux (Terminal){R}")
        print(f"    export FSYNC_CF_TOKEN=\"cfut_...\"")
        print(f"    export FSYNC_SU_KEY=\"eyJ...\"")
        br()
        print(f"  Add to ~/.bashrc, ~/.bash_profile, or ~/.profile for persistence:")
        print(f"    echo 'export FSYNC_CF_TOKEN=\"cfut_...\"' >> ~/.bashrc")
    else:
        print(f"  Unknown OS. Use: windows, macos, or linux")
        print(f"  Example: fsync help --vars-windows")

    br()


def _about():
    br()
    print(f" {C}┌─{R} {W}formseal-sync{R}")
    print(f" {C}───────────────────────────────────────────────────{R}")
    br()
    print(f"  {D}Author:{R} formseal")
    print(f"  {D}License:{R} MIT")
    br()
    print(f"  {D}Project:{R}")
    link("https://github.com/formseal/formseal-sync")
    br()


if __name__ == "__main__":
    main()