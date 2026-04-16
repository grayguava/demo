# Cloudflare provider

from cli.providers import Provider
from cli.providers.cloudflare.account import get_account_id, TokenError, AuthError
from cli.providers.cloudflare.storage.kv import fetch as kv_fetch


class CloudflareProvider(Provider):
    name = "cloudflare"
    display_name = "Cloudflare KV"
    storage_type = "Key-Value"

    def get_config_fields(self):
        return [
            {"key": "namespace", "prompt": "KV Namespace ID", "required": True},
        ]

    def authenticate(self, token):
        return get_account_id(token)

    def fetch(self, token, config, output_path):
        namespace = config.get("namespace")
        if not namespace:
            return 0, 0
        account_id = get_account_id(token)
        return kv_fetch(namespace=namespace, account_id=account_id, token=token, output_path=output_path)


Provider = CloudflareProvider