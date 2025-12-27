"""Repository for user database operations."""

from app.exceptions import DatabaseError
from app.models.db import UserDB
from loguru import logger
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession


class UserRepository:
    """Repository for user database operations."""

    @staticmethod
    async def get_user_by_telegram_id(
        db: AsyncSession, telegram_user_id: int
    ) -> UserDB | None:
        """
        Get user by Telegram user ID.

        Args:
            db: Database session
            telegram_user_id: Telegram user ID

        Returns:
            User database model or None if not found
        """
        query = select(UserDB).where(UserDB.telegram_user_id == telegram_user_id)
        result = await db.execute(query)
        return result.scalar_one_or_none()

    @staticmethod
    async def get_or_create_user_by_telegram_id(
        db: AsyncSession, telegram_user_id: int
    ) -> UserDB:
        """
        Get or create user by Telegram user ID.

        Args:
            db: Database session
            telegram_user_id: Telegram user ID

        Returns:
            User database model

        Raises:
            DatabaseError: If create operation fails
        """
        user = await UserRepository.get_user_by_telegram_id(db, telegram_user_id)
        if user is not None:
            return user

        try:
            user = UserDB(telegram_user_id=telegram_user_id)
            db.add(user)
            await db.flush()
            await db.refresh(user)
            logger.info(
                "User created",
                user_id=str(user.id),
                telegram_user_id=telegram_user_id,
            )
            return user
        except Exception as e:
            logger.error("Failed to create user", error=str(e))
            raise DatabaseError(f"Не удалось создать пользователя: {e}") from e
