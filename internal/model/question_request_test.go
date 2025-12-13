package model

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestGenerateQuestionRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request GenerateQuestionRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			request: GenerateQuestionRequest{
				Language:     LanguagePython,
				Topic:        TopicID("variables_types"),
				Difficulty:   DifficultyBeginner,
				QuestionType: QuestionTypeMultipleChoice,
			},
			wantErr: false,
		},
		{
			name: "missing language",
			request: GenerateQuestionRequest{
				Topic:        TopicID("variables_types"),
				Difficulty:   DifficultyBeginner,
				QuestionType: QuestionTypeMultipleChoice,
			},
			wantErr: true,
			errMsg:  "language is required",
		},
		{
			name: "invalid language",
			request: GenerateQuestionRequest{
				Language:     Language("invalid"),
				Topic:        TopicID("variables_types"),
				Difficulty:   DifficultyBeginner,
				QuestionType: QuestionTypeMultipleChoice,
			},
			wantErr: true,
			errMsg:  "unsupported language",
		},
		{
			name: "missing topic",
			request: GenerateQuestionRequest{
				Language:     LanguagePython,
				Difficulty:   DifficultyBeginner,
				QuestionType: QuestionTypeMultipleChoice,
			},
			wantErr: true,
			errMsg:  "topic is required",
		},
		{
			name: "invalid topic for language",
			request: GenerateQuestionRequest{
				Language:     LanguagePython,
				Topic:        TopicID("ownership"), // тема для Rust, не для Python
				Difficulty:   DifficultyBeginner,
				QuestionType: QuestionTypeMultipleChoice,
			},
			wantErr: true,
			errMsg:  "invalid topic",
		},
		{
			name: "invalid difficulty",
			request: GenerateQuestionRequest{
				Language:     LanguagePython,
				Topic:        TopicID("variables_types"),
				Difficulty:   Difficulty("invalid"),
				QuestionType: QuestionTypeMultipleChoice,
			},
			wantErr: true,
			errMsg:  "invalid difficulty",
		},
		{
			name: "invalid question type",
			request: GenerateQuestionRequest{
				Language:     LanguagePython,
				Topic:        TopicID("variables_types"),
				Difficulty:   DifficultyBeginner,
				QuestionType: QuestionType("invalid"),
			},
			wantErr: true,
			errMsg:  "invalid question type",
		},
		{
			name: "valid with all enum values",
			request: GenerateQuestionRequest{
				Language:     LanguageGo,
				Topic:        TopicID("slices"),
				Difficulty:   DifficultyAdvanced,
				QuestionType: QuestionTypeFreeText,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateQuestionRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" {
				if err == nil || err.Error() == "" {
					t.Errorf("GenerateQuestionRequest.Validate() expected error message containing '%s', got nil or empty", tt.errMsg)
				} else if err.Error() != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("GenerateQuestionRequest.Validate() error message = %v, want containing '%s'", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestGenerateQuestionRequest_JSONTags(t *testing.T) {
	// Тест для проверки JSON тегов структуры
	tests := []struct {
		name     string
		request  GenerateQuestionRequest
		jsonData string
		wantErr  bool
	}{
		{
			name: "marshal to JSON",
			request: GenerateQuestionRequest{
				Language:     LanguagePython,
				Topic:        TopicID("variables_types"),
				Difficulty:   DifficultyBeginner,
				QuestionType: QuestionTypeMultipleChoice,
			},
			wantErr: false,
		},
		{
			name: "unmarshal from JSON",
			jsonData: `{
				"language": "python",
				"topic": "variables_types",
				"difficulty": "beginner",
				"question_type": "multiple_choice"
			}`,
			wantErr: false,
		},
		{
			name: "unmarshal with free text",
			jsonData: `{
				"language": "go",
				"topic": "slices",
				"difficulty": "intermediate",
				"question_type": "free_text"
			}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.jsonData != "" {
				// Тест unmarshal
				var request GenerateQuestionRequest
				err := json.Unmarshal([]byte(tt.jsonData), &request)
				if (err != nil) != tt.wantErr {
					t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !tt.wantErr {
					// Проверяем, что поля заполнены правильно
					if request.Language == "" {
						t.Error("Unmarshal: Language should not be empty")
					}
					if request.Topic == "" {
						t.Error("Unmarshal: Topic should not be empty")
					}
					if request.Difficulty == "" {
						t.Error("Unmarshal: Difficulty should not be empty")
					}
					if request.QuestionType == "" {
						t.Error("Unmarshal: QuestionType should not be empty")
					}
					// Проверяем валидность
					if err := request.Validate(); err != nil {
						t.Errorf("Unmarshal: Validate() error = %v", err)
					}
				}
			} else {
				// Тест marshal
				data, err := json.Marshal(tt.request)
				if (err != nil) != tt.wantErr {
					t.Errorf("json.Marshal() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !tt.wantErr {
					if len(data) == 0 {
						t.Error("Marshal: should produce non-empty JSON")
					}
					// Проверяем, что можем обратно распарсить
					var request GenerateQuestionRequest
					if err := json.Unmarshal(data, &request); err != nil {
						t.Errorf("Marshal: round-trip unmarshal error = %v", err)
					}
					// Проверяем, что значения совпадают
					if request.Language != tt.request.Language {
						t.Errorf("Marshal: Language = %v, want %v", request.Language, tt.request.Language)
					}
					if request.Topic != tt.request.Topic {
						t.Errorf("Marshal: Topic = %v, want %v", request.Topic, tt.request.Topic)
					}
					if request.Difficulty != tt.request.Difficulty {
						t.Errorf("Marshal: Difficulty = %v, want %v", request.Difficulty, tt.request.Difficulty)
					}
					if request.QuestionType != tt.request.QuestionType {
						t.Errorf("Marshal: QuestionType = %v, want %v", request.QuestionType, tt.request.QuestionType)
					}
				}
			}
		})
	}
}

func TestGenerateQuestionRequest_Structure(t *testing.T) {
	// Тест для проверки структуры полей
	request := GenerateQuestionRequest{
		Language:     LanguagePython,
		Topic:        TopicID("variables_types"),
		Difficulty:   DifficultyBeginner,
		QuestionType: QuestionTypeMultipleChoice,
	}

	// Проверяем, что все поля имеют правильные типы
	if request.Language != LanguagePython {
		t.Errorf("Language type check failed: got %T, want Language", request.Language)
	}
	if request.Topic != TopicID("variables_types") {
		t.Errorf("Topic type check failed: got %T, want TopicID", request.Topic)
	}
	if request.Difficulty != DifficultyBeginner {
		t.Errorf("Difficulty type check failed: got %T, want Difficulty", request.Difficulty)
	}
	if request.QuestionType != QuestionTypeMultipleChoice {
		t.Errorf("QuestionType type check failed: got %T, want QuestionType", request.QuestionType)
	}
}


