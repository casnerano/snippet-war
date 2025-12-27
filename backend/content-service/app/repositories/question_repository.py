"""Repository for question database operations."""

import uuid
from datetime import datetime

from app.exceptions import DatabaseError
from app.models.db import QuestionDB, UserQuestionDB
from app.models.enums import Difficulty, Language, QuestionType
from app.models.question import Question
from loguru import logger
from sqlalchemy import and_, select
from sqlalchemy.dialects.postgresql import insert
from sqlalchemy.ext.asyncio import AsyncSession


class QuestionRepository:
    """Repository for question database operations."""

    @staticmethod
    async def save_question(db: AsyncSession, question: Question) -> QuestionDB:
        """
        Save question to database.

        Args:
            db: Database session
            question: Question to save

        Returns:
            Saved question database model

        Raises:
            DatabaseError: If save operation fails
        """
        try:
            question_db = QuestionDB(
                id=uuid.UUID(question.id),
                language=question.language.value,
                topic=question.topic,
                difficulty=question.difficulty.value,
                question_type=question.question_type.value,
                code=question.code,
                question_text=question.question_text,
                options=question.options,
                correct_answers=question.correct_answers,
                case_sensitive=question.case_sensitive,
                explanation=question.explanation,
                created_at=question.created_at,
            )
            db.add(question_db)
            await db.flush()
            await db.refresh(question_db)
            logger.debug("Question saved to database", question_id=question.id)
            return question_db
        except Exception as e:
            logger.error("Failed to save question to database", error=str(e))
            raise DatabaseError(f"Не удалось сохранить вопрос: {e}") from e

    @staticmethod
    async def get_questions_by_filters(
        db: AsyncSession,
        language: Language,
        topic: str,
        difficulty: Difficulty,
        limit: int | None = None,
    ) -> list[QuestionDB]:
        """
        Get questions by filters.

        Args:
            db: Database session
            language: Programming language
            topic: Topic
            difficulty: Difficulty level (optional)
            limit: Maximum number of questions to return

        Returns:
            List of question database models
        """
        query = (
            select(QuestionDB)
            .where(
                and_(
                    QuestionDB.language == language.value,
                    QuestionDB.topic == topic,
                    QuestionDB.difficulty == difficulty.value,
                )
            )
            .order_by(QuestionDB.created_at.desc())
        )

        if limit is not None:
            query = query.limit(limit)

        result = await db.execute(query)
        questions = result.scalars().all()
        return list(questions)

    @staticmethod
    async def get_unseen_questions(
        db: AsyncSession,
        user_id: uuid.UUID,
        language: Language,
        topic: str,
        difficulty: Difficulty,
        limit: int | None = None,
    ) -> list[QuestionDB]:
        """
        Get questions that user hasn't seen yet.

        Args:
            db: Database session
            user_id: User ID
            language: Programming language
            topic: Topic
            difficulty: Difficulty level (optional)
            limit: Maximum number of questions to return

        Returns:
            List of unseen question database models
        """
        # Subquery to get question IDs that user has seen
        seen_subquery = select(UserQuestionDB.question_id).where(
            UserQuestionDB.user_id == user_id
        )

        # Main query to get questions not in seen_subquery
        query = select(QuestionDB).where(
            and_(
                QuestionDB.language == language.value,
                QuestionDB.topic == topic,
                QuestionDB.difficulty == difficulty.value,
                ~QuestionDB.id.in_(seen_subquery),
            )
        )

        if limit is not None:
            query = query.limit(limit)

        result = await db.execute(query)
        questions = result.scalars().all()
        return list(questions)

    @staticmethod
    async def mark_questions_as_seen(
        db: AsyncSession, user_id: uuid.UUID, question_ids: list[uuid.UUID]
    ) -> None:
        """
        Mark questions as seen by user.

        Args:
            db: Database session
            user_id: User ID
            question_ids: List of question IDs to mark as seen

        Raises:
            DatabaseError: If mark operation fails
        """
        try:
            # Use insert with on_conflict_do_nothing to handle duplicates
            stmt = insert(UserQuestionDB).values(
                [
                    {
                        "user_id": user_id,
                        "question_id": qid,
                        "seen_at": datetime.utcnow(),
                    }
                    for qid in question_ids
                ]
            )
            stmt = stmt.on_conflict_do_nothing(
                index_elements=["user_id", "question_id"]
            )
            await db.execute(stmt)
            logger.debug(
                "Questions marked as seen",
                user_id=str(user_id),
                count=len(question_ids),
            )
        except Exception as e:
            logger.error("Failed to mark questions as seen", error=str(e))
            raise DatabaseError(
                f"Не удалось отметить вопросы как просмотренные: {e}"
            ) from e

    @staticmethod
    def question_db_to_model(question_db: QuestionDB) -> Question:
        """
        Convert QuestionDB to Question model.

        Args:
            question_db: Question database model

        Returns:
            Question model
        """
        return Question(
            id=str(question_db.id),
            language=Language(question_db.language),
            topic=question_db.topic,
            difficulty=Difficulty(question_db.difficulty),
            question_type=QuestionType(question_db.question_type),
            code=question_db.code,
            question_text=question_db.question_text,
            options=question_db.options,
            correct_answers=question_db.correct_answers,
            case_sensitive=question_db.case_sensitive,
            explanation=question_db.explanation,
            created_at=question_db.created_at,
        )
