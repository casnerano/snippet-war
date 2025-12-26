"""Question generation service."""

from loguru import logger

from app.clients import LLMClient, build_prompt
from app.models import (
    GenerateQuestionRequest,
    LLMQuestionResponse,
    Question,
)


class QuestionService:
    """Service for generating questions using LLM."""

    def __init__(self, llm_client: LLMClient) -> None:
        """Initialize question service."""
        self.llm_client = llm_client

    async def generate_question(
        self, request: GenerateQuestionRequest
    ) -> Question:
        """
        Generate question based on request using LLM.

        Args:
            request: Question generation request

        Returns:
            Generated question

        Raises:
            ValueError: If request validation fails or LLM response is invalid
        """
        # Request validation is done automatically by Pydantic
        logger.info(
            "generating question",
            language=request.language.value,
            topic=request.topic,
            difficulty=request.difficulty.value,
            question_type=request.question_type.value,
        )

        # Build prompt
        prompt = build_prompt(request)

        # Call LLM API
        try:
            llm_response = await self.llm_client.generate_question(prompt)
        except Exception as e:
            logger.error("failed to generate question from API", error=str(e))
            raise ValueError(
                f"failed to generate question from API: {e}"
            ) from e

        logger.debug(
            "received response from API",
            response_length=len(str(llm_response)),
        )

        # Validate LLM response matches request
        self._validate_llm_response(llm_response, request)

        # Convert LLM response to Question
        try:
            question = llm_response.to_question()
        except Exception as e:
            logger.error(
                "failed to convert LLM response to question", error=str(e)
            )
            raise ValueError(
                f"failed to convert LLM response to question: {e}"
            ) from e

        logger.info("question generated successfully", question_id=question.id)

        return question

    def _validate_llm_response(
        self,
        llm_response: LLMQuestionResponse,
        request: GenerateQuestionRequest,
    ) -> None:
        """
        Validate that LLM response matches the request.

        Args:
            llm_response: Response from LLM
            request: Original request

        Raises:
            ValueError: If response doesn't match request
        """
        if llm_response.language != request.language:
            raise ValueError(
                f"language mismatch: expected {request.language.value}, "
                f"got {llm_response.language.value}"
            )

        if llm_response.topic != request.topic:
            raise ValueError(
                f"topic mismatch: expected {request.topic}, "
                f"got {llm_response.topic}"
            )

        if llm_response.difficulty != request.difficulty:
            raise ValueError(
                f"difficulty mismatch: expected {request.difficulty.value}, "
                f"got {llm_response.difficulty.value}"
            )

        if llm_response.question_type != request.question_type:
            raise ValueError(
                f"question type mismatch: expected "
                f"{request.question_type.value}, "
                f"got {llm_response.question_type.value}"
            )

