package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/casnerano/snippet-war/internal/config"
	"github.com/casnerano/snippet-war/internal/model"
	"github.com/casnerano/snippet-war/internal/proxyapi"
	"github.com/casnerano/snippet-war/internal/service"
)

func main() {
	// Настройка логгера
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	// load from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}
	logger.Info("Starting question generation test")

	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Проверяем наличие конфигурации ProxyAPI
	if cfg.ProxyAPI.APIKey == "" {
		log.Fatalf("ProxyAPI configuration is required. Please set PROXYAPI_API_KEY and PROXYAPI_MODEL environment variables")
	}

	logger.Info("Config loaded successfully",
		"provider", "ProxyAPI",
		"model", cfg.ProxyAPI.Model,
		"base_url", cfg.ProxyAPI.BaseURL,
		"timeout", cfg.ProxyAPI.Timeout,
		"max_tokens", cfg.ProxyAPI.MaxTokens,
	)

	// Создание клиента ProxyAPI
	proxyAPIClient := proxyapi.NewClient(&cfg.ProxyAPI)

	// Создание сервиса генерации вопросов
	questionService := service.NewQuestionService(proxyAPIClient, logger)

	// Создание контекста с таймаутом
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	ctx, timeoutCancel := context.WithTimeout(ctx, 60*time.Second)
	defer timeoutCancel()

	// Создание запроса на генерацию вопроса
	// Пример: Python, начальный уровень, тема "Переменные и типы данных", множественный выбор
	req := &model.GenerateQuestionRequest{
		Language:     model.LanguagePython,
		Topic:        model.TopicID("variables_types"),
		Difficulty:   model.DifficultyBeginner,
		QuestionType: model.QuestionTypeMultipleChoice,
	}

	logger.Info("Generating question with parameters",
		"language", req.Language,
		"topic", req.Topic,
		"difficulty", req.Difficulty,
		"question_type", req.QuestionType,
	)

	// Генерация вопроса
	question, err := questionService.GenerateQuestion(ctx, req)
	if err != nil {
		log.Fatalf("Failed to generate question: %v", err)
	}

	// Вывод результата
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("GENERATED QUESTION")
	fmt.Println(strings.Repeat("=", 80))

	// Красивый вывод вопроса
	printQuestion(question)

	// JSON вывод для копирования
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("JSON OUTPUT")
	fmt.Println(strings.Repeat("=", 80))
	jsonData, err := json.MarshalIndent(question, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal question to JSON: %v", err)
	}
	fmt.Println(string(jsonData))

	logger.Info("Question generated successfully", "question_id", question.ID)
}

func printQuestion(q *model.Question) {
	fmt.Printf("\nID: %s\n", q.ID)
	fmt.Printf("Language: %s\n", q.Language)
	fmt.Printf("Topic: %s\n", q.Topic)
	fmt.Printf("Difficulty: %s\n", q.Difficulty)
	fmt.Printf("Question Type: %s\n", q.QuestionType)
	fmt.Printf("Created At: %s\n", q.CreatedAt.Format(time.RFC3339))

	fmt.Println("\n--- Code ---")
	fmt.Println(q.Code)

	fmt.Println("\n--- Question ---")
	fmt.Println(q.QuestionText)

	if q.QuestionType == model.QuestionTypeMultipleChoice && len(q.Options) > 0 {
		fmt.Println("\n--- Options ---")
		for i, option := range q.Options {
			marker := " "
			if q.CorrectAnswer == fmt.Sprintf("%d", i) {
				marker = "✓"
			}
			fmt.Printf("%s [%d] %s\n", marker, i, option)
		}
		fmt.Printf("\nCorrect Answer Index: %s\n", q.CorrectAnswer)
	} else if q.QuestionType == model.QuestionTypeFreeText {
		fmt.Printf("\n--- Correct Answer ---\n%s\n", q.CorrectAnswer)
		if len(q.AcceptableVariants) > 0 {
			fmt.Println("\n--- Acceptable Variants ---")
			for _, variant := range q.AcceptableVariants {
				fmt.Printf("  - %s\n", variant)
			}
		}
		fmt.Printf("Case Sensitive: %v\n", q.CaseSensitive)
	}

	fmt.Println("\n--- Explanation ---")
	fmt.Println(q.Explanation)
}
