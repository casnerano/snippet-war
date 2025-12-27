"""Question generation service."""

import uuid

from app.clients import LLMClient, build_prompt
from app.models import GenerateQuestionRequest, LLMQuestionResponse, Question
from app.models.enums import Difficulty, Language, QuestionType
from app.repositories import QuestionRepository, UserRepository
from loguru import logger
from sqlalchemy.ext.asyncio import AsyncSession


class QuestionService:
    """Service for generating questions using LLM."""

    def __init__(self, llm_client: LLMClient) -> None:
        """Initialize question service."""
        self.llm_client = llm_client

    async def generate_question(
        self,
        request: GenerateQuestionRequest,
        db_session: AsyncSession | None = None,
    ) -> Question:
        """
        Generate question based on request using LLM.

        Args:
            request: Question generation request
            db_session: Optional database session to save question

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
            raise ValueError(f"failed to generate question from API: {e}") from e

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
            logger.error("failed to convert LLM response to question", error=str(e))
            raise ValueError(f"failed to convert LLM response to question: {e}") from e

        # Save to database if session provided
        if db_session is not None:
            try:
                await QuestionRepository.save_question(db_session, question)
                logger.info("question saved to database", question_id=question.id)
            except Exception as e:
                logger.warning(
                    "failed to save question to database",
                    question_id=question.id,
                    error=str(e),
                )
                # Don't fail the request if save fails, just log warning

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
                f"topic mismatch: expected {request.topic}, got {llm_response.topic}"
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

    async def get_questions_batch(
        self,
        db_session: AsyncSession,
        language: Language,
        topic: str,
        difficulty: Difficulty,
        count: int,
        question_type: QuestionType = QuestionType.MULTIPLE_CHOICE,
        telegram_user_id: int | None = None,
    ) -> list[Question]:
        """
        Get batch of questions, generating missing ones if needed.

        Args:
            db_session: Database session
            language: Programming language
            topic: Topic
            difficulty: Difficulty level
            count: Number of questions to return
            telegram_user_id: Optional Telegram user ID

        Returns:
            List of questions

        Raises:
            ValueError: If count is invalid or generation fails
        """
        if count <= 0:
            raise ValueError("count must be greater than 0")

        logger.info(
            "getting questions batch",
            language=language.value,
            topic=topic,
            difficulty=difficulty.value,
            count=count,
            telegram_user_id=telegram_user_id,
        )

        user_id: uuid.UUID | None = None

        # Get or create user if telegram_user_id provided
        if telegram_user_id is not None:
            user = await UserRepository.get_or_create_user_by_telegram_id(
                db_session, telegram_user_id
            )
            user_id = user.id

        # Get existing questions
        if user_id is not None:
            existing_questions = await QuestionRepository.get_unseen_questions(
                db_session,
                user_id,
                language,
                topic,
                difficulty,
                limit=count,
            )
        else:
            existing_questions = await QuestionRepository.get_questions_by_filters(
                db_session,
                language,
                topic,
                difficulty,
                limit=count,
            )

        # Convert to Question models
        questions = [
            QuestionRepository.question_db_to_model(q) for q in existing_questions
        ]

        # Generate missing questions if needed
        missing_count = count - len(questions)
        if missing_count > 0:
            logger.info(
                "generating missing questions",
                missing_count=missing_count,
            )
            for _ in range(missing_count):
                request = GenerateQuestionRequest(
                    language=language,
                    topic=topic,
                    difficulty=difficulty,
                    question_type=question_type,
                )
                # Generate question with db_session to save it
                question = await self.generate_question(request, db_session=db_session)
                questions.append(question)

        # Take only requested count (in case we generated more)
        questions = questions[:count]

        # Mark questions as seen if user provided
        if user_id is not None and questions:
            question_ids = [uuid.UUID(q.id) for q in questions]
            await QuestionRepository.mark_questions_as_seen(
                db_session, user_id, question_ids
            )

        logger.info(
            "questions batch retrieved",
            count=len(questions),
            telegram_user_id=telegram_user_id,
        )

        return questions
