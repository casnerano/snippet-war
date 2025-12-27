"""SQLAlchemy database models."""

import uuid
from datetime import datetime

from sqlalchemy import (
    BigInteger,
    Boolean,
    CheckConstraint,
    ForeignKey,
    Integer,
    String,
    Text,
    UniqueConstraint,
)
from sqlalchemy.dialects.postgresql import JSONB, TIMESTAMP, UUID
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column, relationship


class Base(DeclarativeBase):
    """Base class for all database models."""

    pass


class QuestionDB(Base):
    """Database model for questions table."""

    __tablename__ = "questions"

    id: Mapped[uuid.UUID] = mapped_column(
        UUID(as_uuid=True),
        primary_key=True,
        default=uuid.uuid4,
        server_default=None,
    )
    language: Mapped[str] = mapped_column(String(50), nullable=False)
    topic: Mapped[str] = mapped_column(String(255), nullable=False)
    difficulty: Mapped[str] = mapped_column(String(20), nullable=False)
    question_type: Mapped[str] = mapped_column(String(20), nullable=False)
    code: Mapped[str] = mapped_column(Text, nullable=False)
    question_text: Mapped[str] = mapped_column(Text, nullable=False)
    options: Mapped[list[str] | None] = mapped_column(JSONB, nullable=True)
    correct_answers: Mapped[list[str]] = mapped_column(JSONB, nullable=False)
    explanation: Mapped[str] = mapped_column(Text, nullable=False)
    likes_count: Mapped[int] = mapped_column(
        Integer, nullable=False, server_default="0"
    )
    created_at: Mapped[datetime] = mapped_column(
        TIMESTAMP(timezone=True),
        nullable=False,
        server_default=None,
        default=datetime.utcnow,
    )

    user_questions: Mapped[list["UserQuestionDB"]] = relationship(
        "UserQuestionDB", back_populates="question", cascade="all, delete-orphan"
    )

    __table_args__ = (
        CheckConstraint(
            "language IN ('python', 'javascript', 'go', 'java', 'cpp', 'rust', 'typescript')",
            name="chk_questions_language",
        ),
        CheckConstraint(
            "difficulty IN ('beginner', 'intermediate', 'advanced')",
            name="chk_questions_difficulty",
        ),
        CheckConstraint(
            "question_type IN ('multiple_choice', 'free_text')",
            name="chk_questions_question_type",
        ),
    )


class UserDB(Base):
    """Database model for users table."""

    __tablename__ = "users"

    id: Mapped[uuid.UUID] = mapped_column(
        UUID(as_uuid=True),
        primary_key=True,
        default=uuid.uuid4,
        server_default=None,
    )
    telegram_user_id: Mapped[int] = mapped_column(
        BigInteger, nullable=False, unique=True
    )
    created_at: Mapped[datetime] = mapped_column(
        TIMESTAMP(timezone=True),
        nullable=False,
        server_default=None,
        default=datetime.utcnow,
    )

    user_questions: Mapped[list["UserQuestionDB"]] = relationship(
        "UserQuestionDB", back_populates="user", cascade="all, delete-orphan"
    )


class UserQuestionDB(Base):
    """Database model for user_questions table."""

    __tablename__ = "user_questions"

    id: Mapped[uuid.UUID] = mapped_column(
        UUID(as_uuid=True),
        primary_key=True,
        default=uuid.uuid4,
        server_default=None,
    )
    user_id: Mapped[uuid.UUID] = mapped_column(
        UUID(as_uuid=True),
        ForeignKey("users.id", ondelete="CASCADE"),
        nullable=False,
    )
    question_id: Mapped[uuid.UUID] = mapped_column(
        UUID(as_uuid=True),
        ForeignKey("questions.id", ondelete="CASCADE"),
        nullable=False,
    )
    seen_at: Mapped[datetime] = mapped_column(
        TIMESTAMP(timezone=True),
        nullable=False,
        server_default=None,
        default=datetime.utcnow,
    )
    answered_at: Mapped[datetime | None] = mapped_column(
        TIMESTAMP(timezone=True), nullable=True
    )
    is_correct: Mapped[bool | None] = mapped_column(Boolean, nullable=True)

    user: Mapped["UserDB"] = relationship("UserDB", back_populates="user_questions")
    question: Mapped["QuestionDB"] = relationship(
        "QuestionDB", back_populates="user_questions"
    )

    __table_args__ = (
        UniqueConstraint(
            "user_id", "question_id", name="idx_user_questions_user_question"
        ),
    )
