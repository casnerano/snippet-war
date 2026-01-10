package bootstrap

import (
	"context"

	content_service_client "github.com/casnerano/snippet-war/internal/clients/content_service"
	"github.com/casnerano/snippet-war/internal/server"
	quiz_service "github.com/casnerano/snippet-war/internal/services/quiz"
)

type Servers struct {
	Quiz *server.Quiz
}

func InitServers(ctx context.Context) (*Servers, error) {
	contentServiceClient := content_service_client.New(ctx)
	quizService := quiz_service.New(contentServiceClient)

	services := Servers{
		Quiz: server.NewQuiz(quizService),
	}

	return &services, nil
}
