package content_service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/casnerano/snippet-war/internal/config"
	models "github.com/casnerano/snippet-war/internal/models/quiz"
)

type ContentService struct {
	baseURL    string
	httpClient *http.Client
}

func New(_ context.Context) *ContentService {
	cfg := config.GetContentServiceConfig()

	return &ContentService{
		httpClient: &http.Client{},
		baseURL:    cfg.Host,
	}
}

type GetQuestionsArgs struct {
	Language   models.Language
	Topics     []string
	Difficulty models.Difficulty
	Limit      uint32
}

func (s *ContentService) GetQuestions(ctx context.Context, args GetQuestionsArgs) ([]*models.Question, error) {
	payload := struct {
		Language   models.Language   `json:"language"`
		Topics     []string          `json:"topics"`
		Difficulty models.Difficulty `json:"difficulty"`
		Count      uint32            `json:"count"`
	}{
		Language:   args.Language,
		Topics:     args.Topics,
		Difficulty: args.Difficulty,
		Count:      args.Limit,
	}

	bPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		s.baseURL+"/api/questions/batch",
		bytes.NewReader(bPayload),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := s.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed request: %w", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", response.StatusCode, string(body))
	}

	var questions Questions
	if err = json.NewDecoder(response.Body).Decode(&questions); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return questions.ToModels(), nil
}

type Question struct {
	ID          string            `json:"id"`
	Language    models.Language   `json:"language"`
	Topic       string            `json:"topic"`
	Difficulty  models.Difficulty `json:"difficulty"`
	Code        string            `json:"code"`
	Question    string            `json:"question"`
	Options     []string          `json:"options,omitempty"`
	Answers     []string          `json:"correct_answers"`
	Explanation string            `json:"explanation"`
	Type        models.AnswerType `json:"question_type"`
}

type Questions []*Question

func (q Questions) ToModels() []*models.Question {
	questions := make([]*models.Question, len(q))
	for idx := range q {
		questions[idx] = q[idx].ToModel()
	}
	return questions
}

func (q *Question) ToModel() *models.Question {
	question := models.Question{
		ID:         q.ID,
		Language:   q.Language,
		Topic:      q.Topic,
		Difficulty: q.Difficulty,
		Content: models.Content{
			Text: q.Question,
		},
		Explanation: q.Explanation,
	}

	if q.Code != "" {
		question.Content.Code = &q.Code
	}

	switch q.Type {
	case models.AnswerTypeFreeText:
		question.Answer = &models.FreeTextAnswer{
			CorrectAnswers: q.Answers,
		}
	case models.AnswerTypeMultipleChoice:
		question.Answer = &models.MultipleChoiceAnswer{
			Options:        q.Options,
			CorrectOptions: q.Answers,
		}
	}

	return &question
}
