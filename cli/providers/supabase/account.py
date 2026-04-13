# cli/providers/supabase/account.py
# Supabase account authentication

import json
import urllib.request
import urllib.error

from cli.ui import fail


def validate_token(url, key):
    """Check if Supabase credentials are valid."""
    req = urllib.request.Request(
        f"{url}/rest/v1/",
        headers={
            "Authorization": f"Bearer {key}",
            "apikey": key,
        }
    )
    try:
        with urllib.request.urlopen(req, timeout=10) as resp:
            return resp.status == 200
    except:
        return False