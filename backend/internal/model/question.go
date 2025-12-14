package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Question представляет вопрос для игры Snippet War.
type Question struct {
	ID                 string       `json:"id"`
	Language           Language     `json:"language"`
	Topic              TopicID      `json:"topic"`
	Difficulty         Difficulty   `json:"difficulty"`
	QuestionType       QuestionType `json:"question_type"`
	Code               string       `json:"code"`
	QuestionText       string       `json:"question"`
	Options            []string     `json:"options,omitempty"`             // для multiple_choice
	CorrectAnswer      string       `json:"correct_answer"`                // текст правильного ответа (для multiple_choice - текст из options, для free_text - правильный ответ)
	AcceptableVariants []string     `json:"acceptable_variants,omitempty"` // для free_text
	CaseSensitive      bool         `json:"case_sensitive,omitempty"`      // для free_text
	Explanation        string       `json:"explanation"`
	CreatedAt          time.Time    `json:"created_at"`
}

// GenerateID генерирует новый UUID для вопроса.
func (q *Question) GenerateID() {
	q.ID = uuid.New().String()
}

// Validate проверяет валидность структуры Question.
func (q *Question) Validate() error {
	if q.ID == "" {
		return errors.New("question ID is required")
	}

	if q.Language == "" {
		return errors.New("language is required")
	}

	if !q.Language.IsValid() {
		return fmt.Errorf("invalid language: %s", q.Language)
	}

	if q.Topic == "" {
		return errors.New("topic is required")
	}

	if !IsValidTopic(q.Language, q.Topic) {
		return fmt.Errorf("invalid topic '%s' for language '%s'", q.Topic, q.Language)
	}

	if !q.Difficulty.IsValid() {
		return fmt.Errorf("invalid difficulty: %s", q.Difficulty)
	}

	if !q.QuestionType.IsValid() {
		return fmt.Errorf("invalid question type: %s", q.QuestionType)
	}

	if q.Code == "" {
		return errors.New("code is required")
	}

	if q.QuestionText == "" {
		return errors.New("question text is required")
	}

	if q.CorrectAnswer == "" {
		return errors.New("correct answer is required")
	}

	if q.Explanation == "" {
		return errors.New("explanation is required")
	}

	// Валидация в зависимости от типа вопроса
	switch q.QuestionType {
	case QuestionTypeMultipleChoice:
		if err := ValidateMultipleChoiceOptions(q.Options); err != nil {
			return err
		}

		if err := ValidateMultipleChoiceAnswerString(q.CorrectAnswer, q.Options); err != nil {
			return err
		}

	case QuestionTypeFreeText:
		if err := ValidateFreeTextAnswerString(q.CorrectAnswer); err != nil {
			return err
		}

	default:
		return fmt.Errorf("unknown question type: %s", q.QuestionType)
	}

	return nil
}
