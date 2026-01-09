package quiz

import (
	"context"
	"time"

	desc "github.com/casnerano/snippet-war/pkg/api/v1/quiz"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Quiz struct {
	desc.UnimplementedQuizServer
}

func (q *Quiz) GetQuestions(context.Context, *desc.GetQuestions_Request) (*desc.GetQuestions_Response, error) {
	questions := []*desc.Question{
		{
			Id:        "1",
			Language:  desc.Language_LANGUAGE_GO,
			Topics:    []string{"example"},
			CreatedAt: timestamppb.New(time.Now()),
		},
	}

	return &desc.GetQuestions_Response{Questions: questions}, nil
}
