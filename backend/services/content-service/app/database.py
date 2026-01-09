"""Database connection and session management."""

from collections.abc import AsyncGenerator

from app.config import DatabaseConfig
from fastapi import Request
from loguru import logger
from sqlalchemy.ext.asyncio import (
    AsyncEngine,
    AsyncSession,
    async_sessionmaker,
    create_async_engine,
)


class Database:
    """Database connection manager."""

    def __init__(self, config: DatabaseConfig) -> None:
        """Initialize database connection."""
        self.config = config
        self.engine: AsyncEngine | None = None
        self.session_factory: async_sessionmaker[AsyncSession] | None = None

    async def init(self) -> None:
        """Initialize database engine and session factory."""
        logger.info("Initializing database connection", url=self.config.url)
        self.engine = create_async_engine(
            self.config.url,
            pool_size=self.config.pool_size,
            max_overflow=self.config.max_overflow,
            echo=False,
        )
        self.session_factory = async_sessionmaker(
            self.engine,
            class_=AsyncSession,
            expire_on_commit=False,
            autoflush=False,
            autocommit=False,
        )
        logger.info("Database connection initialized")

    async def close(self) -> None:
        """Close database connections."""
        if self.engine:
            logger.info("Closing database connections")
            await self.engine.dispose()
            self.engine = None
            self.session_factory = None
            logger.info("Database connections closed")


async def get_db_session(request: Request) -> AsyncGenerator[AsyncSession, None]:
    """
    Dependency for getting database session.

    Args:
        request: FastAPI request object

    Yields:
        Database session

    Raises:
        RuntimeError: If database not initialized
    """
    db: Database = request.app.state.database
    if db is None or db.session_factory is None:
        raise RuntimeError("Database not initialized")

    async with db.session_factory() as session:
        try:
            yield session
            await session.commit()
        except Exception:
            await session.rollback()
            raise
        finally:
            await session.close()
