package proxyapi

import (
	"github.com/casnerano/snippet-war/internal/model"
	"github.com/casnerano/snippet-war/internal/openrouter"
)

// BuildPrompt создает промпт для генерации вопроса на основе параметров запроса.
// Переиспользует функцию из пакета openrouter, так как логика промпта одинакова.
func BuildPrompt(req *model.GenerateQuestionRequest) string {
	// Используем функцию из openrouter, так как она уже реализована
	// и логика промпта не зависит от провайдера
	return openrouter.BuildPrompt(req)
}

