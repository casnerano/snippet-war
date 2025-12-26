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
    correct_answer: str
    acceptable_variants: list[str] | None = None
    case_sensitive: bool = False
    explanation: str
    created_at: datetime = Field(default_factory=datetime.utcnow)

    @field_validator("code")
    @classmethod
    def validate_code(cls, v: str) -> str:
        """Validate code is not empty."""
        if not v:
            raise ValueError("code is required")
        return v

    @field_validator("question")
    @classmethod
    def validate_question_text(cls, v: str) -> str:
        """Validate question text is not empty."""
        if not v:
            raise ValueError("question text is required")
        return v

    @field_validator("correct_answer")
    @classmethod
    def validate_correct_answer(cls, v: str) -> str:
        """Validate correct answer is not empty."""
        if not v:
            raise ValueError("correct answer is required")
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
            validate_multiple_choice_answer(self.correct_answer, self.options)
        elif self.question_type == QuestionType.FREE_TEXT:
            validate_free_text_answer(self.correct_answer)

        if not is_valid_topic(self.language, self.topic):
            raise ValueError(
                f"invalid topic '{self.topic}' for language '{self.language.value}'"
            )

        return self

    class Config:
        """Pydantic config."""

        populate_by_name = True
