package openrouter

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/casnerano/snippet-war/internal/config"
	openrouterapi "github.com/revrost/go-openrouter"
)

// Client представляет клиент для работы с OpenRouter API.
type Client struct {
	apiClient *openrouterapi.Client
	model     string
	timeout   time.Duration
	maxTokens int
}

// NewClient создает новый клиент OpenRouter на основе конфигурации.
func NewClient(cfg *config.OpenRouterConfig) *Client {
	clientConfig := openrouterapi.DefaultConfig(cfg.APIKey)
	if cfg.BaseURL != "" {
		clientConfig.BaseURL = cfg.BaseURL
	}

	apiClient := openrouterapi.NewClientWithConfig(*clientConfig)

	return &Client{
		apiClient: apiClient,
		model:     cfg.Model,
		timeout:   cfg.Timeout,
		maxTokens: cfg.MaxTokens,
	}
}

// GenerateQuestion генерирует вопрос на основе промпта, используя OpenRouter API.
func (c *Client) GenerateQuestion(ctx context.Context, prompt string) (string, error) {
	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// Формируем запрос
	request := openrouterapi.ChatCompletionRequest{
		Model:     c.model,
		Messages:  []openrouterapi.ChatCompletionMessage{openrouterapi.UserMessage(prompt)},
		MaxTokens: c.maxTokens,
	}

	// Отправляем запрос к API
	response, err := c.apiClient.CreateChatCompletion(ctx, request)
	if err != nil {
		// Обработка различных типов ошибок
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("request timeout after %v: %w", c.timeout, err)
		}
		if ctx.Err() == context.Canceled {
			return "", fmt.Errorf("request canceled: %w", err)
		}

		// Проверка на ошибки API
		var apiErr *openrouterapi.APIError
		if errors.As(err, &apiErr) {
			// Обработка специфичных ошибок API
			switch apiErr.HTTPStatusCode {
			case 401:
				return "", fmt.Errorf("authentication failed: %w", err)
			case 429:
				return "", fmt.Errorf("rate limit exceeded: %w", err)
			case 500, 502, 503, 504:
				return "", fmt.Errorf("server error (status %d): %w", apiErr.HTTPStatusCode, err)
			default:
				return "", fmt.Errorf("API error (status %d): %w", apiErr.HTTPStatusCode, err)
			}
		}

		// Проверка на RequestError
		var reqErr *openrouterapi.RequestError
		if errors.As(err, &reqErr) {
			switch reqErr.HTTPStatusCode {
			case 401:
				return "", fmt.Errorf("authentication failed: %w", err)
			case 429:
				return "", fmt.Errorf("rate limit exceeded: %w", err)
			case 500, 502, 503, 504:
				return "", fmt.Errorf("server error (status %d): %w", reqErr.HTTPStatusCode, err)
			default:
				return "", fmt.Errorf("request error (status %d): %w", reqErr.HTTPStatusCode, err)
			}
		}

		// Ошибки сети и другие
		return "", fmt.Errorf("failed to generate question: %w", err)
	}

	// Проверяем наличие ответов
	if len(response.Choices) == 0 {
		return "", fmt.Errorf("empty response from API: no choices returned")
	}

	// Извлекаем текст из первого ответа
	choice := response.Choices[0]
	content := choice.Message.Content

	// Извлекаем текст из Content
	text := content.Text
	if text == "" && len(content.Multi) > 0 {
		// Если текст пустой, но есть multi-part контент, попробуем извлечь текст из частей
		for _, part := range content.Multi {
			if part.Type == openrouterapi.ChatMessagePartTypeText && part.Text != "" {
				text = part.Text
				break
			}
		}
	}

	if text == "" {
		return "", fmt.Errorf("empty response from API: no text content in response")
	}

	return text, nil
}
