package service

import (
	"strings"
	"testing"

	"github.com/casnerano/snippet-war/internal/model"
)

func TestParseLLMResponse(t *testing.T) {
	tests := []struct {
		name      string
		jsonStr   string
		wantErr   bool
		errMsg    string
		checkFunc func(*testing.T, *model.LLMQuestionResponse)
	}{
		{
			name: "parse valid JSON response multiple_choice",
			jsonStr: `{
				"code": "x = 5",
				"question": "What is the value of x?",
				"question_type": "multiple_choice",
				"options": ["3", "4", "5", "6"],
				"correct_answer": 2,
				"explanation": "x is assigned the value 5",
				"difficulty": "beginner",
				"topic": "variables_types",
				"language": "python"
			}`,
			wantErr: false,
			checkFunc: func(t *testing.T, resp *model.LLMQuestionResponse) {
				if resp.Code != "x = 5" {
					t.Errorf("ParseLLMResponse() Code = %v, want 'x = 5'", resp.Code)
				}
				if resp.Question != "What is the value of x?" {
					t.Errorf("ParseLLMResponse() Question = %v, want 'What is the value of x?'", resp.Question)
				}
				if resp.QuestionType != model.QuestionTypeMultipleChoice {
					t.Errorf("ParseLLMResponse() QuestionType = %v, want %v", resp.QuestionType, model.QuestionTypeMultipleChoice)
				}
				if len(resp.Options) != 4 {
					t.Errorf("ParseLLMResponse() Options length = %v, want 4", len(resp.Options))
				}
			},
		},
		{
			name: "parse valid JSON response free_text",
			jsonStr: `{
				"code": "var s []int",
				"question": "What is the type of s?",
				"question_type": "free_text",
				"correct_answer": "slice",
				"acceptable_variants": ["[]int", "slice of int"],
				"case_sensitive": false,
				"explanation": "s is a slice of integers",
				"difficulty": "intermediate",
				"topic": "slices",
				"language": "go"
			}`,
			wantErr: false,
			checkFunc: func(t *testing.T, resp *model.LLMQuestionResponse) {
				if resp.Code != "var s []int" {
					t.Errorf("ParseLLMResponse() Code = %v, want 'var s []int'", resp.Code)
				}
				if resp.QuestionType != model.QuestionTypeFreeText {
					t.Errorf("ParseLLMResponse() QuestionType = %v, want %v", resp.QuestionType, model.QuestionTypeFreeText)
				}
				if resp.CorrectAnswer != "slice" {
					t.Errorf("ParseLLMResponse() CorrectAnswer = %v, want 'slice'", resp.CorrectAnswer)
				}
			},
		},
		{
			name:    "parse JSON with markdown code blocks json",
			jsonStr: "```json\n{\n  \"code\": \"x = 5\",\n  \"question\": \"What is the value of x?\",\n  \"question_type\": \"multiple_choice\",\n  \"options\": [\"3\", \"4\", \"5\"],\n  \"correct_answer\": 2,\n  \"explanation\": \"x is assigned the value 5\",\n  \"difficulty\": \"beginner\",\n  \"topic\": \"variables_types\",\n  \"language\": \"python\"\n}\n```",
			wantErr: false,
			checkFunc: func(t *testing.T, resp *model.LLMQuestionResponse) {
				if resp.Code != "x = 5" {
					t.Errorf("ParseLLMResponse() Code = %v, want 'x = 5'", resp.Code)
				}
			},
		},
		{
			name:    "parse JSON with markdown code blocks without json",
			jsonStr: "```\n{\n  \"code\": \"x = 5\",\n  \"question\": \"What is the value of x?\",\n  \"question_type\": \"multiple_choice\",\n  \"options\": [\"3\", \"4\", \"5\"],\n  \"correct_answer\": 2,\n  \"explanation\": \"x is assigned the value 5\",\n  \"difficulty\": \"beginner\",\n  \"topic\": \"variables_types\",\n  \"language\": \"python\"\n}\n```",
			wantErr: false,
			checkFunc: func(t *testing.T, resp *model.LLMQuestionResponse) {
				if resp.Code != "x = 5" {
					t.Errorf("ParseLLMResponse() Code = %v, want 'x = 5'", resp.Code)
				}
			},
		},
		{
			name:    "handle invalid JSON",
			jsonStr: `{invalid json}`,
			wantErr: true,
			errMsg:  "failed to parse JSON",
		},
		{
			name:    "handle empty response",
			jsonStr: "",
			wantErr: true,
			errMsg:  "empty JSON string",
		},
		{
			name:    "handle response with missing fields",
			jsonStr: `{"code": "x = 5"}`,
			wantErr: false, // парсинг должен пройти, валидация будет в ValidateLLMResponse
			checkFunc: func(t *testing.T, resp *model.LLMQuestionResponse) {
				if resp.Code != "x = 5" {
					t.Errorf("ParseLLMResponse() Code = %v, want 'x = 5'", resp.Code)
				}
				if resp.Question != "" {
					t.Errorf("ParseLLMResponse() Question should be empty, got %v", resp.Question)
				}
			},
		},
		{
			name:    "handle whitespace only",
			jsonStr: "   \n\t  ",
			wantErr: true,
			errMsg:  "failed to parse JSON",
		},
		{
			name: "parse JSON with extra whitespace",
			jsonStr: `
			{
				"code": "x = 5",
				"question": "What is the value of x?",
				"question_type": "multiple_choice",
				"options": ["3", "4", "5"],
				"correct_answer": 2,
				"explanation": "x is assigned the value 5",
				"difficulty": "beginner",
				"topic": "variables_types",
				"language": "python"
			}
			`,
			wantErr: false,
			checkFunc: func(t *testing.T, resp *model.LLMQuestionResponse) {
				if resp.Code != "x = 5" {
					t.Errorf("ParseLLMResponse() Code = %v, want 'x = 5'", resp.Code)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := ParseLLMResponse(tt.jsonStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLLMResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseLLMResponse() expected error, got nil")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ParseLLMResponse() error message = %v, want containing '%s'", err.Error(), tt.errMsg)
				}
			} else {
				if resp == nil {
					t.Fatal("ParseLLMResponse() returned nil response")
				}
				if tt.checkFunc != nil {
					tt.checkFunc(t, resp)
				}
			}
		})
	}
}

