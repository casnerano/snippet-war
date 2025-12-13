package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/casnerano/snippet-war/internal/model"
	"github.com/casnerano/snippet-war/internal/openrouter"
)

// QuestionService представляет сервис для генерации вопросов.
type QuestionService struct {
	llmClient LLMClient
	logger    *slog.Logger
}

// NewQuestionService создает новый экземпляр QuestionService.
func NewQuestionService(client LLMClient, logger *slog.Logger) *QuestionService {
	if logger == nil {
		logger = slog.Default()
	}

	return &QuestionService{
		llmClient: client,
		logger:    logger,
	}
}

// GenerateQuestion генерирует вопрос на основе запроса с использованием OpenRouter API.
func (s *QuestionService) GenerateQuestion(ctx context.Context, req *model.GenerateQuestionRequest) (*model.Question, error) {
	// 1. Валидировать запрос
	if err := req.Validate(); err != nil {
		s.logger.Error("request validation failed", "error", err)
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	s.logger.Info("generating question",
		"language", req.Language,
		"topic", req.Topic,
		"difficulty", req.Difficulty,
		"question_type", req.QuestionType,
	)

	// 2. Построить промпт
	prompt := openrouter.BuildPrompt(req)

	// 3. Отправить запрос к LLM API
	responseText, err := s.llmClient.GenerateQuestion(ctx, prompt)
	if err != nil {
		s.logger.Error("failed to generate question from API", "error", err)
		return nil, fmt.Errorf("failed to generate question from API: %w", err)
	}

	s.logger.Debug("received response from API", "response_length", len(responseText))

	// 4. Распарсить JSON ответ от LLM в LLMQuestionResponse
	llmResponse, err := ParseLLMResponse(responseText)
	if err != nil {
		s.logger.Error("failed to parse LLM response", "error", err, "response", responseText)
		return nil, fmt.Errorf("failed to parse LLM response: %w", err)
	}

	// 5. Валидировать полученный ответ от LLM
	if err := ValidateLLMResponse(llmResponse, req); err != nil {
		s.logger.Error("LLM response validation failed", "error", err)
		return nil, fmt.Errorf("LLM response validation failed: %w", err)
	}

	// 6. Преобразовать LLMQuestionResponse в Question (метод ToQuestion)
	question, err := llmResponse.ToQuestion()
	if err != nil {
		s.logger.Error("failed to convert LLM response to question", "error", err)
		return nil, fmt.Errorf("failed to convert LLM response to question: %w", err)
	}

	// 7. ID и CreatedAt уже установлены в ToQuestion через GenerateID()
	// 8. Финальная валидация также выполняется в ToQuestion

	s.logger.Info("question generated successfully", "question_id", question.ID)

	return question, nil
}
