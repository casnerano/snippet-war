package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/casnerano/snippet-war/internal/model"
	"github.com/casnerano/snippet-war/internal/service"
)

// QuestionHandler обрабатывает HTTP запросы для генерации вопросов.
type QuestionHandler struct {
	service *service.QuestionService
	logger  *slog.Logger
}

// NewQuestionHandler создает новый экземпляр QuestionHandler.
func NewQuestionHandler(questionService *service.QuestionService, logger *slog.Logger) *QuestionHandler {
	if logger == nil {
		logger = slog.Default()
	}

	return &QuestionHandler{
		service: questionService,
		logger:  logger,
	}
}

// GenerateQuestion обрабатывает POST запрос на генерацию вопроса.
func (h *QuestionHandler) GenerateQuestion(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		h.respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Проверяем Content-Type
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		h.respondError(w, http.StatusBadRequest, "content-type must be application/json")
		return
	}

	// Парсим JSON запрос
	var req model.GenerateQuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("failed to decode request body", "error", err)
		h.respondError(w, http.StatusBadRequest, "invalid JSON format")
		return
	}

	// Генерируем вопрос через сервис
	question, err := h.service.GenerateQuestion(r.Context(), &req)
	if err != nil {
		h.logger.Error("failed to generate question", "error", err)
		
		// Определяем тип ошибки для правильного HTTP статуса
		// Ошибки валидации возвращаем как 400, остальные как 500
		statusCode := http.StatusInternalServerError
		if isValidationError(err) {
			statusCode = http.StatusBadRequest
		}
		
		h.respondError(w, statusCode, err.Error())
		return
	}

	// Возвращаем успешный ответ
	h.respondJSON(w, http.StatusOK, question)
}

// respondJSON отправляет JSON ответ с указанным статус кодом.
func (h *QuestionHandler) respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}

// respondError отправляет JSON ответ с ошибкой.
func (h *QuestionHandler) respondError(w http.ResponseWriter, statusCode int, message string) {
	h.respondJSON(w, statusCode, map[string]string{
		"error": message,
	})
}

// isValidationError проверяет, является ли ошибка ошибкой валидации.
func isValidationError(err error) bool {
	// Проверяем наличие ключевых слов в сообщении об ошибке
	errorMsg := strings.ToLower(err.Error())
	return strings.Contains(errorMsg, "invalid") ||
		strings.Contains(errorMsg, "validation") ||
		strings.Contains(errorMsg, "required") ||
		strings.Contains(errorMsg, "unsupported")
}

