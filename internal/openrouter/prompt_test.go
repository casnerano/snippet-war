package openrouter

import (
	"strings"
	"testing"

	"github.com/casnerano/snippet-war/internal/model"
)

func TestBuildPrompt(t *testing.T) {
	tests := []struct {
		name   string
		req    *model.GenerateQuestionRequest
		checks []func(t *testing.T, prompt string)
	}{
		{
			name: "Python beginner multiple choice",
			req: &model.GenerateQuestionRequest{
				Language:     model.LanguagePython,
				Topic:        model.TopicID("variables_types"),
				Difficulty:   model.DifficultyBeginner,
				QuestionType: model.QuestionTypeMultipleChoice,
			},
			checks: []func(t *testing.T, prompt string){
				func(t *testing.T, prompt string) {
					if !strings.Contains(prompt, "Python") {
						t.Error("prompt should contain language name 'Python'")
					}
				},
				func(t *testing.T, prompt string) {
					if !strings.Contains(prompt, "Переменные и типы данных") {
						t.Error("prompt should contain topic name 'Переменные и типы данных'")
					}
				},
				func(t *testing.T, prompt string) {
					if !strings.Contains(prompt, "beginner") {
						t.Error("prompt should contain difficulty 'beginner'")
					}
				},
				func(t *testing.T, prompt string) {
					if !strings.Contains(prompt, "multiple_choice") {
						t.Error("prompt should contain answer type 'multiple_choice'")
					}
				},
				func(t *testing.T, prompt string) {
					desc := model.DifficultyBeginner.GetDescription()
					if !strings.Contains(prompt, desc) {
						t.Errorf("prompt should contain difficulty description: %s", desc)
					}
				},
				func(t *testing.T, prompt string) {
					if !strings.Contains(prompt, `"topic": "variables_types"`) {
						t.Error("prompt should contain topic ID in JSON format")
					}
				},
				func(t *testing.T, prompt string) {
					if !strings.Contains(prompt, `"language": "python"`) {
						t.Error("prompt should contain language ID in JSON format")
					}
				},
			},
		},
		{
			name: "JavaScript intermediate free text",
			req: &model.GenerateQuestionRequest{
				Language:     model.LanguageJavaScript,
				Topic:        model.TopicID("closures"),
				Difficulty:   model.DifficultyIntermediate,
				QuestionType: model.QuestionTypeFreeText,
			},
			checks: []func(t *testing.T, prompt string){
				func(t *testing.T, prompt string) {
					if !strings.Contains(prompt, "JavaScript") {
						t.Error("prompt should contain language name 'JavaScript'")
					}
				},
				func(t *testing.T, prompt string) {
					if !strings.Contains(prompt, "Замыкания") {
						t.Error("prompt should contain topic name 'Замыкания'")
					}
				},
				func(t *testing.T, prompt string) {
					if !strings.Contains(prompt, "intermediate") {
						t.Error("prompt should contain difficulty 'intermediate'")
					}
				},
				func(t *testing.T, prompt string) {
					if !strings.Contains(prompt, "free_text") {
						t.Error("prompt should contain answer type 'free_text'")
					}
				},
				func(t *testing.T, prompt string) {
					desc := model.DifficultyIntermediate.GetDescription()
					if !strings.Contains(prompt, desc) {
						t.Errorf("prompt should contain difficulty description: %s", desc)
					}
				},
			},
		},
		{
			name: "Go advanced multiple choice",
			req: &model.GenerateQuestionRequest{
				Language:     model.LanguageGo,
				Topic:        model.TopicID("goroutines"),
				Difficulty:   model.DifficultyAdvanced,
				QuestionType: model.QuestionTypeMultipleChoice,
			},
			checks: []func(t *testing.T, prompt string){
				func(t *testing.T, prompt string) {
					if !strings.Contains(prompt, "Go") {
						t.Error("prompt should contain language name 'Go'")
					}
				},
				func(t *testing.T, prompt string) {
					if !strings.Contains(prompt, "Горутины") {
						t.Error("prompt should contain topic name 'Горутины'")
					}
				},
				func(t *testing.T, prompt string) {
					if !strings.Contains(prompt, "advanced") {
						t.Error("prompt should contain difficulty 'advanced'")
					}
				},
				func(t *testing.T, prompt string) {
					desc := model.DifficultyAdvanced.GetDescription()
					if !strings.Contains(prompt, desc) {
						t.Errorf("prompt should contain difficulty description: %s", desc)
					}
				},
			},
		},
		{
			name: "Java intermediate free text",
			req: &model.GenerateQuestionRequest{
				Language:     model.LanguageJava,
				Topic:        model.TopicID("generics"),
				Difficulty:   model.DifficultyIntermediate,
				QuestionType: model.QuestionTypeFreeText,
			},
			checks: []func(t *testing.T, prompt string){
				func(t *testing.T, prompt string) {
					if !strings.Contains(prompt, "Java") {
						t.Error("prompt should contain language name 'Java'")
					}
				},
				func(t *testing.T, prompt string) {
					if !strings.Contains(prompt, "Дженерики") {
						t.Error("prompt should contain topic name 'Дженерики'")
					}
				},
			},
		},
		{
			name: "Rust advanced multiple choice",
			req: &model.GenerateQuestionRequest{
				Language:     model.LanguageRust,
				Topic:        model.TopicID("ownership"),
				Difficulty:   model.DifficultyAdvanced,
				QuestionType: model.QuestionTypeMultipleChoice,
			},
			checks: []func(t *testing.T, prompt string){
				func(t *testing.T, prompt string) {
					if !strings.Contains(prompt, "Rust") {
						t.Error("prompt should contain language name 'Rust'")
					}
				},
				func(t *testing.T, prompt string) {
					if !strings.Contains(prompt, "Владение (ownership)") {
						t.Error("prompt should contain topic name 'Владение (ownership)'")
					}
				},
			},
		},
		{
			name: "TypeScript beginner free text",
			req: &model.GenerateQuestionRequest{
				Language:     model.LanguageTypeScript,
				Topic:        model.TopicID("types"),
				Difficulty:   model.DifficultyBeginner,
				QuestionType: model.QuestionTypeFreeText,
			},
			checks: []func(t *testing.T, prompt string){
				func(t *testing.T, prompt string) {
					if !strings.Contains(prompt, "TypeScript") {
						t.Error("prompt should contain language name 'TypeScript'")
					}
				},
				func(t *testing.T, prompt string) {
					if !strings.Contains(prompt, "Типы") {
						t.Error("prompt should contain topic name 'Типы'")
					}
				},
			},
		},
		{
			name: "C++ intermediate multiple choice",
			req: &model.GenerateQuestionRequest{
				Language:     model.LanguageCpp,
				Topic:        model.TopicID("pointers_references"),
				Difficulty:   model.DifficultyIntermediate,
				QuestionType: model.QuestionTypeMultipleChoice,
			},
			checks: []func(t *testing.T, prompt string){
				func(t *testing.T, prompt string) {
					if !strings.Contains(prompt, "C++") {
						t.Error("prompt should contain language name 'C++'")
					}
				},
				func(t *testing.T, prompt string) {
					if !strings.Contains(prompt, "Указатели и ссылки") {
						t.Error("prompt should contain topic name 'Указатели и ссылки'")
					}
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt := BuildPrompt(tt.req)

			// Проверка, что промпт не пустой
			if prompt == "" {
				t.Fatal("prompt should not be empty")
			}

			// Проверка обязательных частей промпта
			requiredParts := []string{
				"Ты - эксперт по программированию",
				"Snippet War",
				"Параметры вопроса:",
				"Требования к вопросу:",
				"Требования к ответу:",
				"Формат ответа (JSON):",
				"Важно:",
			}

			for _, part := range requiredParts {
				if !strings.Contains(prompt, part) {
					t.Errorf("prompt should contain required part: %s", part)
				}
			}

			// Проверка отсутствия плейсхолдеров
			placeholders := []string{
				"{language}",
				"{topic}",
				"{difficulty}",
				"{difficulty_description}",
				"{answer_type}",
				"{topic_id}",
				"{language_id}",
			}

			for _, placeholder := range placeholders {
				if strings.Contains(prompt, placeholder) {
					t.Errorf("prompt should not contain placeholder: %s", placeholder)
				}
			}

			// Выполнение специфичных проверок
			for _, check := range tt.checks {
				check(t, prompt)
			}
		})
	}
}

func TestBuildPrompt_AllParametersSubstituted(t *testing.T) {
	req := &model.GenerateQuestionRequest{
		Language:     model.LanguagePython,
		Topic:        model.TopicID("functions"),
		Difficulty:   model.DifficultyBeginner,
		QuestionType: model.QuestionTypeMultipleChoice,
	}

	prompt := BuildPrompt(req)

	// Проверка подстановки языка
	if !strings.Contains(prompt, "Python") {
		t.Error("language name should be substituted")
	}
	if !strings.Contains(prompt, `"language": "python"`) {
		t.Error("language ID should be substituted in JSON")
	}

	// Проверка подстановки темы
	if !strings.Contains(prompt, "Функции") {
		t.Error("topic name should be substituted")
	}
	if !strings.Contains(prompt, `"topic": "functions"`) {
		t.Error("topic ID should be substituted in JSON")
	}

	// Проверка подстановки сложности
	if !strings.Contains(prompt, "beginner") {
		t.Error("difficulty should be substituted")
	}
	if !strings.Contains(prompt, model.DifficultyBeginner.GetDescription()) {
		t.Error("difficulty description should be substituted")
	}
	if !strings.Contains(prompt, `"difficulty": "beginner"`) {
		t.Error("difficulty should be substituted in JSON")
	}

	// Проверка подстановки типа ответа
	if !strings.Contains(prompt, "multiple_choice") {
		t.Error("answer type should be substituted")
	}
	if !strings.Contains(prompt, `"question_type": "multiple_choice"`) {
		t.Error("answer type should be substituted in JSON")
	}
}

func TestBuildPrompt_DifferentDifficulties(t *testing.T) {
	difficulties := []model.Difficulty{
		model.DifficultyBeginner,
		model.DifficultyIntermediate,
		model.DifficultyAdvanced,
	}

	for _, difficulty := range difficulties {
		t.Run(string(difficulty), func(t *testing.T) {
			req := &model.GenerateQuestionRequest{
				Language:     model.LanguagePython,
				Topic:        model.TopicID("variables_types"),
				Difficulty:   difficulty,
				QuestionType: model.QuestionTypeMultipleChoice,
			}

			prompt := BuildPrompt(req)

			// Проверка подстановки сложности
			if !strings.Contains(prompt, string(difficulty)) {
				t.Errorf("prompt should contain difficulty: %s", difficulty)
			}

			// Проверка подстановки описания сложности
			desc := difficulty.GetDescription()
			if desc != "" && !strings.Contains(prompt, desc) {
				t.Errorf("prompt should contain difficulty description: %s", desc)
			}
		})
	}
}

func TestBuildPrompt_DifferentQuestionTypes(t *testing.T) {
	questionTypes := []model.QuestionType{
		model.QuestionTypeMultipleChoice,
		model.QuestionTypeFreeText,
	}

	for _, questionType := range questionTypes {
		t.Run(string(questionType), func(t *testing.T) {
			req := &model.GenerateQuestionRequest{
				Language:     model.LanguagePython,
				Topic:        model.TopicID("variables_types"),
				Difficulty:   model.DifficultyBeginner,
				QuestionType: questionType,
			}

			prompt := BuildPrompt(req)

			// Проверка подстановки типа ответа
			if !strings.Contains(prompt, string(questionType)) {
				t.Errorf("prompt should contain question type: %s", questionType)
			}

			// Проверка подстановки в JSON
			expectedJSON := `"question_type": "` + string(questionType) + `"`
			if !strings.Contains(prompt, expectedJSON) {
				t.Errorf("prompt should contain question type in JSON: %s", expectedJSON)
			}
		})
	}
}
