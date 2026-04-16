# List providers

from pathlib import Path

from cli.ui import br, C, G, Y, W, D, R
from cli.providers import get_providers


def _load_version():
    p = Path(__file__).parent.parent.parent / "version.txt"
    if p.exists():
        return p.read_text().strip()
    return "dev"


VERSION = _load_version()


def run(args):
    _list_providers()


def _list_providers():
    providers = get_providers()

    br()
    print(f"{C} \u250c\u2500 {R}{W}formseal-fetch{R}  {Y}v{VERSION}{R}")
    print(G + " " + "\u2500" * 52 + R)
    br()

    print(f"  {G}Available providers:{R}")
    br()

    for name, provider in providers.items():
        print(f"    {W}>{provider.display_name}{R}   -  {provider.storage_type}")

    br()