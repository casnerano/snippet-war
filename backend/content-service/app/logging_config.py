"""Logging configuration for the application."""

import sys
from pathlib import Path

from loguru import logger


def setup_logging(
    log_level: str = "INFO",
    log_file: str | None = None,
    log_dir: Path | None = None,
    enable_json: bool = True,
    enable_console: bool = True,
) -> None:
    """
    Configure loguru logger.

    Args:
        log_level: Logging level (DEBUG, INFO, WARNING, ERROR, CRITICAL)
        log_file: Log file name (default: app.log)
        log_dir: Directory for log files (default: ./logs)
        enable_json: Enable JSON format for file logs
        enable_console: Enable console output
    """
    # Remove default handler
    logger.remove()

    # Console handler with colored output
    # Note: Named arguments are automatically added to extra dict
    # and will be visible in JSON file logs with serialize=True
    if enable_console:
        logger.add(
            sys.stderr,
            format=(
                "<green>{time:YYYY-MM-DD HH:mm:ss}</green> | "
                "<level>{level: <8}</level> | "
                "<cyan>{name}</cyan>:<cyan>{function}</cyan>:<cyan>{line}</cyan> | "
                "<level>{message}</level>"
            ),
            level=log_level,
            colorize=True,
        )

    # File handler with structured JSON format
    if log_file or log_dir:
        if log_dir is None:
            log_dir = Path("./logs")
        else:
            log_dir = Path(log_dir)

        # Create log directory if it doesn't exist
        log_dir.mkdir(parents=True, exist_ok=True)

        if log_file is None:
            log_file = "app.log"

        log_path = log_dir / log_file

        # JSON format for structured logging (includes all named arguments)
        if enable_json:
            logger.add(
                log_path,
                format="{time} | {level} | {name}:{function}:{line} | {message} | {extra}",
                level=log_level,
                rotation="10 MB",
                retention="7 days",
                compression="zip",
                serialize=True,  # JSON format
                encoding="utf-8",
                backtrace=True,
                diagnose=True,
            )
        else:
            # Human-readable format for file
            logger.add(
                log_path,
                format=(
                    "{time:YYYY-MM-DD HH:mm:ss} | {level: <8} | "
                    "{name}:{function}:{line} | {message} | {extra}"
                ),
                level=log_level,
                rotation="10 MB",
                retention="7 days",
                compression="zip",
                encoding="utf-8",
                backtrace=True,
                diagnose=True,
            )

        logger.info(
            "Logging configured",
            log_file=str(log_path),
            log_level=log_level,
            json_format=enable_json,
        )

