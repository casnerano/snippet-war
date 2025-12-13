package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/casnerano/snippet-war/internal/model"
)

// ParseLLMResponse парсит JSON строку ответа от LLM в структуру LLMQuestionResponse.
// Удаляет markdown code blocks, если LLM обернул ответ в ```json ... ```.
func ParseLLMResponse(jsonStr string) (*model.LLMQuestionResponse, error) {
	if jsonStr == "" {
		return nil, errors.New("empty JSON string")
	}

	// Удаляем markdown code blocks если есть
	cleaned := jsonStr
	cleaned = strings.TrimSpace(cleaned)

	// Удаляем ```json ... ``` или ``` ... ```
	if strings.HasPrefix(cleaned, "```") {
		lines := strings.Split(cleaned, "\n")
		// Удаляем первую строку с ```json или ```
		if len(lines) > 1 {
			lines = lines[1:]
		}
		// Удаляем последнюю строку с ```
		if len(lines) > 0 && strings.HasPrefix(strings.TrimSpace(lines[len(lines)-1]), "```") {
			lines = lines[:len(lines)-1]
		}
		cleaned = strings.Join(lines, "\n")
	}

	cleaned = strings.TrimSpace(cleaned)

	var resp model.LLMQuestionResponse
	if err := json.Unmarshal([]byte(cleaned), &resp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &resp, nil
}

// ValidateLLMResponse валидирует ответ от LLM, проверяя соответствие запросу.
func ValidateLLMResponse(resp *model.LLMQuestionResponse, req *model.GenerateQuestionRequest) error {
	if resp == nil {
		return errors.New("response is nil")
	}

	if req == nil {
		return errors.New("request is nil")
	}

	// Проверяем наличие всех обязательных полей
	if resp.Code == "" {
		return errors.New("code is required")
	}

	if resp.Question == "" {
		return errors.New("question is required")
	}

	if resp.Explanation == "" {
		return errors.New("explanation is required")
	}

	// Проверяем соответствие языка, темы, сложности и типа ответа запросу
	if resp.Language != req.Language {
		return fmt.Errorf("language mismatch: expected %s, got %s", req.Language, resp.Language)
	}

	if resp.Topic != req.Topic {
		return fmt.Errorf("topic mismatch: expected %s, got %s", req.Topic, resp.Topic)
	}

	if resp.Difficulty != req.Difficulty {
		return fmt.Errorf("difficulty mismatch: expected %s, got %s", req.Difficulty, resp.Difficulty)
	}

	if resp.QuestionType != req.QuestionType {
		return fmt.Errorf("question type mismatch: expected %s, got %s", req.QuestionType, resp.QuestionType)
	}

	// Проверяем структуру в зависимости от типа ответа
	switch resp.QuestionType {
	case model.QuestionTypeMultipleChoice:
		// Проверяем наличие options (минимум 2, максимум 5)
		if err := model.ValidateMultipleChoiceOptions(resp.Options); err != nil {
			return fmt.Errorf("invalid options: %w", err)
		}

		// Проверяем валидный correct_answer (текст из options)
		if _, err := model.ValidateMultipleChoiceAnswer(resp.CorrectAnswer, resp.Options); err != nil {
			return fmt.Errorf("invalid correct answer: %w", err)
		}

	case model.QuestionTypeFreeText:
		// Проверяем наличие correct_answer (непустая строка)
		if _, err := model.ValidateFreeTextAnswer(resp.CorrectAnswer); err != nil {
			return fmt.Errorf("invalid correct answer: %w", err)
		}

	default:
		return fmt.Errorf("unknown question type: %s", resp.QuestionType)
	}

	// Используем встроенную валидацию модели для дополнительных проверок
	if err := resp.Validate(); err != nil {
		return fmt.Errorf("response validation failed: %w", err)
	}

	return nil
}
