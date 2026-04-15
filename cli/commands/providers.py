# cli/commands/providers.py
# List available providers

from pathlib import Path

from cli.ui import br, C, G, Y, W, D, R


def _load_version():
    p = Path(__file__).parent.parent.parent / "version.txt"
    if p.exists():
        return p.read_text().strip()
    return "dev"


VERSION = _load_version()


def run(args):
    _list_providers()


def _list_providers():
    br()
    print(f"{C} \u250c\u2500 {R}{W}formseal-fetch{R}  {Y}v{VERSION}{R}")
    print(G + " " + "\u2500" * 52 + R)
    br()

    print(f"  {G}Available providers:{R}")
    br()
    print(f"    {W}> Cloudflare{R}   -  Cloudflare KV")
    br()
    print(f"    {W}> Supabase{R}     -  PostgreSQL DB (coming soon)")
    br()


def _help():
    return [
        ("fsf providers", "list available providers"),
    ]