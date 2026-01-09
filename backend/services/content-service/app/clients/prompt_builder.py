"""Prompt builder for question generation."""

from dataclasses import dataclass

from app.models.question import GenerateQuestionRequest
from app.models.topics import get_topics_for_language


@dataclass
class RequestParams:
    """Parameters extracted from question generation request."""

    language_name: str
    language_id: str
    topic_name: str
    topic_id: str
    difficulty: str
    difficulty_description: str
    answer_type: str


def _get_request_params(req: GenerateQuestionRequest) -> RequestParams:
    """
    Extract common parameters from request for prompt building.

    Args:
        req: Question generation request

    Returns:
        RequestParams with extracted parameters
    """
    language_name = req.language.get_name()
    language_id = req.language.value

    topic_id = req.topic
    topic_name = str(topic_id)
    topics = get_topics_for_language(req.language)
    for topic in topics:
        if topic.topic_id == req.topic:
            topic_name = topic.name
            break

    difficulty = req.difficulty.value
    difficulty_description = req.difficulty.get_description()
    answer_type = req.question_type.value

    return RequestParams(
        language_name=language_name,
        language_id=language_id,
        topic_name=topic_name,
        topic_id=topic_id,
        difficulty=difficulty,
        difficulty_description=difficulty_description,
        answer_type=answer_type,
    )


def build_prompt(req: GenerateQuestionRequest) -> str:
    """
    Build prompt for single question generation.

    Args:
        req: Question generation request

    Returns:
        Prompt string for LLM
    """
    params = _get_request_params(req)

    prompt_template = f"""Ты - эксперт по программированию на языке {params.language_name}. Твоя задача - создать вопрос для игры "Snippet War", где игроки угадывают вывод кода.

Параметры вопроса:
- Язык программирования: {params.language_name}
- Тема: {params.topic_name}
- Уровень сложности: {params.difficulty} ({params.difficulty_description})
- Тип ответа: {params.answer_type}

Требования к вопросу:
1. Создай короткий, но понятный фрагмент кода (не более 30 строк), если это необходимо для вопроса
2. Код должен демонстрировать концепцию из указанной темы
3. Сложность должна соответствовать уровню {params.difficulty}
4. Вопрос должен быть интересным и обучающим
5. Код должен быть валидным и исполняемым (если код присутствует)
6. Теоретически вопрос может быть без кода - в таком случае поле "code" должно быть пустой строкой "", но ключ "code" всегда должен присутствовать в JSON ответе

Требования к ответу:
- Если тип ответа "multiple_choice": создай 4 варианта ответа, где правильных ответов может быть несколько. Варианты должны быть правдоподобными, но отличаться от правильных. В поле "correct_answers" укажи массив с ТОЧНЫМИ ТЕКСТАМИ КАЖДОГО правильного варианта из массива "options" (не индексы, а сами тексты ответов).
- Если тип ответа "free_text": в поле "correct_answers" укажи массив всех ПОДХОДЯЩИХ ВАРИАНТОВ правильного ответа (т.к. правильный ответ можно сформулировать по-разному, с учетом регистра, пробелов и т.д.)

Формат ответа (JSON):
{{
  "code": "фрагмент кода на языке {params.language_name}", // если кода нет - пустая строка "", но ключ всегда должен быть
  "question": "Что выведет этот код? (или другой вопрос)",
  "question_type": "{params.answer_type}",
  "options": ["вариант1", "вариант2", "вариант3", "вариант4"], // только для multiple_choice
  "correct_answers": ["точный текст правильного ответа"], // для multiple_choice - массив ТОЧНЫХ ТЕКСТОВ из options, для free_text - массив всех подходящих вариантов
  "explanation": "подробное объяснение, почему правильный ответ именно такой. Объясни, как работает код, какие концепции демонстрируются.",
  "difficulty": "{params.difficulty}",
  "topic": "{params.topic_id}",
  "language": "{params.language_id}"
}}

Важно:
- Код должен быть правильно отформатирован (если код присутствует)
- Если кода нет, поле "code" должно быть пустой строкой "", но ключ "code" всегда должен присутствовать в JSON
- Объяснение должно быть понятным и обучающим
- Для уровня Beginner код должен быть простым и понятным
- Для уровня Advanced можно использовать неочевидные особенности языка
- Всегда возвращай валидный JSON без дополнительного текста"""

    return prompt_template


