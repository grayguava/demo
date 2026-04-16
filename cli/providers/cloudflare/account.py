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
    
    token = token.strip()
    if not token:
        raise TokenError("Token is blank after stripping whitespace")
    
    req = urllib.request.Request(
        "https://api.cloudflare.com/client/v4/accounts",
        headers={"Authorization": f"Bearer {token}"}
    )
    
    try:
        with urllib.request.urlopen(req, timeout=10) as resp:
            data = json.loads(resp.read().decode())
    except urllib.error.HTTPError as e:
        body = e.read().decode("utf-8", errors="replace")
        raise AuthError(f"HTTP {e.code}: {body}")
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