"""Enums for question generation models."""

from enum import Enum


class Language(str, Enum):
    """Supported programming languages."""

    PYTHON = "python"
    JAVASCRIPT = "javascript"
    GO = "go"
    JAVA = "java"
    CPP = "cpp"
    RUST = "rust"
    TYPESCRIPT = "typescript"

    def get_name(self) -> str:
        """Get display name for the language."""
        names = {
            Language.PYTHON: "Python",
            Language.JAVASCRIPT: "JavaScript",
            Language.GO: "Go",
            Language.JAVA: "Java",
            Language.CPP: "C++",
            Language.RUST: "Rust",
            Language.TYPESCRIPT: "TypeScript",
        }
        return names.get(self, self.value)


class Difficulty(str, Enum):
    """Question difficulty levels."""

    BEGINNER = "beginner"
    INTERMEDIATE = "intermediate"
    ADVANCED = "advanced"

    def get_description(self) -> str:
        """Get description for the difficulty level."""
        descriptions = {
            Difficulty.BEGINNER: (
                "Базовые операции и синтаксис. Простые типы данных, "
                "базовые структуры данных, простые условия и циклы, "
                "простые функции без сложной логики, базовые операции "
                "со строками и числами."
            ),
            Difficulty.INTERMEDIATE: (
                "Более сложные структуры данных, вложенные циклы и условия, "
                "функции высшего порядка, работа с коллекциями, базовое ООП, "
                "обработка исключений, базовые паттерны проектирования."
            ),
            Difficulty.ADVANCED: (
                "Сложные алгоритмы и оптимизация, продвинутые концепции языка, "
                "конкурентное/параллельное программирование, продвинутые паттерны "
                "проектирования, неочевидное поведение языка, оптимизация "
                "производительности, работа с памятью и указателями."
            ),
        }
        return descriptions.get(self, "")


class QuestionType(str, Enum):
    """Question types."""

    MULTIPLE_CHOICE = "multiple_choice"
    FREE_TEXT = "free_text"
