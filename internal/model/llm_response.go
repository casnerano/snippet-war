package model

import (
	"errors"
	"fmt"
	"time"
)

// LLMQuestionResponse представляет ответ от LLM для генерации вопроса.
type LLMQuestionResponse struct {
	Code               string       `json:"code"`
	Question           string       `json:"question"`
	QuestionType       QuestionType `json:"question_type"`
	Options            []string     `json:"options,omitempty"`
	CorrectAnswer      interface{}  `json:"correct_answer"` // строка с текстом правильного ответа (для multiple_choice - текст из options)
	AcceptableVariants []string     `json:"acceptable_variants,omitempty"`
	CaseSensitive      bool         `json:"case_sensitive,omitempty"`
	Explanation        string       `json:"explanation"`
	Difficulty         Difficulty   `json:"difficulty"`
	Topic              TopicID      `json:"topic"`
	Language           Language     `json:"language"`
}

// Validate проверяет валидность ответа от LLM.
func (r *LLMQuestionResponse) Validate() error {
	if r.Code == "" {
		return errors.New("code is required")
	}

	if r.Question == "" {
		return errors.New("question is required")
	}

	if !r.QuestionType.IsValid() {
		return fmt.Errorf("invalid question type: %s", r.QuestionType)
	}

	if !r.Difficulty.IsValid() {
		return fmt.Errorf("invalid difficulty: %s", r.Difficulty)
	}

	if !r.Language.IsValid() {
		return fmt.Errorf("invalid language: %s", r.Language)
	}

	if r.Topic == "" {
		return errors.New("topic is required")
	}

	if !IsValidTopic(r.Language, r.Topic) {
		return fmt.Errorf("invalid topic '%s' for language '%s'", r.Topic, r.Language)
	}

	if r.Explanation == "" {
		return errors.New("explanation is required")
	}

	if r.CorrectAnswer == nil {
		return errors.New("correct answer is required")
	}

	// Валидация в зависимости от типа вопроса
	switch r.QuestionType {
	case QuestionTypeMultipleChoice:
		if err := ValidateMultipleChoiceOptions(r.Options); err != nil {
			return err
		}

		if _, err := ValidateMultipleChoiceAnswer(r.CorrectAnswer, r.Options); err != nil {
			return err
		}

	case QuestionTypeFreeText:
		if _, err := ValidateFreeTextAnswer(r.CorrectAnswer); err != nil {
			return err
		}

	default:
		return fmt.Errorf("unknown question type: %s", r.QuestionType)
	}

	return nil
}

// ToQuestion преобразует ответ от LLM в модель Question.
func (r *LLMQuestionResponse) ToQuestion() (*Question, error) {
	if err := r.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	question := &Question{
		Language:           r.Language,
		Topic:              r.Topic,
		Difficulty:         r.Difficulty,
		QuestionType:       r.QuestionType,
		Code:               r.Code,
		QuestionText:       r.Question,
		Explanation:        r.Explanation,
		CreatedAt:          time.Now(),
		AcceptableVariants: r.AcceptableVariants,
		CaseSensitive:      r.CaseSensitive,
	}

	// Генерируем ID
	question.GenerateID()

	// Обработка correct_answer в зависимости от типа вопроса
	switch r.QuestionType {
	case QuestionTypeMultipleChoice:
		question.Options = r.Options

		// Проверяем, что correct_answer является валидным текстом из options
		answer, err := ValidateMultipleChoiceAnswer(r.CorrectAnswer, r.Options)
		if err != nil {
			return nil, fmt.Errorf("failed to validate correct answer: %w", err)
		}

		question.CorrectAnswer = answer

	case QuestionTypeFreeText:
		// Для free_text correct_answer должен быть строкой
		answer, err := ValidateFreeTextAnswer(r.CorrectAnswer)
		if err != nil {
			return nil, fmt.Errorf("failed to validate correct answer: %w", err)
		}

		question.CorrectAnswer = answer
	}

	// Финальная валидация созданного вопроса
	if err := question.Validate(); err != nil {
		return nil, fmt.Errorf("question validation failed: %w", err)
	}

	return question, nil
}
