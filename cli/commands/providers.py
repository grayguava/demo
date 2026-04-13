# cli/commands/providers.py
# List available providers and storage types

from cli.ui import br, rule, cmd_line, link, C, G, Y, O, W, D, R


def run(args):
    if not args:
        _list_providers()
        return

    # Handle --<provider> flags
    if args[0] == "--cloudflare":
        _show_cloudflare_types()
    elif args[0] == "--supabase":
        _show_supabase_types()
    else:
        _list_providers()


def _list_providers():
    br()
    print(f"{C} \u250c\u2500 {R}{W}formseal-sync{R}  {G}providers{R}")
    print(G + " " + "\u2500" * 52 + R)
    br()

    print(f"  {G}Available providers:{R}")
    print(f"    {W}cloudflare{R}   - Cloudflare KV")
    print(f"    {W}supabase{R}     - Supabase DB")
    br()

    print(f"  {D}Storage types:{R}")
    print(f"    Run: fsync providers --<provider>")
    br()
    print(f"    {W}fsync providers --cloudflare{R}")
    print(f"    {W}fsync providers --supabase{R}")
    br()


def _show_cloudflare_types():
    br()
    print(f"{C} \u250c\u2500 {R}{W}formseal-sync{R}  {G}cloudflare storage{R}")
    print(G + " " + "\u2500" * 52 + R)
    br()

    print(f"  {W}kv{R}       - Key-Value Store")
    print(f"     {D}Most common for form submissions{R}")
    br()
    print(f"  {W}d1{R}       - D1 Database (coming soon)")
    br()
    print(f"  {W}r2{R}       - R2 Object Storage (coming soon)")
    br()


def _show_supabase_types():
    br()
    print(f"{C} \u250c\u2500 {R}{W}formseal-sync{R}  {G}supabase storage{R}")
    print(G + " " + "\u2500" * 52 + R)
    br()

    print(f"  {W}db{R}       - PostgreSQL Database")
    print(f"     {D}Direct table access{R}")
    br()
    print(f"  {W}storage{R}  - Supabase Storage (coming soon)")
    br()


def _help():
    return [
        ("fsync providers", "list available providers"),
        ("fsync providers --cloudflare", "show Cloudflare storage types"),
        ("fsync providers --supabase", "show Supabase storage types"),
    ]