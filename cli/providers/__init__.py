# Provider base class

from abc import ABC, abstractmethod


class Provider(ABC):
    name = None
    display_name = None
    storage_type = None

    @abstractmethod
    def get_config_fields(self):
        pass

    @abstractmethod
    def authenticate(self, token):
        pass

    @abstractmethod
    def fetch(self, token, config, output_path):
        pass


def get_providers():
    providers = {}
    for mod in ["cloudflare"]:
        try:
            cls = __import__(f"cli.providers.{mod}", fromlist=["Provider"]).Provider
            providers[mod] = cls()
        except Exception:
            pass
    return providers


def get_provider(name):
    providers = get_providers()
    return providers.get(name)


def get_provider_schema():
    providers = get_providers()
    schema = {}
    for name, provider in providers.items():
        schema[name] = {"config_fields": provider.get_config_fields()}
    return schema