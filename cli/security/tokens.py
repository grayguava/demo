# Credential storage (keyring + JSON fallback)

import json
import base64
from pathlib import Path

try:
    import keyring
    HAS_KEYRING = True
except ImportError:
    HAS_KEYRING = False

SERVICE = "formseal-fetch"

CONFIG_DIR = Path.home() / ".config" / "formseal-fetch"
SECRETS_FILE = CONFIG_DIR / "secrets.json"


def _ensure_config_dir():
    CONFIG_DIR.mkdir(parents=True, exist_ok=True)


def _load_secrets() -> dict:
    if not SECRETS_FILE.exists():
        return {}
    try:
        return json.loads(SECRETS_FILE.read_text())
    except:
        return {}


def _save_secrets(secrets: dict):
    _ensure_config_dir()
    SECRETS_FILE.write_text(json.dumps(secrets, indent=2))


def save_token(provider: str, token: str) -> bool:
    """Save API token to OS keyring. Falls back to JSON if keyring fails."""
    if HAS_KEYRING:
        try:
            keyring.set_password(SERVICE, f"{provider}:api-token", token)
            return True
        except Exception:
            pass
    
    # Fallback to JSON file
    secrets = _load_secrets()
    secrets[f"{provider}:token"] = base64.b64encode(token.encode()).decode()
    _save_secrets(secrets)
    return True


def load_token(provider: str) -> str | None:
    """Load API token from OS keyring. Falls back to JSON if not in keyring."""
    if HAS_KEYRING:
        try:
            token = keyring.get_password(SERVICE, f"{provider}:api-token")
            if token:
                return token
        except Exception:
            pass
    
    # Fallback to JSON file
    secrets = _load_secrets()
    encoded = secrets.get(f"{provider}:token")
    if encoded:
        return base64.b64decode(encoded.encode()).decode().strip()
    return None


def delete_token(provider: str):
    """Delete API token from keyring. Falls back to JSON if keyring fails."""
    if HAS_KEYRING:
        try:
            keyring.delete_password(SERVICE, f"{provider}:api-token")
            return
        except Exception:
            pass
    
    secrets = _load_secrets()
    secrets.pop(f"{provider}:token", None)
    _save_secrets(secrets)


def save_namespace(provider: str, namespace: str) -> bool:
    """Save KV namespace ID to OS keyring. Falls back to JSON."""
    if HAS_KEYRING:
        try:
            keyring.set_password(SERVICE, f"{provider}:kv-namespace", namespace)
            return True
        except Exception:
            pass
    
    # Fallback to JSON file
    secrets = _load_secrets()
    secrets[f"{provider}:namespace"] = base64.b64encode(namespace.encode()).decode()
    _save_secrets(secrets)
    return True


def load_namespace(provider: str) -> str | None:
    """Load KV namespace ID from OS keyring. Falls back to JSON."""
    if HAS_KEYRING:
        try:
            namespace = keyring.get_password(SERVICE, f"{provider}:kv-namespace")
            if namespace:
                return namespace
        except Exception:
            pass
    
    # Fallback to JSON file
    secrets = _load_secrets()
    encoded = secrets.get(f"{provider}:namespace")
    if encoded:
        return base64.b64decode(encoded.encode()).decode().strip()
    return None


def delete_namespace(provider: str):
    """Delete namespace from keyring. Falls back to JSON if keyring fails."""
    if HAS_KEYRING:
        try:
            keyring.delete_password(SERVICE, f"{provider}:kv-namespace")
            return
        except Exception:
            pass
    
    secrets = _load_secrets()
    secrets.pop(f"{provider}:namespace", None)
    _save_secrets(secrets)


def token_location(provider: str) -> str:
    """Check where token is stored."""
    if HAS_KEYRING:
        try:
            if keyring.get_password(SERVICE, f"{provider}:api-token"):
                return "OS Keychain"
        except Exception:
            pass
    
    secrets = _load_secrets()
    if f"{provider}:token" in secrets:
        return "Config File"
    return "Not set"


def namespace_location(provider: str) -> str:
    """Check where namespace is stored."""
    if HAS_KEYRING:
        try:
            if keyring.get_password(SERVICE, f"{provider}:kv-namespace"):
                return "OS Keychain"
        except Exception:
            pass
    
    secrets = _load_secrets()
    if f"{provider}:namespace" in secrets:
        return "Config File"
    return "Not set"


def clear_all(provider: str):
    """Clear all secrets for a provider."""
    delete_token(provider)
    delete_namespace(provider)