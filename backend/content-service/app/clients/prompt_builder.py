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
1. Создай короткий, но понятный фрагмент кода (не более 30 строк)
2. Код должен демонстрировать концепцию из указанной темы
3. Сложность должна соответствовать уровню {difficulty}
4. Вопрос должен быть интересным и обучающим
5. Код должен быть валидным и исполняемым

Требования к ответу:
- Если тип ответа "multiple_choice": создай 4 варианта ответа, где только один правильный. Варианты должны быть правдоподобными, но отличаться от правильного. В поле "correct_answer" укажи ТОЧНЫЙ ТЕКСТ правильного варианта из массива "options" (не индекс, а сам текст ответа).
- Если тип ответа "free_text": укажи точный правильный ответ и возможные варианты написания (с учетом регистра, пробелов и т.д.)

Формат ответа (JSON):
{{
  "code": "фрагмент кода на языке {language_name}",
  "question": "Что выведет этот код? (или другой вопрос)",
  "question_type": "{answer_type}",
  "options": ["вариант1", "вариант2", "вариант3", "вариант4"], // только для multiple_choice
  "correct_answer": "точный текст правильного ответа из options", // для multiple_choice - ТОЧНЫЙ ТЕКСТ из options, для free_text - правильный ответ
  "acceptable_variants": ["вариант1", "вариант2"], // только для free_text, опционально
  "case_sensitive": false, // только для free_text
  "explanation": "подробное объяснение, почему правильный ответ именно такой. Объясни, как работает код, какие концепции демонстрируются.",
  "difficulty": "{difficulty}",
  "topic": "{topic_id}",
  "language": "{language_id}"
}}

Важно:
- Код должен быть правильно отформатирован
- Объяснение должно быть понятным и обучающим
- Для уровня Beginner код должен быть простым и понятным
- Для уровня Advanced можно использовать неочевидные особенности языка
- Всегда возвращай валидный JSON без дополнительного текста"""

    return prompt_template
