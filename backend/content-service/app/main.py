"""Main FastAPI application."""

from contextlib import asynccontextmanager

from app.clients import ProxyAPIClient
from app.config import load_config
from app.exceptions import BusinessLogicError
from app.routers import questions_router
from fastapi import FastAPI, Request, status
from fastapi.responses import JSONResponse
from loguru import logger


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Lifespan context manager for app initialization."""
    # Load configuration
    try:
        config = load_config()
        logger.info(
            "Config loaded successfully",
            provider="ProxyAPI",
            model=config.proxyapi.model,
            base_url=config.proxyapi.base_url,
            timeout=config.proxyapi.timeout,
            max_tokens=config.proxyapi.max_tokens,
        )
    except Exception as e:
        logger.error(f"Failed to load config: {e}")
        raise

    # Check ProxyAPI configuration
    if not config.proxyapi.api_key:
        logger.error("ProxyAPI configuration is required")
        raise ValueError(
            "ProxyAPI configuration is required. "
            "Please set PROXYAPI_API_KEY and PROXYAPI_MODEL "
            "environment variables"
        )

    # Create LLM client and store in app state
    llm_client = ProxyAPIClient(config.proxyapi)
    app.state.proxyapi_client = llm_client

    logger.info("Application started successfully")

    yield

    # Close the client on shutdown
    await llm_client.close()
    logger.info("Application shutting down")


# Create FastAPI app
app = FastAPI(
    title="Snippet War Platform Service",
    description="Platform service for generating questions for Snippet War",
    version="0.1.0",
    lifespan=lifespan,
)

# Include routers
app.include_router(questions_router)


# Exception handlers
@app.exception_handler(BusinessLogicError)
async def business_logic_error_handler(
    request: Request, exc: BusinessLogicError
) -> JSONResponse:
    """Handle business logic errors."""
    return JSONResponse(
        status_code=exc.status_code,
        content={"error": exc.detail},
    )


@app.exception_handler(Exception)
async def general_exception_handler(request: Request, exc: Exception) -> JSONResponse:
    """Handle general exceptions."""
    logger.error(f"Unhandled exception: {exc}", exc_info=True)
    return JSONResponse(
        status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
        content={"error": "Internal server error"},
    )


@app.get("/health")
async def health_check() -> dict[str, str]:
    """Health check endpoint."""
    return {"status": "ok"}
