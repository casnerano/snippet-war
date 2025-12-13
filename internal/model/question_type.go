package model

// QuestionType представляет тип вопроса.
type QuestionType string

const (
	// QuestionTypeMultipleChoice - вопрос с множественным выбором.
	QuestionTypeMultipleChoice QuestionType = "multiple_choice"
	// QuestionTypeFreeText - вопрос со свободным текстовым ответом.
	QuestionTypeFreeText QuestionType = "free_text"
)

// AllQuestionTypes возвращает список всех поддерживаемых типов вопросов.
func AllQuestionTypes() []QuestionType {
	return []QuestionType{
		QuestionTypeMultipleChoice,
		QuestionTypeFreeText,
	}
}

// String возвращает строковое представление типа вопроса.
func (qt QuestionType) String() string {
	return string(qt)
}

// IsValid проверяет, является ли значение валидным типом вопроса.
func (qt QuestionType) IsValid() bool {
	switch qt {
	case QuestionTypeMultipleChoice, QuestionTypeFreeText:
		return true
	default:
		return false
	}
}
