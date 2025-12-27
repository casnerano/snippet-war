"""LLM client interface."""

from typing import Protocol

from app.models.llm_response import LLMQuestionResponse


class LLMClient(Protocol):
    """Protocol for LLM clients."""

    async def generate_question(self, prompt: str) -> LLMQuestionResponse:
        """
        Generate question response from LLM.

        Args:
            prompt: Prompt text for LLM

        Returns:
            LLMQuestionResponse with generated question data

        Raises:
            Exception: If LLM request fails
        """
        ...
