package server

import (
	"context"
	"embed"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed templates/*
var webFS embed.FS

type Server struct {
	httpServer *http.Server
	router     *chi.Mux
}

func New(addr string) *Server {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Handle("/", http.FileServer(http.FS(webFS)))

	server := Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: router,
		},
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
