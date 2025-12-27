"""LLM response models."""

from typing import TYPE_CHECKING

from app.models.enums import Difficulty, Language, QuestionType
from app.models.topics import is_valid_topic
from app.models.validation import (
    validate_free_text_answer,
    validate_multiple_choice_answer,
    validate_multiple_choice_options,
)
from pydantic import BaseModel, field_validator, model_validator

if TYPE_CHECKING:
    from app.models.question import Question


class LLMQuestionResponse(BaseModel):
    """Response model from LLM for question generation."""

    code: str
    question: str
    question_type: QuestionType
    options: list[str] | None = None
    correct_answer: str
    acceptable_variants: list[str] | None = None
    case_sensitive: bool = False
    explanation: str
    difficulty: Difficulty
    topic: str
    language: Language

    @field_validator("code")
    @classmethod
    def validate_code(cls, v: str) -> str:
        """Validate code is not empty."""
        if not v:
            raise ValueError("code is required")
        return v

    @field_validator("question")
    @classmethod
    def validate_question(cls, v: str) -> str:
        """Validate question is not empty."""
        if not v:
            raise ValueError("question is required")
        return v

    @field_validator("explanation")
    @classmethod
    def validate_explanation(cls, v: str) -> str:
        """Validate explanation is not empty."""
        if not v:
            raise ValueError("explanation is required")
        return v

    @model_validator(mode="after")
    def validate_response(self) -> "LLMQuestionResponse":
        """Validate LLM response structure."""
        if not is_valid_topic(self.language, self.topic):
            raise ValueError(
                f"invalid topic '{self.topic}' for language '{self.language.value}'"
            )

        if self.question_type == QuestionType.MULTIPLE_CHOICE:
            if not self.options:
                raise ValueError("options are required for multiple choice questions")
            validate_multiple_choice_options(self.options)
            validate_multiple_choice_answer(self.correct_answer, self.options)
        elif self.question_type == QuestionType.FREE_TEXT:
            validate_free_text_answer(self.correct_answer)

        return self

    def to_question(self) -> "Question":
        """Convert LLM response to Question model."""
        # Import here to avoid circular dependency
        from app.models.question import Question  # noqa: PLC0415

        question = Question(
            language=self.language,
            topic=self.topic,
            difficulty=self.difficulty,
            question_type=self.question_type,
            code=self.code,
            question_text=self.question,
            options=self.options,
            correct_answer=self.correct_answer,
            acceptable_variants=self.acceptable_variants,
            case_sensitive=self.case_sensitive,
            explanation=self.explanation,
        )

        return question
