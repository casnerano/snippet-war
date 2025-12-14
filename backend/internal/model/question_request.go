package model

import (
	"errors"
	"fmt"
)

// GenerateQuestionRequest представляет запрос на генерацию вопроса.
type GenerateQuestionRequest struct {
	Language     Language     `json:"language" binding:"required"`
	Topic        TopicID      `json:"topic" binding:"required"`
	Difficulty   Difficulty   `json:"difficulty" binding:"required"`
	QuestionType QuestionType `json:"question_type" binding:"required"`
}

// Validate проверяет валидность запроса на генерацию вопроса.
func (r *GenerateQuestionRequest) Validate() error {
	if r.Language == "" {
		return errors.New("language is required")
	}

	if !r.Language.IsValid() {
		return fmt.Errorf("unsupported language: %s", r.Language)
	}

	if r.Topic == "" {
		return errors.New("topic is required")
	}

	if !IsValidTopic(r.Language, r.Topic) {
		return fmt.Errorf("invalid topic '%s' for language '%s'", r.Topic, r.Language)
	}

	if !r.Difficulty.IsValid() {
		return fmt.Errorf("invalid difficulty: %s", r.Difficulty)
	}

	if !r.QuestionType.IsValid() {
		return fmt.Errorf("invalid question type: %s", r.QuestionType)
	}

	return nil
}
