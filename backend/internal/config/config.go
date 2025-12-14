package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	OpenRouterModel     = "deepseek/deepseek-chat"
	OpenRouterBaseURL   = "https://openrouter.ai/api/v1"
	OpenRouterTimeout   = "30s"
	OpenRouterMaxTokens = 2000

	ProxyAPIBaseURL   = "https://api.proxyapi.ru/openai/v1"
	ProxyAPITimeout   = "30s"
	ProxyAPIMaxTokens = 2000
	ProxyAPIModel     = "gpt-4.1-mini"
)

// OpenRouterConfig содержит конфигурацию для OpenRouter API.
type OpenRouterConfig struct {
	APIKey    string
	Model     string
	BaseURL   string
	Timeout   time.Duration
	MaxTokens int
}

// ProxyAPIConfig содержит конфигурацию для ProxyAPI (OpenAI-совместимый API).
type ProxyAPIConfig struct {
	APIKey    string
	Model     string
	BaseURL   string
	Timeout   time.Duration
	MaxTokens int
}

// Config содержит всю конфигурацию приложения.
type Config struct {
	OpenRouter OpenRouterConfig
	ProxyAPI   ProxyAPIConfig
}

// Load загружает конфигурацию из переменных окружения.
// Возвращает ошибку, если отсутствуют обязательные поля (API_KEY, MODEL).
func Load() (*Config, error) {
	cfg := &Config{}

	// Загрузка конфигурации OpenRouter
	openRouterCfg, err := loadOpenRouterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load OpenRouter config: %w", err)
	}
	cfg.OpenRouter = *openRouterCfg

	// Загрузка конфигурации ProxyAPI (опционально)
	proxyAPICfg, err := loadProxyAPIConfig()
	if err != nil {
		// ProxyAPI не обязателен, поэтому просто логируем ошибку, но не прерываем загрузку
		// В реальном приложении можно использовать логгер
		_ = err
	} else {
		cfg.ProxyAPI = *proxyAPICfg
	}

	return cfg, nil
}

// loadOpenRouterConfig загружает конфигурацию OpenRouter из переменных окружения.
func loadOpenRouterConfig() (*OpenRouterConfig, error) {
	cfg := &OpenRouterConfig{}

	// Обязательные поля
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return nil, errors.New("OPENROUTER_API_KEY is required")
	}
	cfg.APIKey = apiKey

	model := os.Getenv("OPENROUTER_MODEL")
	if model == "" {
		return nil, errors.New("OPENROUTER_MODEL is required")
	}
	cfg.Model = model

	// Опциональные поля с значениями по умолчанию
	baseURL := os.Getenv("OPENROUTER_BASE_URL")
	if baseURL == "" {
		baseURL = "https://openrouter.ai/api/v1"
	}
	cfg.BaseURL = baseURL

	timeoutStr := os.Getenv("OPENROUTER_TIMEOUT")
	if timeoutStr == "" {
		timeoutStr = "30s"
	}
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return nil, fmt.Errorf("invalid OPENROUTER_TIMEOUT format: %w", err)
	}
	cfg.Timeout = timeout

	maxTokensStr := os.Getenv("OPENROUTER_MAX_TOKENS")
	if maxTokensStr == "" {
		maxTokensStr = "2000"
	}
	maxTokens, err := strconv.Atoi(maxTokensStr)
	if err != nil {
		return nil, fmt.Errorf("invalid OPENROUTER_MAX_TOKENS format: %w", err)
	}
	cfg.MaxTokens = maxTokens

	return cfg, nil
}

// loadProxyAPIConfig загружает конфигурацию ProxyAPI из переменных окружения.
func loadProxyAPIConfig() (*ProxyAPIConfig, error) {
	cfg := &ProxyAPIConfig{}

	// Обязательные поля
	apiKey := os.Getenv("PROXYAPI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("PROXYAPI_API_KEY is required")
	}
	cfg.APIKey = apiKey

	model := os.Getenv("PROXYAPI_MODEL")
	if model == "" {
		return nil, errors.New("PROXYAPI_MODEL is required")
	}
	cfg.Model = model

	// Опциональные поля с значениями по умолчанию
	baseURL := os.Getenv("PROXYAPI_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.proxyapi.ru/openai/v1"
	}
	cfg.BaseURL = baseURL

	timeoutStr := os.Getenv("PROXYAPI_TIMEOUT")
	if timeoutStr == "" {
		timeoutStr = "30s"
	}
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return nil, fmt.Errorf("invalid PROXYAPI_TIMEOUT format: %w", err)
	}
	cfg.Timeout = timeout

	maxTokensStr := os.Getenv("PROXYAPI_MAX_TOKENS")
	if maxTokensStr == "" {
		maxTokensStr = "2000"
	}
	maxTokens, err := strconv.Atoi(maxTokensStr)
	if err != nil {
		return nil, fmt.Errorf("invalid PROXYAPI_MAX_TOKENS format: %w", err)
	}
	cfg.MaxTokens = maxTokens

	return cfg, nil
}
