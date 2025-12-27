"""Application configuration."""

from loguru import logger
from pydantic_settings import BaseSettings, SettingsConfigDict


class ProxyAPIConfig(BaseSettings):
    """ProxyAPI configuration."""

    model_config = SettingsConfigDict(
        env_file=[
            ".env.root",
            ".env.service",
            ".env",
        ],  # Load in order: root, service, local
        env_prefix="PROXYAPI_",
        case_sensitive=False,
        env_file_encoding="utf-8",
        enable_decoding=False,
        extra="ignore",
    )

    api_key: str = ""
    model: str = "gpt-4.1-mini"
    base_url: str = "https://api.proxyapi.ru/openai/v1"
    timeout: int = 30
    max_tokens: int = 2000


class Config:
    """Application configuration."""

    def __init__(self) -> None:
        """Initialize config with nested settings."""
        # Load nested configs
        self.proxyapi = ProxyAPIConfig()


def load_config() -> Config:
    """Load configuration from environment variables."""
    logger.info("Loading configuration from environment variables")
    cfg = Config()
    logger.info(f"Config: {cfg.proxyapi}")
    return cfg
