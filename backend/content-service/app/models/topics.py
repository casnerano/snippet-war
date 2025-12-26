"""Topic definitions for programming languages."""

from app.models.enums import Language
from pydantic import BaseModel


class Topic(BaseModel):
    """Topic model with ID and name."""

    topic_id: str
    name: str

    class Config:
        """Pydantic config."""

        populate_by_name = True


# Topic definitions for each language
LANGUAGE_TOPICS: dict[Language, list[Topic]] = {
    Language.PYTHON: [
        Topic(topic_id="variables_types", name="Переменные и типы данных"),
        Topic(topic_id="lists_arrays", name="Списки и массивы"),
        Topic(topic_id="dictionaries", name="Словари"),
        Topic(topic_id="functions", name="Функции"),
        Topic(topic_id="closures", name="Замыкания"),
        Topic(topic_id="decorators", name="Декораторы"),
        Topic(topic_id="generators", name="Генераторы"),
        Topic(topic_id="classes_oop", name="Классы и ООП"),
        Topic(topic_id="exceptions", name="Обработка исключений"),
        Topic(topic_id="context_managers", name="Контекстные менеджеры"),
        Topic(topic_id="async_await", name="Асинхронное программирование"),
    ],
    Language.JAVASCRIPT: [
        Topic(topic_id="variables_types", name="Переменные и типы"),
        Topic(topic_id="arrays", name="Массивы"),
        Topic(topic_id="objects", name="Объекты"),
        Topic(topic_id="functions", name="Функции"),
        Topic(topic_id="closures", name="Замыкания"),
        Topic(topic_id="this_binding", name="Контекст выполнения (this)"),
        Topic(topic_id="prototypes", name="Прототипы"),
        Topic(topic_id="classes", name="Классы (ES6+)"),
        Topic(topic_id="promises_async", name="Промисы и async/await"),
        Topic(topic_id="event_loop", name="Event Loop"),
        Topic(topic_id="destructuring", name="Деструктуризация"),
        Topic(topic_id="modules", name="Модули (ES6+)"),
    ],
    Language.GO: [
        Topic(topic_id="variables_types", name="Переменные и типы"),
        Topic(topic_id="slices", name="Срезы"),
        Topic(topic_id="maps", name="Мапы"),
        Topic(topic_id="functions", name="Функции"),
        Topic(topic_id="methods", name="Методы"),
        Topic(topic_id="interfaces", name="Интерфейсы"),
        Topic(topic_id="goroutines", name="Горутины"),
        Topic(topic_id="channels", name="Каналы"),
        Topic(topic_id="select", name="Select statement"),
        Topic(topic_id="defer_panic_recover", name="Defer, panic, recover"),
        Topic(topic_id="pointers", name="Указатели"),
        Topic(topic_id="structs", name="Структуры"),
    ],
    Language.JAVA: [
        Topic(topic_id="variables_types", name="Переменные и типы"),
        Topic(topic_id="arrays_lists", name="Массивы и списки"),
        Topic(topic_id="collections", name="Коллекции (Set, Map)"),
        Topic(topic_id="methods", name="Методы"),
        Topic(topic_id="classes_objects", name="Классы и объекты"),
        Topic(topic_id="inheritance", name="Наследование"),
        Topic(topic_id="interfaces", name="Интерфейсы"),
        Topic(topic_id="generics", name="Дженерики"),
        Topic(topic_id="exceptions", name="Исключения"),
        Topic(topic_id="streams", name="Streams API"),
        Topic(topic_id="lambda_expressions", name="Lambda выражения"),
        Topic(topic_id="concurrency", name="Многопоточность"),
    ],
    Language.CPP: [
        Topic(topic_id="variables_types", name="Переменные и типы"),
        Topic(topic_id="pointers_references", name="Указатели и ссылки"),
        Topic(topic_id="arrays_vectors", name="Массивы и векторы"),
        Topic(topic_id="functions", name="Функции"),
        Topic(topic_id="classes_objects", name="Классы и объекты"),
        Topic(topic_id="inheritance", name="Наследование"),
        Topic(topic_id="templates", name="Шаблоны"),
        Topic(topic_id="smart_pointers", name="Умные указатели"),
        Topic(topic_id="stl", name="STL контейнеры и алгоритмы"),
        Topic(topic_id="move_semantics", name="Move семантика"),
        Topic(topic_id="lambda", name="Lambda выражения"),
        Topic(topic_id="multithreading", name="Многопоточность"),
    ],
    Language.RUST: [
        Topic(topic_id="variables_types", name="Переменные и типы"),
        Topic(topic_id="ownership", name="Владение (ownership)"),
        Topic(topic_id="borrowing", name="Заимствование (borrowing)"),
        Topic(topic_id="lifetimes", name="Время жизни"),
        Topic(topic_id="vectors", name="Векторы"),
        Topic(topic_id="hashmaps", name="HashMap"),
        Topic(topic_id="functions", name="Функции"),
        Topic(topic_id="structs", name="Структуры"),
        Topic(topic_id="enums", name="Перечисления"),
        Topic(topic_id="pattern_matching", name="Сопоставление с образцом"),
        Topic(
            topic_id="error_handling",
            name="Обработка ошибок (Result, Option)",
        ),
        Topic(topic_id="concurrency", name="Многопоточность"),
    ],
    Language.TYPESCRIPT: [
        Topic(topic_id="types", name="Типы"),
        Topic(topic_id="interfaces", name="Интерфейсы"),
        Topic(topic_id="generics", name="Дженерики"),
        Topic(
            topic_id="unions_intersections",
            name="Объединения и пересечения типов",
        ),
        Topic(topic_id="type_guards", name="Защитники типов"),
        Topic(topic_id="decorators", name="Декораторы"),
        Topic(topic_id="utility_types", name="Утилитарные типы"),
        Topic(topic_id="modules", name="Модули"),
        Topic(topic_id="async_promises", name="Асинхронность"),
        Topic(topic_id="classes", name="Классы"),
        Topic(topic_id="namespaces", name="Пространства имен"),
    ],
}


def get_topics_for_language(language: Language) -> list[Topic]:
    """Get topics for a specific language."""
    return LANGUAGE_TOPICS.get(language, [])


def is_valid_topic(language: Language, topic: str) -> bool:
    """Check if topic is valid for the given language."""
    topics = LANGUAGE_TOPICS.get(language, [])
    return any(t.topic_id == topic for t in topics)
