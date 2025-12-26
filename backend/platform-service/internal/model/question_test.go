package model

import (
	"strings"
	"testing"
	"time"
)

func TestQuestion_Validate(t *testing.T) {
	tests := []struct {
		name     string
		question Question
		wantErr  bool
		errMsg   string
	}{
		{
			name: "valid multiple choice question",
			question: Question{
				ID:            "test-id",
				Language:      LanguagePython,
				Topic:         TopicID("variables_types"),
				Difficulty:    DifficultyBeginner,
				QuestionType:  QuestionTypeMultipleChoice,
				Code:          "x = 5",
				QuestionText:  "What is the value of x?",
				Options:       []string{"3", "4", "5", "6"},
				CorrectAnswer: "2",
				Explanation:   "x is assigned the value 5",
				CreatedAt:     time.Now(),
			},
			wantErr: false,
		},
		{
			name: "valid free text question",
			question: Question{
				ID:                 "test-id-2",
				Language:           LanguageGo,
				Topic:              TopicID("slices"),
				Difficulty:         DifficultyIntermediate,
				QuestionType:       QuestionTypeFreeText,
				Code:               "var s []int",
				QuestionText:       "What is the type of s?",
				CorrectAnswer:      "slice",
				AcceptableVariants: []string{"[]int", "slice of int"},
				CaseSensitive:      false,
				Explanation:        "s is a slice of integers",
				CreatedAt:          time.Now(),
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			question: Question{
				Language:      LanguagePython,
				Topic:         TopicID("variables_types"),
				Difficulty:    DifficultyBeginner,
				QuestionType:  QuestionTypeMultipleChoice,
				Code:          "x = 5",
				QuestionText:  "What is the value of x?",
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: "2",
				Explanation:   "x is assigned the value 5",
				CreatedAt:     time.Now(),
			},
			wantErr: true,
			errMsg:  "question ID is required",
		},
		{
			name: "missing language",
			question: Question{
				ID:            "test-id",
				Topic:         TopicID("variables_types"),
				Difficulty:    DifficultyBeginner,
				QuestionType:  QuestionTypeMultipleChoice,
				Code:          "x = 5",
				QuestionText:  "What is the value of x?",
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: "2",
				Explanation:   "x is assigned the value 5",
				CreatedAt:     time.Now(),
			},
			wantErr: true,
			errMsg:  "language is required",
		},
		{
			name: "invalid language",
			question: Question{
				ID:            "test-id",
				Language:      Language("invalid"),
				Topic:         TopicID("variables_types"),
				Difficulty:    DifficultyBeginner,
				QuestionType:  QuestionTypeMultipleChoice,
				Code:          "x = 5",
				QuestionText:  "What is the value of x?",
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: "2",
				Explanation:   "x is assigned the value 5",
				CreatedAt:     time.Now(),
			},
			wantErr: true,
			errMsg:  "invalid language",
		},
		{
			name: "missing topic",
			question: Question{
				ID:            "test-id",
				Language:      LanguagePython,
				Difficulty:    DifficultyBeginner,
				QuestionType:  QuestionTypeMultipleChoice,
				Code:          "x = 5",
				QuestionText:  "What is the value of x?",
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: "2",
				Explanation:   "x is assigned the value 5",
				CreatedAt:     time.Now(),
			},
			wantErr: true,
			errMsg:  "topic is required",
		},
		{
			name: "invalid topic for language",
			question: Question{
				ID:            "test-id",
				Language:      LanguagePython,
				Topic:         TopicID("ownership"), // это тема для Rust, не для Python
				Difficulty:    DifficultyBeginner,
				QuestionType:  QuestionTypeMultipleChoice,
				Code:          "x = 5",
				QuestionText:  "What is the value of x?",
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: "2",
				Explanation:   "x is assigned the value 5",
				CreatedAt:     time.Now(),
			},
			wantErr: true,
			errMsg:  "invalid topic",
		},
		{
			name: "invalid difficulty",
			question: Question{
				ID:            "test-id",
				Language:      LanguagePython,
				Topic:         TopicID("variables_types"),
				Difficulty:    Difficulty("invalid"),
				QuestionType:  QuestionTypeMultipleChoice,
				Code:          "x = 5",
				QuestionText:  "What is the value of x?",
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: "2",
				Explanation:   "x is assigned the value 5",
				CreatedAt:     time.Now(),
			},
			wantErr: true,
			errMsg:  "invalid difficulty",
		},
		{
			name: "invalid question type",
			question: Question{
				ID:            "test-id",
				Language:      LanguagePython,
				Topic:         TopicID("variables_types"),
				Difficulty:    DifficultyBeginner,
				QuestionType:  QuestionType("invalid"),
				Code:          "x = 5",
				QuestionText:  "What is the value of x?",
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: "2",
				Explanation:   "x is assigned the value 5",
				CreatedAt:     time.Now(),
			},
			wantErr: true,
			errMsg:  "invalid question type",
		},
		{
			name: "missing code",
			question: Question{
				ID:            "test-id",
				Language:      LanguagePython,
				Topic:         TopicID("variables_types"),
				Difficulty:    DifficultyBeginner,
				QuestionType:  QuestionTypeMultipleChoice,
				QuestionText:  "What is the value of x?",
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: "2",
				Explanation:   "x is assigned the value 5",
				CreatedAt:     time.Now(),
			},
			wantErr: true,
			errMsg:  "code is required",
		},
		{
			name: "missing question text",
			question: Question{
				ID:            "test-id",
				Language:      LanguagePython,
				Topic:         TopicID("variables_types"),
				Difficulty:    DifficultyBeginner,
				QuestionType:  QuestionTypeMultipleChoice,
				Code:          "x = 5",
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: "2",
				Explanation:   "x is assigned the value 5",
				CreatedAt:     time.Now(),
			},
			wantErr: true,
			errMsg:  "question text is required",
		},
		{
			name: "missing correct answer",
			question: Question{
				ID:           "test-id",
				Language:     LanguagePython,
				Topic:        TopicID("variables_types"),
				Difficulty:   DifficultyBeginner,
				QuestionType: QuestionTypeMultipleChoice,
				Code:         "x = 5",
				QuestionText: "What is the value of x?",
				Options:      []string{"3", "4", "5"},
				Explanation:  "x is assigned the value 5",
				CreatedAt:    time.Now(),
			},
			wantErr: true,
			errMsg:  "correct answer is required",
		},
		{
			name: "missing explanation",
			question: Question{
				ID:            "test-id",
				Language:      LanguagePython,
				Topic:         TopicID("variables_types"),
				Difficulty:    DifficultyBeginner,
				QuestionType:  QuestionTypeMultipleChoice,
				Code:          "x = 5",
				QuestionText:  "What is the value of x?",
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: "2",
				CreatedAt:     time.Now(),
			},
			wantErr: true,
			errMsg:  "explanation is required",
		},
		{
			name: "multiple choice with too few options",
			question: Question{
				ID:            "test-id",
				Language:      LanguagePython,
				Topic:         TopicID("variables_types"),
				Difficulty:    DifficultyBeginner,
				QuestionType:  QuestionTypeMultipleChoice,
				Code:          "x = 5",
				QuestionText:  "What is the value of x?",
				Options:       []string{"3"},
				CorrectAnswer: "0",
				Explanation:   "x is assigned the value 5",
				CreatedAt:     time.Now(),
			},
			wantErr: true,
			errMsg:  "multiple choice question must have at least 2 options",
		},
		{
			name: "multiple choice with too many options",
			question: Question{
				ID:            "test-id",
				Language:      LanguagePython,
				Topic:         TopicID("variables_types"),
				Difficulty:    DifficultyBeginner,
				QuestionType:  QuestionTypeMultipleChoice,
				Code:          "x = 5",
				QuestionText:  "What is the value of x?",
				Options:       []string{"1", "2", "3", "4", "5", "6"},
				CorrectAnswer: "0",
				Explanation:   "x is assigned the value 5",
				CreatedAt:     time.Now(),
			},
			wantErr: true,
			errMsg:  "multiple choice question must have at most 5 options",
		},
		{
			name: "multiple choice with invalid answer index",
			question: Question{
				ID:            "test-id",
				Language:      LanguagePython,
				Topic:         TopicID("variables_types"),
				Difficulty:    DifficultyBeginner,
				QuestionType:  QuestionTypeMultipleChoice,
				Code:          "x = 5",
				QuestionText:  "What is the value of x?",
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: "5", // индекс вне диапазона
				Explanation:   "x is assigned the value 5",
				CreatedAt:     time.Now(),
			},
			wantErr: true,
			errMsg:  "correct answer index",
		},
		{
			name: "multiple choice with non-numeric answer",
			question: Question{
				ID:            "test-id",
				Language:      LanguagePython,
				Topic:         TopicID("variables_types"),
				Difficulty:    DifficultyBeginner,
				QuestionType:  QuestionTypeMultipleChoice,
				Code:          "x = 5",
				QuestionText:  "What is the value of x?",
				Options:       []string{"3", "4", "5"},
				CorrectAnswer: "abc", // не число
				Explanation:   "x is assigned the value 5",
				CreatedAt:     time.Now(),
			},
			wantErr: true,
			errMsg:  "correct answer must be a valid index",
		},
		{
			name: "free text with empty answer",
			question: Question{
				ID:            "test-id",
				Language:      LanguageGo,
				Topic:         TopicID("slices"),
				Difficulty:    DifficultyIntermediate,
				QuestionType:  QuestionTypeFreeText,
				Code:          "var s []int",
				QuestionText:  "What is the type of s?",
				CorrectAnswer: "",
				Explanation:   "s is a slice of integers",
				CreatedAt:     time.Now(),
			},
			wantErr: true,
			errMsg:  "correct answer is required",
		},
		{
			name: "unknown question type",
			question: Question{
				ID:            "test-id",
				Language:      LanguagePython,
				Topic:         TopicID("variables_types"),
				Difficulty:    DifficultyBeginner,
				QuestionType:  QuestionType("unknown_type"),
				Code:          "x = 5",
				QuestionText:  "What is the value of x?",
				CorrectAnswer: "2",
				Explanation:   "x is assigned the value 5",
				CreatedAt:     time.Now(),
			},
			wantErr: true,
			errMsg:  "invalid question type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.question.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Question.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" {
				if err == nil || err.Error() == "" {
					t.Errorf("Question.Validate() expected error message containing '%s', got nil or empty", tt.errMsg)
				} else if err.Error() != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Question.Validate() error message = %v, want containing '%s'", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestQuestion_GenerateID(t *testing.T) {
	question := &Question{}
	question.GenerateID()

	if question.ID == "" {
		t.Error("GenerateID() should generate a non-empty ID")
	}

	// Проверяем, что ID является валидным UUID
	if len(question.ID) != 36 {
		t.Errorf("GenerateID() should generate a UUID (36 chars), got length %d", len(question.ID))
	}

	// Генерируем еще раз и проверяем, что ID изменился
	oldID := question.ID
	question.GenerateID()
	if question.ID == oldID {
		t.Error("GenerateID() should generate a new ID each time")
	}
}
