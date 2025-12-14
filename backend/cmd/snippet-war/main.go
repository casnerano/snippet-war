package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/casnerano/snippet-war/internal/config"
	"github.com/casnerano/snippet-war/internal/proxyapi"
	"github.com/casnerano/snippet-war/internal/server"
	"github.com/casnerano/snippet-war/internal/service"
)

const defaultAddr = ":8081"

var addr string

func main() {
	// load from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Настройка логгера
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	flag.StringVar(&addr, "addr", defaultAddr, "Server address")
	flag.Parse()

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

	// Создание сервера с сервисом
	srv := server.New(addr, questionService)

	if err := srv.Run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
