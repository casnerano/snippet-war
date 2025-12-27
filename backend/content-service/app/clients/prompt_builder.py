"""Prompt builder for question generation."""

from app.models.question import GenerateQuestionRequest
from app.models.topics import get_topics_for_language


def build_prompt(req: GenerateQuestionRequest) -> str:
    """Build prompt for question generation based on request parameters."""
    # Get language name
    language_name = req.language.get_name()
    language_id = req.language.value

    # Get topic name and ID
    topic_id = req.topic
    topic_name = str(topic_id)
    topics = get_topics_for_language(req.language)
    for topic in topics:
        if topic.topic_id == req.topic:
            topic_name = topic.name
            break

    # Get difficulty level and description
    difficulty = req.difficulty.value
    difficulty_description = req.difficulty.get_description()

    # Get answer type
    answer_type = req.question_type.value

    # Prompt template
    prompt_template = f"""Ты - эксперт по программированию на языке {language_name}. Твоя задача - создать вопрос для игры "Snippet War", где игроки угадывают вывод кода.

Параметры вопроса:
- Язык программирования: {language_name}
- Тема: {topic_name}
- Уровень сложности: {difficulty} ({difficulty_description})
- Тип ответа: {answer_type}

Требования к вопросу:
1. Создай короткий, но понятный фрагмент кода (не более 30 строк), если это необходимо для вопроса
2. Код должен демонстрировать концепцию из указанной темы
3. Сложность должна соответствовать уровню {difficulty}
4. Вопрос должен быть интересным и обучающим
5. Код должен быть валидным и исполняемым (если код присутствует)
6. Теоретически вопрос может быть без кода - в таком случае поле "code" должно быть пустой строкой "", но ключ "code" всегда должен присутствовать в JSON ответе

Требования к ответу:
- Если тип ответа "multiple_choice": создай 4 варианта ответа, где правильных ответов может быть несколько. Варианты должны быть правдоподобными, но отличаться от правильных. В поле "correct_answers" укажи массив с ТОЧНЫМИ ТЕКСТАМИ КАЖДОГО правильного варианта из массива "options" (не индексы, а сами тексты ответов).
- Если тип ответа "free_text": в поле "correct_answers" укажи массив всех ПОДХОДЯЩИХ ВАРИАНТОВ правильного ответа (т.к. правильный ответ можно сформулировать по-разному, с учетом регистра, пробелов и т.д.)

Формат ответа (JSON):
{{
  "code": "фрагмент кода на языке {language_name}", // если кода нет - пустая строка "", но ключ всегда должен быть
  "question": "Что выведет этот код? (или другой вопрос)",
  "question_type": "{answer_type}",
  "options": ["вариант1", "вариант2", "вариант3", "вариант4"], // только для multiple_choice
  "correct_answers": ["точный текст правильного ответа"], // для multiple_choice - массив ТОЧНЫХ ТЕКСТОВ из options, для free_text - массив всех подходящих вариантов
  "explanation": "подробное объяснение, почему правильный ответ именно такой. Объясни, как работает код, какие концепции демонстрируются.",
  "difficulty": "{difficulty}",
  "topic": "{topic_id}",
  "language": "{language_id}"
}}

Важно:
- Код должен быть правильно отформатирован (если код присутствует)
- Если кода нет, поле "code" должно быть пустой строкой "", но ключ "code" всегда должен присутствовать в JSON
- Объяснение должно быть понятным и обучающим
- Для уровня Beginner код должен быть простым и понятным
- Для уровня Advanced можно использовать неочевидные особенности языка
- Всегда возвращай валидный JSON без дополнительного текста"""

    return prompt_template
