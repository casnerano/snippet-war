"""Validation functions for question models."""


def validate_multiple_choice_options(options: list[str]) -> None:
    """Validate options for multiple choice question."""
    if len(options) < 2:
        raise ValueError("multiple choice question must have at least 2 options")
    if len(options) > 5:
        raise ValueError("multiple choice question must have at most 5 options")


def validate_multiple_choice_answer(
    correct_answers: list[str], options: list[str]
) -> None:
    """Validate correct answer for multiple choice question."""
    if not correct_answers or len(correct_answers) == 0:
        raise ValueError("correct answer must be a non-empty list for multiple choice")

    for answer in correct_answers:
        if answer not in options:
            raise ValueError(
                f"correct answer '{answer}' must be one of the options: {options}"
            )


def validate_free_text_answer(correct_answers: list[str]) -> None:
    """Validate correct answer for free text question."""
    if not correct_answers or len(correct_answers) == 0:
        raise ValueError("correct answer must be a non-empty list for free text")
