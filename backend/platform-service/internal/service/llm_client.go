package service

import "context"

// LLMClient представляет интерфейс для работы с различными LLM провайдерами.
// Это позволяет использовать разные провайдеры (OpenRouter, ProxyAPI и т.д.)
// через единый интерфейс.
type LLMClient interface {
	// GenerateQuestion генерирует ответ на основе промпта.
	// Возвращает текстовый ответ от LLM или ошибку.
	GenerateQuestion(ctx context.Context, prompt string) (string, error)
}

