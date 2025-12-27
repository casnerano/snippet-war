"""Application configuration."""

from loguru import logger
from pydantic_settings import BaseSettings, SettingsConfigDict


class ProxyAPIConfig(BaseSettings):
    """ProxyAPI configuration."""

    model_config = SettingsConfigDict(
        env_file=[
            ".env",
        ],
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


class DatabaseConfig(BaseSettings):
    """Database configuration."""

    model_config = SettingsConfigDict(
        env_file=[
            ".env",
        ],
        env_prefix="DATABASE_",
        case_sensitive=False,
        env_file_encoding="utf-8",
        enable_decoding=False,
        extra="ignore",
    )

    host: str = "localhost"
    port: int = 5432
    user: str = "snippet_war"
    password: str = "snippet_war"
    database: str = "snippet_war"
    pool_size: int = 5
    max_overflow: int = 10

    @property
    def url(self) -> str:
        """Get database URL for asyncpg."""
        return (
            f"postgresql+asyncpg://{self.user}:{self.password}@"
            f"{self.host}:{self.port}/{self.database}"
        )


class Config:
    """Application configuration."""

    def __init__(self) -> None:
        """Initialize config with nested settings."""
        # Load nested configs
        self.proxyapi = ProxyAPIConfig()
        self.database = DatabaseConfig()


def load_config() -> Config:
    """Load configuration from environment variables."""
    logger.info("Loading configuration from environment variables")
    cfg = Config()
    logger.info(f"Config: {cfg.proxyapi}")
    logger.info(
        "Database config loaded",
        host=cfg.database.host,
        port=cfg.database.port,
        database=cfg.database.database,
    )
    return cfg
