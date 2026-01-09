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
    correct_answers: list[str]
    explanation: str
    difficulty: Difficulty
    topic: str
    language: Language

    @field_validator("code")
    @classmethod
    def validate_code(cls, v: str) -> str:
        """Validate code field exists (can be empty string)."""
        # Code can be empty string for theoretical questions without code
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
            validate_multiple_choice_answer(self.correct_answers, self.options)
        elif self.question_type == QuestionType.FREE_TEXT:
            validate_free_text_answer(self.correct_answers)

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
            correct_answers=self.correct_answers,
            explanation=self.explanation,
        )

        return question


class LLMQuestionsResponse(BaseModel):
    """Response model from LLM for multiple questions generation."""

    questions: list[LLMQuestionResponse]

    @field_validator("questions")
    @classmethod
    def validate_questions(
        cls, v: list[LLMQuestionResponse]
    ) -> list[LLMQuestionResponse]:
        """Validate questions list is not empty."""
        if not v or len(v) == 0:
            raise ValueError("questions list must not be empty")
        return v

    def to_questions(self) -> list["Question"]:
        """Convert LLM response to list of Question models."""
        return [q.to_question() for q in self.questions]
