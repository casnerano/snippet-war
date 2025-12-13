package model

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestLLMQuestionResponse_Validate(t *testing.T) {
	tests := []struct {
		name    string
		response LLMQuestionResponse
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid multiple choice response",
			response: LLMQuestionResponse{
				Code:          "x = 5",
				Question:      "What is the value of x?",
				QuestionType:  QuestionTypeMultipleChoice,
				Options:       []string{"3", "4", "5", "6"},
				CorrectAnswer: 2,
				Explanation:   "x is assigned the value 5",
				Difficulty:    DifficultyBeginner,
				Topic:         TopicID("variables_types"),
				Language:      LanguagePython,
			},
			wantErr: false,
		},
		{
			name: "valid free text response",
			response: LLMQuestionResponse{
				Code:               "var s []int",
				Question:           "What is the type of s?",
				QuestionType:       QuestionTypeFreeText,
				CorrectAnswer:      "slice",
				AcceptableVariants: []string{"[]int", "slice of int"},
				CaseSensitive:      false,
				Explanation:        "s is a slice of integers",
				Difficulty:         DifficultyIntermediate,
				Topic:              TopicID("slices"),
				Language:           LanguageGo,
			},
			wantErr: false,
		},
		{
			name: "missing code",
			response: LLMQuestionResponse{
				Question:      "What is the value of x?",
				QuestionType:  QuestionTypeMultipleChoice,
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: 2,
				Explanation:   "x is assigned the value 5",
				Difficulty:    DifficultyBeginner,
				Topic:         TopicID("variables_types"),
				Language:      LanguagePython,
			},
			wantErr: true,
			errMsg:  "code is required",
		},
		{
			name: "missing question",
			response: LLMQuestionResponse{
				Code:          "x = 5",
				QuestionType:  QuestionTypeMultipleChoice,
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: 2,
				Explanation:   "x is assigned the value 5",
				Difficulty:    DifficultyBeginner,
				Topic:         TopicID("variables_types"),
				Language:      LanguagePython,
			},
			wantErr: true,
			errMsg:  "question is required",
		},
		{
			name: "invalid question type",
			response: LLMQuestionResponse{
				Code:          "x = 5",
				Question:      "What is the value of x?",
				QuestionType:  QuestionType("invalid"),
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: 2,
				Explanation:   "x is assigned the value 5",
				Difficulty:    DifficultyBeginner,
				Topic:         TopicID("variables_types"),
				Language:      LanguagePython,
			},
			wantErr: true,
			errMsg:  "invalid question type",
		},
		{
			name: "invalid difficulty",
			response: LLMQuestionResponse{
				Code:          "x = 5",
				Question:      "What is the value of x?",
				QuestionType:  QuestionTypeMultipleChoice,
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: 2,
				Explanation:   "x is assigned the value 5",
				Difficulty:    Difficulty("invalid"),
				Topic:         TopicID("variables_types"),
				Language:      LanguagePython,
			},
			wantErr: true,
			errMsg:  "invalid difficulty",
		},
		{
			name: "invalid language",
			response: LLMQuestionResponse{
				Code:          "x = 5",
				Question:      "What is the value of x?",
				QuestionType:  QuestionTypeMultipleChoice,
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: 2,
				Explanation:   "x is assigned the value 5",
				Difficulty:    DifficultyBeginner,
				Topic:         TopicID("variables_types"),
				Language:      Language("invalid"),
			},
			wantErr: true,
			errMsg:  "invalid language",
		},
		{
			name: "missing topic",
			response: LLMQuestionResponse{
				Code:          "x = 5",
				Question:      "What is the value of x?",
				QuestionType:  QuestionTypeMultipleChoice,
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: 2,
				Explanation:   "x is assigned the value 5",
				Difficulty:    DifficultyBeginner,
				Language:      LanguagePython,
			},
			wantErr: true,
			errMsg:  "topic is required",
		},
		{
			name: "invalid topic for language",
			response: LLMQuestionResponse{
				Code:          "x = 5",
				Question:      "What is the value of x?",
				QuestionType:  QuestionTypeMultipleChoice,
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: 2,
				Explanation:   "x is assigned the value 5",
				Difficulty:    DifficultyBeginner,
				Topic:         TopicID("ownership"), // тема для Rust, не для Python
				Language:      LanguagePython,
			},
			wantErr: true,
			errMsg:  "invalid topic",
		},
		{
			name: "missing explanation",
			response: LLMQuestionResponse{
				Code:          "x = 5",
				Question:      "What is the value of x?",
				QuestionType:  QuestionTypeMultipleChoice,
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: 2,
				Difficulty:    DifficultyBeginner,
				Topic:         TopicID("variables_types"),
				Language:      LanguagePython,
			},
			wantErr: true,
			errMsg:  "explanation is required",
		},
		{
			name: "missing correct answer",
			response: LLMQuestionResponse{
				Code:         "x = 5",
				Question:     "What is the value of x?",
				QuestionType: QuestionTypeMultipleChoice,
				Options:      []string{"3", "4", "5"},
				Explanation:  "x is assigned the value 5",
				Difficulty:   DifficultyBeginner,
				Topic:        TopicID("variables_types"),
				Language:     LanguagePython,
			},
			wantErr: true,
			errMsg:  "correct answer is required",
		},
		{
			name: "multiple choice with too few options",
			response: LLMQuestionResponse{
				Code:          "x = 5",
				Question:      "What is the value of x?",
				QuestionType:  QuestionTypeMultipleChoice,
				Options:       []string{"3"},
				CorrectAnswer: 0,
				Explanation:   "x is assigned the value 5",
				Difficulty:    DifficultyBeginner,
				Topic:         TopicID("variables_types"),
				Language:      LanguagePython,
			},
			wantErr: true,
			errMsg:  "multiple choice question must have at least 2 options",
		},
		{
			name: "multiple choice with invalid answer index",
			response: LLMQuestionResponse{
				Code:          "x = 5",
				Question:      "What is the value of x?",
				QuestionType:  QuestionTypeMultipleChoice,
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: 5, // индекс вне диапазона
				Explanation:   "x is assigned the value 5",
				Difficulty:    DifficultyBeginner,
				Topic:         TopicID("variables_types"),
				Language:      LanguagePython,
			},
			wantErr: true,
			errMsg:  "correct answer index",
		},
		{
			name: "free text with empty answer",
			response: LLMQuestionResponse{
				Code:         "var s []int",
				Question:     "What is the type of s?",
				QuestionType: QuestionTypeFreeText,
				CorrectAnswer: "",
				Explanation:  "s is a slice of integers",
				Difficulty:   DifficultyIntermediate,
				Topic:        TopicID("slices"),
				Language:     LanguageGo,
			},
			wantErr: true,
			errMsg:  "correct answer must be a non-empty string",
		},
		{
			name: "free text with non-string answer",
			response: LLMQuestionResponse{
				Code:         "var s []int",
				Question:     "What is the type of s?",
				QuestionType: QuestionTypeFreeText,
				CorrectAnswer: 123, // не строка
				Explanation:  "s is a slice of integers",
				Difficulty:   DifficultyIntermediate,
				Topic:        TopicID("slices"),
				Language:     LanguageGo,
			},
			wantErr: true,
			errMsg:  "correct answer must be a string",
		},
		{
			name: "unknown question type",
			response: LLMQuestionResponse{
				Code:          "x = 5",
				Question:      "What is the value of x?",
				QuestionType:  QuestionType("unknown_type"),
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: 2,
				Explanation:   "x is assigned the value 5",
				Difficulty:    DifficultyBeginner,
				Topic:         TopicID("variables_types"),
				Language:      LanguagePython,
			},
			wantErr: true,
			errMsg:  "invalid question type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.response.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("LLMQuestionResponse.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" {
				if err == nil || err.Error() == "" {
					t.Errorf("LLMQuestionResponse.Validate() expected error message containing '%s', got nil or empty", tt.errMsg)
				} else if err.Error() != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("LLMQuestionResponse.Validate() error message = %v, want containing '%s'", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestLLMQuestionResponse_ToQuestion_MultipleChoice(t *testing.T) {
	tests := []struct {
		name      string
		response  LLMQuestionResponse
		wantErr   bool
		checkFunc func(*testing.T, *Question)
	}{
		{
			name: "valid multiple choice with int index",
			response: LLMQuestionResponse{
				Code:          "x = 5",
				Question:      "What is the value of x?",
				QuestionType:  QuestionTypeMultipleChoice,
				Options:       []string{"3", "4", "5", "6"},
				CorrectAnswer: 2,
				Explanation:   "x is assigned the value 5",
				Difficulty:    DifficultyBeginner,
				Topic:         TopicID("variables_types"),
				Language:      LanguagePython,
			},
			wantErr: false,
			checkFunc: func(t *testing.T, q *Question) {
				if q.CorrectAnswer != "2" {
					t.Errorf("ToQuestion() CorrectAnswer = %v, want '2'", q.CorrectAnswer)
				}
				if len(q.Options) != 4 {
					t.Errorf("ToQuestion() Options length = %v, want 4", len(q.Options))
				}
				if q.QuestionType != QuestionTypeMultipleChoice {
					t.Errorf("ToQuestion() QuestionType = %v, want %v", q.QuestionType, QuestionTypeMultipleChoice)
				}
			},
		},
		{
			name: "valid multiple choice with float64 index",
			response: LLMQuestionResponse{
				Code:          "x = 5",
				Question:      "What is the value of x?",
				QuestionType:  QuestionTypeMultipleChoice,
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: float64(1), // JSON числа парсятся как float64
				Explanation:   "x is assigned the value 5",
				Difficulty:    DifficultyBeginner,
				Topic:         TopicID("variables_types"),
				Language:      LanguagePython,
			},
			wantErr: false,
			checkFunc: func(t *testing.T, q *Question) {
				if q.CorrectAnswer != "1" {
					t.Errorf("ToQuestion() CorrectAnswer = %v, want '1'", q.CorrectAnswer)
				}
			},
		},
		{
			name: "valid multiple choice with string index",
			response: LLMQuestionResponse{
				Code:          "x = 5",
				Question:      "What is the value of x?",
				QuestionType:  QuestionTypeMultipleChoice,
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: "0",
				Explanation:   "x is assigned the value 5",
				Difficulty:    DifficultyBeginner,
				Topic:         TopicID("variables_types"),
				Language:      LanguagePython,
			},
			wantErr: false,
			checkFunc: func(t *testing.T, q *Question) {
				if q.CorrectAnswer != "0" {
					t.Errorf("ToQuestion() CorrectAnswer = %v, want '0'", q.CorrectAnswer)
				}
			},
		},
		{
			name: "multiple choice with invalid answer type",
			response: LLMQuestionResponse{
				Code:          "x = 5",
				Question:      "What is the value of x?",
				QuestionType:  QuestionTypeMultipleChoice,
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: []string{"invalid"}, // неверный тип
				Explanation:   "x is assigned the value 5",
				Difficulty:    DifficultyBeginner,
				Topic:         TopicID("variables_types"),
				Language:      LanguagePython,
			},
			wantErr: true,
		},
		{
			name: "multiple choice with out of range index",
			response: LLMQuestionResponse{
				Code:          "x = 5",
				Question:      "What is the value of x?",
				QuestionType:  QuestionTypeMultipleChoice,
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: 10, // индекс вне диапазона
				Explanation:   "x is assigned the value 5",
				Difficulty:    DifficultyBeginner,
				Topic:         TopicID("variables_types"),
				Language:      LanguagePython,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			question, err := tt.response.ToQuestion()
			if (err != nil) != tt.wantErr {
				t.Errorf("LLMQuestionResponse.ToQuestion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if question == nil {
					t.Fatal("LLMQuestionResponse.ToQuestion() returned nil question")
				}
				if tt.checkFunc != nil {
					tt.checkFunc(t, question)
				}
				// Проверяем, что все поля скопированы правильно
				if question.Code != tt.response.Code {
					t.Errorf("ToQuestion() Code = %v, want %v", question.Code, tt.response.Code)
				}
				if question.QuestionText != tt.response.Question {
					t.Errorf("ToQuestion() QuestionText = %v, want %v", question.QuestionText, tt.response.Question)
				}
				if question.Language != tt.response.Language {
					t.Errorf("ToQuestion() Language = %v, want %v", question.Language, tt.response.Language)
				}
				if question.Topic != tt.response.Topic {
					t.Errorf("ToQuestion() Topic = %v, want %v", question.Topic, tt.response.Topic)
				}
				if question.Difficulty != tt.response.Difficulty {
					t.Errorf("ToQuestion() Difficulty = %v, want %v", question.Difficulty, tt.response.Difficulty)
				}
				if question.ID == "" {
					t.Error("ToQuestion() should generate an ID")
				}
				if question.CreatedAt.IsZero() {
					t.Error("ToQuestion() should set CreatedAt")
				}
			}
		})
	}
}

func TestLLMQuestionResponse_ToQuestion_FreeText(t *testing.T) {
	tests := []struct {
		name      string
		response  LLMQuestionResponse
		wantErr   bool
		checkFunc func(*testing.T, *Question)
	}{
		{
			name: "valid free text with string answer",
			response: LLMQuestionResponse{
				Code:               "var s []int",
				Question:           "What is the type of s?",
				QuestionType:       QuestionTypeFreeText,
				CorrectAnswer:      "slice",
				AcceptableVariants: []string{"[]int", "slice of int"},
				CaseSensitive:      false,
				Explanation:        "s is a slice of integers",
				Difficulty:         DifficultyIntermediate,
				Topic:              TopicID("slices"),
				Language:           LanguageGo,
			},
			wantErr: false,
			checkFunc: func(t *testing.T, q *Question) {
				if q.CorrectAnswer != "slice" {
					t.Errorf("ToQuestion() CorrectAnswer = %v, want 'slice'", q.CorrectAnswer)
				}
				if len(q.AcceptableVariants) != 2 {
					t.Errorf("ToQuestion() AcceptableVariants length = %v, want 2", len(q.AcceptableVariants))
				}
				if q.CaseSensitive != false {
					t.Errorf("ToQuestion() CaseSensitive = %v, want false", q.CaseSensitive)
				}
				if q.QuestionType != QuestionTypeFreeText {
					t.Errorf("ToQuestion() QuestionType = %v, want %v", q.QuestionType, QuestionTypeFreeText)
				}
			},
		},
		{
			name: "free text with empty answer",
			response: LLMQuestionResponse{
				Code:         "var s []int",
				Question:     "What is the type of s?",
				QuestionType: QuestionTypeFreeText,
				CorrectAnswer: "",
				Explanation:  "s is a slice of integers",
				Difficulty:   DifficultyIntermediate,
				Topic:        TopicID("slices"),
				Language:     LanguageGo,
			},
			wantErr: true,
		},
		{
			name: "free text with non-string answer",
			response: LLMQuestionResponse{
				Code:         "var s []int",
				Question:     "What is the type of s?",
				QuestionType: QuestionTypeFreeText,
				CorrectAnswer: 123, // не строка
				Explanation:  "s is a slice of integers",
				Difficulty:   DifficultyIntermediate,
				Topic:        TopicID("slices"),
				Language:     LanguageGo,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			question, err := tt.response.ToQuestion()
			if (err != nil) != tt.wantErr {
				t.Errorf("LLMQuestionResponse.ToQuestion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if question == nil {
					t.Fatal("LLMQuestionResponse.ToQuestion() returned nil question")
				}
				if tt.checkFunc != nil {
					tt.checkFunc(t, question)
				}
			}
		})
	}
}

func TestLLMQuestionResponse_JSONParsing(t *testing.T) {
	// Тест для проверки преобразования строковых значений в enum типы при парсинге JSON
	tests := []struct {
		name     string
		jsonData string
		wantErr  bool
		checkFunc func(*testing.T, *LLMQuestionResponse)
	}{
		{
			name: "parse JSON with string enums",
			jsonData: `{
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
			checkFunc: func(t *testing.T, r *LLMQuestionResponse) {
				if r.QuestionType != QuestionTypeMultipleChoice {
					t.Errorf("JSON parsing QuestionType = %v, want %v", r.QuestionType, QuestionTypeMultipleChoice)
				}
				if r.Difficulty != DifficultyBeginner {
					t.Errorf("JSON parsing Difficulty = %v, want %v", r.Difficulty, DifficultyBeginner)
				}
				if r.Language != LanguagePython {
					t.Errorf("JSON parsing Language = %v, want %v", r.Language, LanguagePython)
				}
				if r.Topic != TopicID("variables_types") {
					t.Errorf("JSON parsing Topic = %v, want %v", r.Topic, TopicID("variables_types"))
				}
			},
		},
		{
			name: "parse JSON with free text",
			jsonData: `{
				"code": "var s []int",
				"question": "What is the type of s?",
				"question_type": "free_text",
				"correct_answer": "slice",
				"explanation": "s is a slice",
				"difficulty": "intermediate",
				"topic": "slices",
				"language": "go"
			}`,
			wantErr: false,
			checkFunc: func(t *testing.T, r *LLMQuestionResponse) {
				if r.QuestionType != QuestionTypeFreeText {
					t.Errorf("JSON parsing QuestionType = %v, want %v", r.QuestionType, QuestionTypeFreeText)
				}
				if r.Difficulty != DifficultyIntermediate {
					t.Errorf("JSON parsing Difficulty = %v, want %v", r.Difficulty, DifficultyIntermediate)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var response LLMQuestionResponse
			err := json.Unmarshal([]byte(tt.jsonData), &response)
			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if tt.checkFunc != nil {
					tt.checkFunc(t, &response)
				}
				// Проверяем валидность после парсинга
				if err := response.Validate(); err != nil {
					t.Errorf("Response.Validate() after JSON parsing error = %v", err)
				}
			}
		})
	}
}


