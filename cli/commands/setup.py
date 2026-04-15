# cli/commands/setup.py
# Interactive connect commands

import sys
import json
from pathlib import Path

from cli.ui import br, fail, ok, info, warn, G, W, D, C, Y, O, R
from cli.commands.config import load_config, save_config
from cli.security import tokens


SETUP_SCHEMA = {
    "cloudflare": {
        "kv": {
            "config_fields": [
                {"key": "namespace", "prompt": "KV Namespace ID", "required": True},
            ],
        }
    }
}


def _parse_args(args):
    parsed = {}
    for arg in args:
        if ":" not in arg:
            fail(f"Invalid format: {arg}\n           Use flag:value (e.g., provider:cloudflare)")
        key, value = arg.split(":", 1)
        parsed[key] = value
    return parsed


def run(args):
    if not args:
        fail("Usage: fsf connect provider:cloudflare [namespace:<id>] [output:<path>]")

    parsed = _parse_args(args)

    if "provider" not in parsed:
        fail("provider is required.\n           Usage: fsf connect provider:<cloudflare> [...]")

    provider = parsed["provider"].lower()
    if provider not in SETUP_SCHEMA:
        fail(f"Unknown provider: {provider}\n           Run fsf providers to see available.")

    _setup_flow(provider, parsed)


def _setup_flow(provider, parsed):
    print()
    print(f"{C} \u250c\u2500 {R}{W}formseal-fetch{R}  {G}setup{R}")
    print(G + " " + "\u2500" * 52 + R)
    print()

    cfg = load_config()
    cfg["provider"] = provider

    # Get config fields from schema
    fields = SETUP_SCHEMA[provider]["kv"]["config_fields"]

    # Ask for namespace if not provided
    namespace = parsed.get("namespace")
    if not namespace:
        for field in fields:
            prompt = field["prompt"]
            try:
                sys.stdout.write(f"  {prompt}: ")
                sys.stdout.flush()
                namespace = input().strip()
            except KeyboardInterrupt:
                br()
                info("Cancelled.")
                br()
                return
            if field.get("required") and not namespace:
                fail(f"{prompt} is required")
            break

    if namespace:
        tokens.save_namespace(provider, namespace)

    # Ask for token
    if "token" in parsed:
        token = parsed["token"]
    else:
        try:
            sys.stdout.write("  Account API Token: ")
            sys.stdout.flush()
            token = sys.stdin.readline().strip()
        except KeyboardInterrupt:
            br()
            info("Cancelled.")
            br()
            return
        if not token:
            fail("Account API Token is required")

    token = "".join(c for c in token if c.isprintable()).strip()
    if not token:
        fail("Account API Token is required")

    tokens.save_token(provider, token)

    # Ask for output folder
    if "output" in parsed:
        output_folder = parsed["output"]
    else:
        try:
            sys.stdout.write(f"  Output Folder [{D}data{R}]: ")
            sys.stdout.flush()
            output_folder = input().strip()
            if not output_folder:
                output_folder = "data"
        except KeyboardInterrupt:
            br()
            info("Cancelled.")
            br()
            return
    cfg["output_folder"] = output_folder
    print()

    save_config(cfg)

    print(f"{G} \u2713{R} Saved!")
    print()
    print(f"  Run {W}fsf fetch{R} to download ciphertexts")
    print()