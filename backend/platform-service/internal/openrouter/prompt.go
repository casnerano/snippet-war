package openrouter

import (
	"strings"

	"github.com/casnerano/snippet-war/internal/model"
)

// BuildPrompt создает промпт для генерации вопроса на основе параметров запроса.
func BuildPrompt(req *model.GenerateQuestionRequest) string {
	// Получаем название языка
	languageName := req.Language.GetName()
	languageID := req.Language.String()

	// Получаем название темы и ID
	topicName := req.Topic.String()
	topicID := req.Topic.String()
	topics := model.GetTopicsForLanguage(req.Language)
	for _, topic := range topics {
		if topic.ID == req.Topic {
			topicName = topic.Name
			break
		}
	}

	// Получаем уровень сложности и его описание
	difficulty := req.Difficulty.String()
	difficultyDescription := req.Difficulty.GetDescription()

	// Получаем тип ответа
	answerType := req.QuestionType.String()

	// Шаблон промпта
	promptTemplate := `Ты - эксперт по программированию на языке {language}. Твоя задача - создать вопрос для игры "Snippet War", где игроки угадывают вывод кода.

Параметры вопроса:
- Язык программирования: {language}
- Тема: {topic}
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
{
  "code": "фрагмент кода на языке {language}",
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
}

Важно:
- Код должен быть правильно отформатирован
- Объяснение должно быть понятным и обучающим
- Для уровня Beginner код должен быть простым и понятным
- Для уровня Advanced можно использовать неочевидные особенности языка
- Всегда возвращай валидный JSON без дополнительного текста`

	// Подставляем параметры
	prompt := promptTemplate
	prompt = strings.ReplaceAll(prompt, "{language}", languageName)
	prompt = strings.ReplaceAll(prompt, "{topic}", topicName)
	prompt = strings.ReplaceAll(prompt, "{difficulty}", difficulty)
	prompt = strings.ReplaceAll(prompt, "{difficulty_description}", difficultyDescription)
	prompt = strings.ReplaceAll(prompt, "{answer_type}", answerType)
	prompt = strings.ReplaceAll(prompt, "{topic_id}", topicID)
	prompt = strings.ReplaceAll(prompt, "{language_id}", languageID)

	return prompt
}
