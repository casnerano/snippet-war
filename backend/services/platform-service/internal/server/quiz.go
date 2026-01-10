package server

import (
	"context"
	"log/slog"

	quiz_models "github.com/casnerano/snippet-war/internal/models/quiz"
	quiz_service "github.com/casnerano/snippet-war/internal/services/quiz"
	desc "github.com/casnerano/snippet-war/pkg/api/v1/quiz"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type quizService interface {
	GetQuestions(ctx context.Context, args quiz_service.GetQuestionsArgs) ([]*quiz_models.Question, error)
}

type Quiz struct {
	desc.UnimplementedQuizServer

	quizService quizService
}

func NewQuiz(quizService quizService) *Quiz {
	return &Quiz{
		quizService: quizService,
	}
}

func (q *Quiz) ListQuestions(ctx context.Context, request *desc.ListQuestions_Request) (*desc.ListQuestions_Response, error) {
	questions, err := q.quizService.GetQuestions(ctx, quiz_service.GetQuestionsArgs{
		Language:   ProtoToLanguage(request.Language),
		Topics:     request.Topics,
		Difficulty: ProtoToDifficulty(request.Difficulty),
		Limit:      request.Limit,
	})
	if err != nil {
		var (
			statusCode = codes.Internal
			logLevel   = slog.LevelError
		)

		slog.Log(ctx, logLevel, "failed get questions", "error", err)

		return nil, status.Error(statusCode, statusCode.String())
	}

	response := desc.ListQuestions_Response{
		Questions: QuestionsToProto(questions),
	}

	return &response, nil
}

func (q *Quiz) Register(server *grpc.Server) {
	desc.RegisterQuizServer(server, q)
}
