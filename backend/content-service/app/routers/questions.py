"""Question generation API routes."""

from app.clients import ProxyAPIClient
from app.database import get_db_session
from app.exceptions import BusinessLogicError, ValidationError
from app.models import GenerateQuestionRequest, GetQuestionsBatchRequest, Question
from app.services import QuestionService
from fastapi import APIRouter, Depends, HTTPException, Request, status
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


@router.post(
    "/batch",
    response_model=list[Question],
    status_code=status.HTTP_200_OK,
)
async def get_questions_batch(
    request: GetQuestionsBatchRequest,
    service: QuestionService = Depends(get_question_service),
    db: AsyncSession = Depends(get_db_session),
) -> list[Question]:
    """
    Get batch of questions, generating missing ones if needed.

    Args:
        request: Batch request parameters
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
            language=request.language,
            topic=request.topic,
            difficulty=request.difficulty,
            count=request.count,
            question_type=request.question_type,
            telegram_user_id=request.telegram_user_id,
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