func TestValidateLLMResponse(t *testing.T) {
	validRequest := &model.GenerateQuestionRequest{
		Language:     model.LanguagePython,
		Topic:        model.TopicID("variables_types"),
		Difficulty:   model.DifficultyBeginner,
		QuestionType: model.QuestionTypeMultipleChoice,
	}

	validFreeTextRequest := &model.GenerateQuestionRequest{
		Language:     model.LanguageGo,
		Topic:        model.TopicID("slices"),
		Difficulty:   model.DifficultyIntermediate,
		QuestionType: model.QuestionTypeFreeText,
	}

	tests := []struct {
		name    string
		resp    *model.LLMQuestionResponse
		req     *model.GenerateQuestionRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid multiple_choice response",
			resp: &model.LLMQuestionResponse{
				Code:          "x = 5",
				Question:      "What is the value of x?",
				QuestionType:  model.QuestionTypeMultipleChoice,
				Options:       []string{"3", "4", "5", "6"},
				CorrectAnswer: 2,
				Explanation:   "x is assigned the value 5",
				Difficulty:    model.DifficultyBeginner,
				Topic:         model.TopicID("variables_types"),
				Language:      model.LanguagePython,
			},
			req:     validRequest,
			wantErr: false,
		},
		{
			name: "valid free_text response",
			resp: &model.LLMQuestionResponse{
				Code:          "var s []int",
				Question:      "What is the type of s?",
				QuestionType:  model.QuestionTypeFreeText,
				CorrectAnswer: "slice",
				Explanation:   "s is a slice of integers",
				Difficulty:    model.DifficultyIntermediate,
				Topic:         model.TopicID("slices"),
				Language:      model.LanguageGo,
			},
			req:     validFreeTextRequest,
			wantErr: false,
		},
		{
			name:    "nil response",
			resp:    nil,
			req:     validRequest,
			wantErr: true,
			errMsg:  "response is nil",
		},
		{
			name: "nil request",
			resp: &model.LLMQuestionResponse{
				Code:          "x = 5",
				Question:      "What is the value of x?",
				QuestionType:  model.QuestionTypeMultipleChoice,
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: 2,
				Explanation:   "explanation",
				Difficulty:    model.DifficultyBeginner,
				Topic:         model.TopicID("variables_types"),
				Language:      model.LanguagePython,
			},
			req:     nil,
			wantErr: true,
			errMsg:  "request is nil",
		},
		{
			name: "missing code",
			resp: &model.LLMQuestionResponse{
				Question:      "What is the value of x?",
				QuestionType:  model.QuestionTypeMultipleChoice,
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: 2,
				Explanation:   "explanation",
				Difficulty:    model.DifficultyBeginner,
				Topic:         model.TopicID("variables_types"),
				Language:      model.LanguagePython,
			},
			req:     validRequest,
			wantErr: true,
			errMsg:  "code is required",
		},
		{
			name: "missing question",
			resp: &model.LLMQuestionResponse{
				Code:          "x = 5",
				QuestionType:  model.QuestionTypeMultipleChoice,
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: 2,
				Explanation:   "explanation",
				Difficulty:    model.DifficultyBeginner,
				Topic:         model.TopicID("variables_types"),
				Language:      model.LanguagePython,
			},
			req:     validRequest,
			wantErr: true,
			errMsg:  "question is required",
		},
		{
			name: "missing explanation",
			resp: &model.LLMQuestionResponse{
				Code:          "x = 5",
				Question:      "What is the value of x?",
				QuestionType:  model.QuestionTypeMultipleChoice,
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: 2,
				Difficulty:    model.DifficultyBeginner,
				Topic:         model.TopicID("variables_types"),
				Language:      model.LanguagePython,
			},
			req:     validRequest,
			wantErr: true,
			errMsg:  "explanation is required",
		},
		{
			name: "language mismatch",
			resp: &model.LLMQuestionResponse{
				Code:          "x = 5",
				Question:      "What is the value of x?",
				QuestionType:  model.QuestionTypeMultipleChoice,
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: 2,
				Explanation:   "explanation",
				Difficulty:    model.DifficultyBeginner,
				Topic:         model.TopicID("variables_types"),
				Language:      model.LanguageGo, // не совпадает с запросом
			},
			req:     validRequest,
			wantErr: true,
			errMsg:  "language mismatch",
		},
		{
			name: "topic mismatch",
			resp: &model.LLMQuestionResponse{
				Code:          "x = 5",
				Question:      "What is the value of x?",
				QuestionType:  model.QuestionTypeMultipleChoice,
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: 2,
				Explanation:   "explanation",
				Difficulty:    model.DifficultyBeginner,
				Topic:         model.TopicID("slices"), // не совпадает с запросом
				Language:      model.LanguagePython,
			},
			req:     validRequest,
			wantErr: true,
			errMsg:  "topic mismatch",
		},
		{
			name: "difficulty mismatch",
			resp: &model.LLMQuestionResponse{
				Code:          "x = 5",
				Question:      "What is the value of x?",
				QuestionType:  model.QuestionTypeMultipleChoice,
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: 2,
				Explanation:   "explanation",
				Difficulty:    model.DifficultyAdvanced, // не совпадает с запросом
				Topic:         model.TopicID("variables_types"),
				Language:      model.LanguagePython,
			},
			req:     validRequest,
			wantErr: true,
			errMsg:  "difficulty mismatch",
		},
		{
			name: "question type mismatch",
			resp: &model.LLMQuestionResponse{
				Code:          "x = 5",
				Question:      "What is the value of x?",
				QuestionType:  model.QuestionTypeFreeText, // не совпадает с запросом
				CorrectAnswer: "5",
				Explanation:   "explanation",
				Difficulty:    model.DifficultyBeginner,
				Topic:         model.TopicID("variables_types"),
				Language:      model.LanguagePython,
			},
			req:     validRequest,
			wantErr: true,
			errMsg:  "question type mismatch",
		},
		{
			name: "multiple_choice with too few options",
			resp: &model.LLMQuestionResponse{
				Code:          "x = 5",
				Question:      "What is the value of x?",
				QuestionType:  model.QuestionTypeMultipleChoice,
				Options:       []string{"3"}, // только 1 опция, минимум 2
				CorrectAnswer: 0,
				Explanation:   "explanation",
				Difficulty:    model.DifficultyBeginner,
				Topic:         model.TopicID("variables_types"),
				Language:      model.LanguagePython,
			},
			req:     validRequest,
			wantErr: true,
			errMsg:  "invalid options",
		},
		{
			name: "multiple_choice with too many options",
			resp: &model.LLMQuestionResponse{
				Code:          "x = 5",
				Question:      "What is the value of x?",
				QuestionType:  model.QuestionTypeMultipleChoice,
				Options:       []string{"1", "2", "3", "4", "5", "6"}, // 6 опций, максимум 5
				CorrectAnswer: 0,
				Explanation:   "explanation",
				Difficulty:    model.DifficultyBeginner,
				Topic:         model.TopicID("variables_types"),
				Language:      model.LanguagePython,
			},
			req:     validRequest,
			wantErr: true,
			errMsg:  "invalid options",
		},
		{
			name: "multiple_choice with invalid answer index",
			resp: &model.LLMQuestionResponse{
				Code:          "x = 5",
				Question:      "What is the value of x?",
				QuestionType:  model.QuestionTypeMultipleChoice,
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: 5, // индекс вне диапазона [0, 3)
				Explanation:   "explanation",
				Difficulty:    model.DifficultyBeginner,
				Topic:         model.TopicID("variables_types"),
				Language:      model.LanguagePython,
			},
			req:     validRequest,
			wantErr: true,
			errMsg:  "invalid correct answer",
		},
		{
			name: "multiple_choice with negative answer index",
			resp: &model.LLMQuestionResponse{
				Code:          "x = 5",
				Question:      "What is the value of x?",
				QuestionType:  model.QuestionTypeMultipleChoice,
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: -1, // отрицательный индекс
				Explanation:   "explanation",
				Difficulty:    model.DifficultyBeginner,
				Topic:         model.TopicID("variables_types"),
				Language:      model.LanguagePython,
			},
			req:     validRequest,
			wantErr: true,
			errMsg:  "invalid correct answer",
		},
		{
			name: "free_text with empty answer",
			resp: &model.LLMQuestionResponse{
				Code:          "var s []int",
				Question:      "What is the type of s?",
				QuestionType:  model.QuestionTypeFreeText,
				CorrectAnswer: "", // пустой ответ
				Explanation:   "explanation",
				Difficulty:    model.DifficultyIntermediate,
				Topic:         model.TopicID("slices"),
				Language:      model.LanguageGo,
			},
			req:     validFreeTextRequest,
			wantErr: true,
			errMsg:  "invalid correct answer",
		},
		{
			name: "free_text with non-string answer",
			resp: &model.LLMQuestionResponse{
				Code:          "var s []int",
				Question:      "What is the type of s?",
				QuestionType:  model.QuestionTypeFreeText,
				CorrectAnswer: 123, // не строка
				Explanation:   "explanation",
				Difficulty:    model.DifficultyIntermediate,
				Topic:         model.TopicID("slices"),
				Language:      model.LanguageGo,
			},
			req:     validFreeTextRequest,
			wantErr: true,
			errMsg:  "invalid correct answer",
		},
		{
			name: "multiple_choice with valid minimum options (2)",
			resp: &model.LLMQuestionResponse{
				Code:          "x = 5",
				Question:      "What is the value of x?",
				QuestionType:  model.QuestionTypeMultipleChoice,
				Options:       []string{"3", "5"}, // минимум 2 опции
				CorrectAnswer: 1,
				Explanation:   "explanation",
				Difficulty:    model.DifficultyBeginner,
				Topic:         model.TopicID("variables_types"),
				Language:      model.LanguagePython,
			},
			req:     validRequest,
			wantErr: false,
		},
		{
			name: "multiple_choice with valid maximum options (5)",
			resp: &model.LLMQuestionResponse{
				Code:          "x = 5",
				Question:      "What is the value of x?",
				QuestionType:  model.QuestionTypeMultipleChoice,
				Options:       []string{"1", "2", "3", "4", "5"}, // максимум 5 опций
				CorrectAnswer: 4,
				Explanation:   "explanation",
				Difficulty:    model.DifficultyBeginner,
				Topic:         model.TopicID("variables_types"),
				Language:      model.LanguagePython,
			},
			req:     validRequest,
			wantErr: false,
		},
		{
			name: "multiple_choice with valid answer index at boundary",
			resp: &model.LLMQuestionResponse{
				Code:          "x = 5",
				Question:      "What is the value of x?",
				QuestionType:  model.QuestionTypeMultipleChoice,
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: 2, // последний валидный индекс (len - 1)
				Explanation:   "explanation",
				Difficulty:    model.DifficultyBeginner,
				Topic:         model.TopicID("variables_types"),
				Language:      model.LanguagePython,
			},
			req:     validRequest,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLLMResponse(tt.resp, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateLLMResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateLLMResponse() expected error, got nil")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateLLMResponse() error message = %v, want containing '%s'", err.Error(), tt.errMsg)
				}
			}
		})
	}
}
