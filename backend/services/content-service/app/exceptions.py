"""Custom exceptions for the application."""


class BusinessLogicError(Exception):
    """Base exception for business logic errors."""

    def __init__(self, detail: str, status_code: int = 500) -> None:
        """Initialize business logic error."""
        self.detail = detail
        self.status_code = status_code
        super().__init__(self.detail)


class ValidationError(BusinessLogicError):
    """Validation error (400)."""

    def __init__(self, detail: str) -> None:
        """Initialize validation error."""
        super().__init__(detail, status_code=400)


class DatabaseError(BusinessLogicError):
    """Database operation error (500)."""

    def __init__(self, detail: str) -> None:
        """Initialize database error."""
        super().__init__(detail, status_code=500)
