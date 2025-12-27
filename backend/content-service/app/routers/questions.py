"""Question generation API routes."""

from app.clients import ProxyAPIClient
from app.database import get_db_session
from app.exceptions import BusinessLogicError, ValidationError
from app.models import GenerateQuestionRequest, Question
from app.models.enums import Difficulty, Language, QuestionType
from app.services import QuestionService
from fastapi import APIRouter, Depends, HTTPException, Query, Request, status
from sqlalchemy.ext.asyncio import AsyncSession


def get_question_service(request: Request) -> QuestionService:
    """Get question service dependency from app state."""
    llm_client: ProxyAPIClient = request.app.state.proxyapi_client
    return QuestionService(llm_client)


router = APIRouter(prefix="/api/questions", tags=["questions"])


@router.post(
    "/generate",
    response_model=Question,
    status_code=status.HTTP_200_OK,
)
async def generate_question(
    request: GenerateQuestionRequest,
    service: QuestionService = Depends(get_question_service),
    db: AsyncSession = Depends(get_db_session),
) -> Question:
    """
    Generate a question based on request parameters.

    Args:
        request: Question generation request
        service: Question service dependency
        db: Database session

    Returns:
        Generated question

    Raises:
        HTTPException: If generation fails
    """
    try:
        return await service.generate_question(request, db_session=db)
    except ValidationError as e:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST, detail=e.detail
        ) from e
    except BusinessLogicError as e:
        raise HTTPException(status_code=e.status_code, detail=e.detail) from e
    except ValueError as e:
        # Convert ValueError to ValidationError
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST, detail=str(e)
        ) from e
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Internal server error: {str(e)}",
        ) from e


@router.get(
    "/batch",
    response_model=list[Question],
    status_code=status.HTTP_200_OK,
)
async def get_questions_batch(
    language: Language = Query(..., description="Programming language"),
    topic: str = Query(..., description="Topic"),
    difficulty: Difficulty = Query(..., description="Difficulty level"),
    count: int = Query(..., ge=1, description="Number of questions"),
    question_type: QuestionType = Query(
        QuestionType.MULTIPLE_CHOICE, description="Question type"
    ),
    telegram_user_id: int | None = Query(
        None, description="Telegram user ID (optional)"
    ),
    service: QuestionService = Depends(get_question_service),
    db: AsyncSession = Depends(get_db_session),
) -> list[Question]:
    """
    Get batch of questions, generating missing ones if needed.

    Args:
        language: Programming language
        topic: Topic
        difficulty: Difficulty level
        count: Number of questions to return
        question_type: Question type
        telegram_user_id: Optional Telegram user ID
        service: Question service dependency
        db: Database session

    Returns:
        List of questions

    Raises:
        HTTPException: If request fails
    """
    try:
        return await service.get_questions_batch(
            db_session=db,
            language=language,
            topic=topic,
            difficulty=difficulty,
            count=count,
            question_type=question_type,
            telegram_user_id=telegram_user_id,
        )
    except ValidationError as e:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST, detail=e.detail
        ) from e
    except BusinessLogicError as e:
        raise HTTPException(status_code=e.status_code, detail=e.detail) from e
    except ValueError as e:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST, detail=str(e)
        ) from e
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Internal server error: {str(e)}",
        ) from e