def build_prompt_multiple(req: GenerateQuestionRequest, count: int) -> str:
    """
    Build prompt for multiple questions generation.

    Args:
        req: Question generation request
        count: Number of questions to generate

    Returns:
        Prompt string for LLM
    """
    params = _get_request_params(req)

    prompt_template = f"""Ты - эксперт по программированию на языке {params.language_name}. Твоя задача - создать {count} вопросов для игры "Snippet War", где игроки угадывают вывод кода.

Параметры вопросов:
- Язык программирования: {params.language_name}
- Тема: {params.topic_name}
- Уровень сложности: {params.difficulty} ({params.difficulty_description})
- Тип ответа: {params.answer_type}
- Количество вопросов: {count}

Требования к вопросам:
1. Создай {count} РАЗНЫХ вопросов по указанной теме
2. Каждый вопрос должен быть уникальным и демонстрировать разные аспекты темы
3. Для каждого вопроса создай короткий, но понятный фрагмент кода (не более 30 строк), если это необходимо для вопроса
4. Код должен демонстрировать концепцию из указанной темы
5. Сложность должна соответствовать уровню {params.difficulty}
6. Вопросы должны быть интересными и обучающими
7. Код должен быть валидным и исполняемым (если код присутствует)
8. Теоретически вопрос может быть без кода - в таком случае поле "code" должно быть пустой строкой "", но ключ "code" всегда должен присутствовать в JSON ответе

Требования к ответу:
- Если тип ответа "multiple_choice": для каждого вопроса создай 4 варианта ответа, где правильных ответов может быть несколько. Варианты должны быть правдоподобными, но отличаться от правильных. В поле "correct_answers" укажи массив с ТОЧНЫМИ ТЕКСТАМИ КАЖДОГО правильного варианта из массива "options" (не индексы, а сами тексты ответов).
- Если тип ответа "free_text": в поле "correct_answers" укажи массив всех ПОДХОДЯЩИХ ВАРИАНТОВ правильного ответа (т.к. правильный ответ можно сформулировать по-разному, с учетом регистра, пробелов и т.д.)

Формат ответа (JSON):
{{
  "questions": [
    {{
      "code": "фрагмент кода на языке {params.language_name}", // если кода нет - пустая строка "", но ключ всегда должен быть
      "question": "Что выведет этот код? (или другой вопрос)",
      "question_type": "{params.answer_type}",
      "options": ["вариант1", "вариант2", "вариант3", "вариант4"], // только для multiple_choice
      "correct_answers": ["точный текст правильного ответа"], // для multiple_choice - массив ТОЧНЫХ ТЕКСТОВ из options, для free_text - массив всех подходящих вариантов
      "explanation": "подробное объяснение, почему правильный ответ именно такой. Объясни, как работает код, какие концепции демонстрируются.",
      "difficulty": "{params.difficulty}",
      "topic": "{params.topic_id}",
      "language": "{params.language_id}"
    }}
    // ... еще {count - 1} вопросов
  ]
}}

Важно:
- Код должен быть правильно отформатирован (если код присутствует)
- Если кода нет, поле "code" должно быть пустой строкой "", но ключ "code" всегда должен присутствовать в JSON
- Объяснение должно быть понятным и обучающим
- Для уровня Beginner код должен быть простым и понятным
- Для уровня Advanced можно использовать неочевидные особенности языка
- Все вопросы должны быть РАЗНЫМИ и уникальными
- Всегда возвращай валидный JSON без дополнительного текста
- Массив "questions" должен содержать ровно {count} вопросов"""

    return prompt_template
