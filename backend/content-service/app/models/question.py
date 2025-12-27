"""Question models."""

import uuid
from datetime import datetime

from app.models.enums import Difficulty, Language, QuestionType
from app.models.topics import is_valid_topic
from app.models.validation import (
    validate_free_text_answer,
    validate_multiple_choice_answer,
    validate_multiple_choice_options,
)
from pydantic import BaseModel, Field, field_validator, model_validator


class GenerateQuestionRequest(BaseModel):
    """Request model for question generation."""

    language: Language
    topic: str
    difficulty: Difficulty
    question_type: QuestionType

    @model_validator(mode="after")
    def validate_topic(self) -> "GenerateQuestionRequest":
        """Validate that topic is valid for the given language."""
        if not is_valid_topic(self.language, self.topic):
            raise ValueError(
                f"invalid topic '{self.topic}' for language '{self.language.value}'"
            )
        return self


class GetQuestionsBatchRequest(BaseModel):
    """Request model for getting batch of questions."""

    language: Language
    topics: list[str] = Field(..., min_length=1, description="List of topics")
    difficulty: Difficulty
    count: int = Field(..., ge=1, description="Number of questions")
    question_type: QuestionType = Field(
        default=QuestionType.MULTIPLE_CHOICE, description="Question type"
    )
    telegram_user_id: int | None = Field(
        default=None, description="Telegram user ID (optional)"
    )

    @model_validator(mode="after")
    def validate_topics(self) -> "GetQuestionsBatchRequest":
        """Validate that all topics are valid for the given language."""
        for topic in self.topics:
            if not is_valid_topic(self.language, topic):
                raise ValueError(
                    f"invalid topic '{topic}' for language '{self.language.value}'"
                )
        return self


class Question(BaseModel):
    """Question model."""

    id: str = Field(default_factory=lambda: str(uuid.uuid4()))
    language: Language
    topic: str
    difficulty: Difficulty
    question_type: QuestionType
    code: str
    question_text: str
    options: list[str] | None = None
    correct_answers: list[str]
    explanation: str
    created_at: datetime = Field(default_factory=datetime.utcnow)

    @field_validator("code")
    @classmethod
    def validate_code(cls, v: str) -> str:
        """Validate code field exists (can be empty string)."""
        # Code can be empty string for theoretical questions without code
        return v

    @field_validator("question_text")
    @classmethod
    def validate_question_text(cls, v: str) -> str:
        """Validate question text is not empty."""
        if not v:
            raise ValueError("question text is required")
        return v

    @field_validator("correct_answers")
    @classmethod
    def validate_correct_answers(cls, v: list[str]) -> list[str]:
        """Validate correct answer is not empty."""
        if not v or len(v) == 0:
            raise ValueError("correct answers must be a non-empty list")
        return v

    @field_validator("explanation")
    @classmethod
    def validate_explanation(cls, v: str) -> str:
        """Validate explanation is not empty."""
        if not v:
            raise ValueError("explanation is required")
        return v

    @model_validator(mode="after")
    def validate_question_type_specific(self) -> "Question":
        """Validate question type specific fields."""
        if self.question_type == QuestionType.MULTIPLE_CHOICE:
            if not self.options:
                raise ValueError("options are required for multiple choice questions")
            validate_multiple_choice_options(self.options)
            validate_multiple_choice_answer(self.correct_answers, self.options)
        elif self.question_type == QuestionType.FREE_TEXT:
            validate_free_text_answer(self.correct_answers)

        if not is_valid_topic(self.language, self.topic):
            raise ValueError(
                f"invalid topic '{self.topic}' for language '{self.language.value}'"
            )

        return self

    class Config:
        """Pydantic config."""

        populate_by_name = True
