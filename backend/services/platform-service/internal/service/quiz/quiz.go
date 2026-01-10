package quiz

import (
	"context"

	desc "github.com/casnerano/snippet-war/pkg/api/v1/quiz"
)

type Quiz struct {
	desc.UnimplementedQuizServer
}

func (q *Quiz) ListQuestions(context.Context, *desc.ListQuestions_Request) (*desc.ListQuestions_Response, error) {
	questions := []*desc.Question{
		{
			Id:       "1",
			Language: desc.Language_LANGUAGE_GO,
		},
	}

	return &desc.ListQuestions_Response{Questions: questions}, nil
}
