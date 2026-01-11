package quiz

import (
	"context"

	"github.com/casnerano/snippet-war/internal/client/content_service"
	models "github.com/casnerano/snippet-war/internal/model/quiz"
)

type contentProvider interface {
	GetQuestions(ctx context.Context, tgUserID string, args content_service.GetQuestionsArgs) ([]*models.Question, error)
}

type Quiz struct {
	contentProvider contentProvider
}

func New(contentProvider contentProvider) *Quiz {
	return &Quiz{
		contentProvider: contentProvider,
	}
}

type GetQuestionsArgs struct {
	Language   models.Language
	Topics     []string
	Difficulty models.Difficulty
	Limit      uint32
}

func (q *Quiz) GetQuestions(ctx context.Context, args GetQuestionsArgs) ([]*models.Question, error) {
	tgUserID := "1"

	return q.contentProvider.GetQuestions(ctx, tgUserID, content_service.GetQuestionsArgs{
		Language:   args.Language,
		Topics:     args.Topics,
		Difficulty: args.Difficulty,
		Limit:      args.Limit,
	})
}
