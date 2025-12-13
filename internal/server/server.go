package server

import (
	"context"
	"embed"
	"io/fs"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/casnerano/snippet-war/internal/service"
)

//go:embed templates/*
var webFS embed.FS

type Server struct {
	httpServer      *http.Server
	router          *chi.Mux
	questionService *service.QuestionService
}

func New(addr string, questionService *service.QuestionService) *Server {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Создаем handler для вопросов
	questionHandler := NewQuestionHandler(questionService, slog.Default())

	// Регистрируем API routes
	router.Route("/api", func(r chi.Router) {
		r.Post("/questions/generate", questionHandler.GenerateQuestion)
	})

	subFS, _ := fs.Sub(webFS, "templates")

	router.Handle("/*", http.FileServer(http.FS(subFS)))

	server := Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: router,
		},
		questionService: questionService,
	}

	return &server
}

func (s *Server) Run(ctx context.Context) error {
	slog.Info("Starting server on " + s.httpServer.Addr)

	go func() {
		<-ctx.Done()
		slog.Info("Shutting down server...")
		_ = s.httpServer.Shutdown(context.Background())
	}()

	return s.httpServer.ListenAndServe()
}
