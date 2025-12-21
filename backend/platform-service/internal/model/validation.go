package model

import (
	"errors"
	"fmt"
)

// ValidateMultipleChoiceOptions проверяет валидность опций для вопроса с множественным выбором.
func ValidateMultipleChoiceOptions(options []string) error {
	if len(options) < 2 {
		return errors.New("multiple choice question must have at least 2 options")
	}
	if len(options) > 5 {
		return errors.New("multiple choice question must have at most 5 options")
	}
	return nil
}

// ValidateMultipleChoiceAnswer проверяет, что correct_answer является валидным текстом из options для multiple_choice.
// Принимает correctAnswer как interface{} и options, возвращает валидированный текст ответа.
func ValidateMultipleChoiceAnswer(correctAnswer interface{}, options []string) (string, error) {
	var answer string
	switch v := correctAnswer.(type) {
	case string:
		answer = v
	case fmt.Stringer:
		answer = v.String()
	default:
		return "", fmt.Errorf("correct answer must be a string for multiple choice, got %T", correctAnswer)
	}

	if answer == "" {
		return "", errors.New("correct answer must be a non-empty string for multiple choice")
	}

	// Проверяем, что ответ есть в списке options
	found := false
	for _, option := range options {
		if option == answer {
			found = true
			break
		}
	}

	if !found {
		return "", fmt.Errorf("correct answer '%s' must be one of the options: %v", answer, options)
	}

	return answer, nil
}

// ValidateMultipleChoiceAnswerString проверяет, что correctAnswer (строка) является валидным текстом из options для multiple_choice.
func ValidateMultipleChoiceAnswerString(correctAnswer string, options []string) error {
	if correctAnswer == "" {
		return errors.New("correct answer must be a non-empty string for multiple choice")
	}

	// Проверяем, что ответ есть в списке options
	for _, option := range options {
		if option == correctAnswer {
			return nil
		}
	}

	return fmt.Errorf("correct answer '%s' must be one of the options: %v", correctAnswer, options)
}

// ValidateFreeTextAnswer проверяет, что correctAnswer является валидной строкой для free_text.
// Принимает correctAnswer как interface{} и возвращает строку.
func ValidateFreeTextAnswer(correctAnswer interface{}) (string, error) {
	var answer string
	switch v := correctAnswer.(type) {
	case string:
		answer = v
	case fmt.Stringer:
		answer = v.String()
	default:
		return "", fmt.Errorf("correct answer must be a string for free text, got %T", correctAnswer)
	}

	if answer == "" {
		return "", errors.New("correct answer must be a non-empty string for free text")
	}

	return answer, nil
}

// ValidateFreeTextAnswerString проверяет, что correctAnswer (строка) является валидной для free_text.
func ValidateFreeTextAnswerString(correctAnswer string) error {
	if correctAnswer == "" {
		return errors.New("correct answer must be a non-empty string for free text")
	}
	return nil
}
