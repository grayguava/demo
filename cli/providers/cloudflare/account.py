# cli/providers/cloudflare/account.py
# Cloudflare account authentication

import json
import urllib.request
import urllib.error

from cli.ui import fail


class AuthError(Exception):
    """Raised when Cloudflare authentication fails."""
    pass


class TokenError(Exception):
    """Raised when token is invalid or missing."""
    pass


def get_account_id(token):
    """Fetch account_id from Cloudflare API."""
    if not token:
        raise TokenError("Token is empty")
    
    req = urllib.request.Request(
        "https://api.cloudflare.com/client/v4/accounts",
        headers={"Authorization": f"Bearer {token}"}
    )
    
    try:
        with urllib.request.urlopen(req, timeout=10) as resp:
            data = json.loads(resp.read().decode())
    except urllib.error.HTTPError as e:
        if e.code == 401:
            raise AuthError("Invalid token (unauthorized)")
        elif e.code == 403:
            raise AuthError("Token lacks required permissions")
        elif e.code == 429:
            raise AuthError("Rate limit exceeded")
        else:
            raise AuthError(f"HTTP {e.code}: {e.reason}")
    except json.JSONDecodeError as e:
        raise AuthError(f"Invalid API response: {e}")
    except urllib.error.URLError as e:
        raise AuthError(f"Network error: {e.reason}")
    
    if not data.get("success"):
        errors = data.get("errors", [])
        if errors:
            msg = errors[0].get("message", "Unknown error")
            if "token" in msg.lower():
                raise TokenError(f"Invalid token: {msg}")
            raise AuthError(msg)
        raise AuthError("Unknown API error")
    
    accounts = data.get("result", [])
    if not accounts:
        raise AuthError("No accounts found. Token needs 'Account Settings: Read' scope.")
    
    return accounts[0]["id"]


def validate_token(token):
    """Check if token is valid. Returns True/False, never raises."""
    try:
        get_account_id(token)
        return True
    except (TokenError, AuthError) as e:
        return False
    except Exception:
        return False