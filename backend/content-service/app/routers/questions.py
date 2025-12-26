"""Question generation API routes."""

from app.clients import ProxyAPIClient
from app.exceptions import BusinessLogicError, ValidationError
from app.models import GenerateQuestionRequest, Question
from app.services import QuestionService
from fastapi import APIRouter, Depends, HTTPException, Request, status


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
) -> Question:
    """
    Generate a question based on request parameters.

    Args:
        request: Question generation request
        service: Question service dependency

    Returns:
        Generated question

    Raises:
        HTTPException: If generation fails
    """
    try:
        return await service.generate_question(request)
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
