package config

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		setupEnv    func()
		cleanupEnv  func()
		wantErr     bool
		errContains string
		validate    func(*testing.T, *Config)
	}{
		{
			name: "success with all fields",
			setupEnv: func() {
				os.Setenv("OPENROUTER_API_KEY", "test-api-key")
				os.Setenv("OPENROUTER_MODEL", "test-model")
				os.Setenv("OPENROUTER_BASE_URL", "https://custom.api.com/v1")
				os.Setenv("OPENROUTER_TIMEOUT", "60s")
				os.Setenv("OPENROUTER_MAX_TOKENS", "3000")
			},
			cleanupEnv: func() {
				os.Unsetenv("OPENROUTER_API_KEY")
				os.Unsetenv("OPENROUTER_MODEL")
				os.Unsetenv("OPENROUTER_BASE_URL")
				os.Unsetenv("OPENROUTER_TIMEOUT")
				os.Unsetenv("OPENROUTER_MAX_TOKENS")
			},
			wantErr: false,
			validate: func(t *testing.T, cfg *Config) {
				if cfg.OpenRouter.APIKey != "test-api-key" {
					t.Errorf("APIKey = %v, want test-api-key", cfg.OpenRouter.APIKey)
				}
				if cfg.OpenRouter.Model != "test-model" {
					t.Errorf("Model = %v, want test-model", cfg.OpenRouter.Model)
				}
				if cfg.OpenRouter.BaseURL != "https://custom.api.com/v1" {
					t.Errorf("BaseURL = %v, want https://custom.api.com/v1", cfg.OpenRouter.BaseURL)
				}
				if cfg.OpenRouter.Timeout != 60*time.Second {
					t.Errorf("Timeout = %v, want 60s", cfg.OpenRouter.Timeout)
				}
				if cfg.OpenRouter.MaxTokens != 3000 {
					t.Errorf("MaxTokens = %v, want 3000", cfg.OpenRouter.MaxTokens)
				}
			},
		},
		{
			name: "success with required fields only (defaults applied)",
			setupEnv: func() {
				os.Setenv("OPENROUTER_API_KEY", "test-api-key")
				os.Setenv("OPENROUTER_MODEL", "test-model")
			},
			cleanupEnv: func() {
				os.Unsetenv("OPENROUTER_API_KEY")
				os.Unsetenv("OPENROUTER_MODEL")
			},
			wantErr: false,
			validate: func(t *testing.T, cfg *Config) {
				if cfg.OpenRouter.APIKey != "test-api-key" {
					t.Errorf("APIKey = %v, want test-api-key", cfg.OpenRouter.APIKey)
				}
				if cfg.OpenRouter.Model != "test-model" {
					t.Errorf("Model = %v, want test-model", cfg.OpenRouter.Model)
				}
				// Проверка значений по умолчанию
				if cfg.OpenRouter.BaseURL != "https://openrouter.ai/api/v1" {
					t.Errorf("BaseURL = %v, want https://openrouter.ai/api/v1", cfg.OpenRouter.BaseURL)
				}
				if cfg.OpenRouter.Timeout != 30*time.Second {
					t.Errorf("Timeout = %v, want 30s", cfg.OpenRouter.Timeout)
				}
				if cfg.OpenRouter.MaxTokens != 2000 {
					t.Errorf("MaxTokens = %v, want 2000", cfg.OpenRouter.MaxTokens)
				}
			},
		},
		{
			name: "missing OPENROUTER_API_KEY",
			setupEnv: func() {
				os.Unsetenv("OPENROUTER_API_KEY")
				os.Setenv("OPENROUTER_MODEL", "test-model")
			},
			cleanupEnv: func() {
				os.Unsetenv("OPENROUTER_API_KEY")
				os.Unsetenv("OPENROUTER_MODEL")
			},
			wantErr:     true,
			errContains: "OPENROUTER_API_KEY is required",
		},
		{
			name: "missing OPENROUTER_MODEL",
			setupEnv: func() {
				os.Setenv("OPENROUTER_API_KEY", "test-api-key")
				os.Unsetenv("OPENROUTER_MODEL")
			},
			cleanupEnv: func() {
				os.Unsetenv("OPENROUTER_API_KEY")
				os.Unsetenv("OPENROUTER_MODEL")
			},
			wantErr:     true,
			errContains: "OPENROUTER_MODEL is required",
		},
		{
			name: "empty OPENROUTER_API_KEY",
			setupEnv: func() {
				os.Setenv("OPENROUTER_API_KEY", "")
				os.Setenv("OPENROUTER_MODEL", "test-model")
			},
			cleanupEnv: func() {
				os.Unsetenv("OPENROUTER_API_KEY")
				os.Unsetenv("OPENROUTER_MODEL")
			},
			wantErr:     true,
			errContains: "OPENROUTER_API_KEY is required",
		},
		{
			name: "empty OPENROUTER_MODEL",
			setupEnv: func() {
				os.Setenv("OPENROUTER_API_KEY", "test-api-key")
				os.Setenv("OPENROUTER_MODEL", "")
			},
			cleanupEnv: func() {
				os.Unsetenv("OPENROUTER_API_KEY")
				os.Unsetenv("OPENROUTER_MODEL")
			},
			wantErr:     true,
			errContains: "OPENROUTER_MODEL is required",
		},
		{
			name: "invalid OPENROUTER_TIMEOUT format",
			setupEnv: func() {
				os.Setenv("OPENROUTER_API_KEY", "test-api-key")
				os.Setenv("OPENROUTER_MODEL", "test-model")
				os.Setenv("OPENROUTER_TIMEOUT", "invalid-duration")
			},
			cleanupEnv: func() {
				os.Unsetenv("OPENROUTER_API_KEY")
				os.Unsetenv("OPENROUTER_MODEL")
				os.Unsetenv("OPENROUTER_TIMEOUT")
			},
			wantErr:     true,
			errContains: "invalid OPENROUTER_TIMEOUT format",
		},
		{
			name: "invalid OPENROUTER_MAX_TOKENS format",
			setupEnv: func() {
				os.Setenv("OPENROUTER_API_KEY", "test-api-key")
				os.Setenv("OPENROUTER_MODEL", "test-model")
				os.Setenv("OPENROUTER_MAX_TOKENS", "not-a-number")
			},
			cleanupEnv: func() {
				os.Unsetenv("OPENROUTER_API_KEY")
				os.Unsetenv("OPENROUTER_MODEL")
				os.Unsetenv("OPENROUTER_MAX_TOKENS")
			},
			wantErr:     true,
			errContains: "invalid OPENROUTER_MAX_TOKENS format",
		},
		{
			name: "valid timeout with different units",
			setupEnv: func() {
				os.Setenv("OPENROUTER_API_KEY", "test-api-key")
				os.Setenv("OPENROUTER_MODEL", "test-model")
				os.Setenv("OPENROUTER_TIMEOUT", "2m")
			},
			cleanupEnv: func() {
				os.Unsetenv("OPENROUTER_API_KEY")
				os.Unsetenv("OPENROUTER_MODEL")
				os.Unsetenv("OPENROUTER_TIMEOUT")
			},
			wantErr: false,
			validate: func(t *testing.T, cfg *Config) {
				if cfg.OpenRouter.Timeout != 2*time.Minute {
					t.Errorf("Timeout = %v, want 2m", cfg.OpenRouter.Timeout)
				}
			},
		},
		{
			name: "valid max tokens as string number",
			setupEnv: func() {
				os.Setenv("OPENROUTER_API_KEY", "test-api-key")
				os.Setenv("OPENROUTER_MODEL", "test-model")
				os.Setenv("OPENROUTER_MAX_TOKENS", "5000")
			},
			cleanupEnv: func() {
				os.Unsetenv("OPENROUTER_API_KEY")
				os.Unsetenv("OPENROUTER_MODEL")
				os.Unsetenv("OPENROUTER_MAX_TOKENS")
			},
			wantErr: false,
			validate: func(t *testing.T, cfg *Config) {
				if cfg.OpenRouter.MaxTokens != 5000 {
					t.Errorf("MaxTokens = %v, want 5000", cfg.OpenRouter.MaxTokens)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Сохраняем текущие значения переменных окружения
			originalAPIKey := os.Getenv("OPENROUTER_API_KEY")
			originalModel := os.Getenv("OPENROUTER_MODEL")
			originalBaseURL := os.Getenv("OPENROUTER_BASE_URL")
			originalTimeout := os.Getenv("OPENROUTER_TIMEOUT")
			originalMaxTokens := os.Getenv("OPENROUTER_MAX_TOKENS")

			// Настраиваем окружение для теста
			tt.setupEnv()

			// Восстанавливаем окружение после теста
			defer func() {
				tt.cleanupEnv()
				// Восстанавливаем оригинальные значения, если они были
				if originalAPIKey != "" {
					os.Setenv("OPENROUTER_API_KEY", originalAPIKey)
				}
				if originalModel != "" {
					os.Setenv("OPENROUTER_MODEL", originalModel)
				}
				if originalBaseURL != "" {
					os.Setenv("OPENROUTER_BASE_URL", originalBaseURL)
				}
				if originalTimeout != "" {
					os.Setenv("OPENROUTER_TIMEOUT", originalTimeout)
				}
				if originalMaxTokens != "" {
					os.Setenv("OPENROUTER_MAX_TOKENS", originalMaxTokens)
				}
			}()

			// Выполняем тест
			cfg, err := Load()

			// Проверяем ошибки
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err != nil && tt.errContains != "" {
					if !strings.Contains(err.Error(), tt.errContains) {
						t.Errorf("Load() error = %v, should contain %v", err, tt.errContains)
					}
				}
				return
			}

			// Проверяем результат, если ошибки нет
			if cfg == nil {
				t.Fatal("Load() returned nil config without error")
			}

			if tt.validate != nil {
				tt.validate(t, cfg)
			}
		})
	}
}
