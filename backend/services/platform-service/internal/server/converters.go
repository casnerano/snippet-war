package server

import (
	models "github.com/casnerano/snippet-war/internal/models/quiz"
	desc "github.com/casnerano/snippet-war/pkg/api/v1/quiz"
)

func QuestionToProto(question *models.Question) *desc.Question {
	if question == nil {
		return nil
	}

	pb := &desc.Question{
		Id:          question.ID,
		Language:    LanguageToProto(question.Language),
		Topic:       question.Topic,
		Difficulty:  DifficultyToProto(question.Difficulty),
		Explanation: question.Explanation,
		Content: &desc.Question_Content{
			Text: question.Content.Text,
			Code: question.Content.Code,
		},
	}

	switch answer := question.Answer.(type) {
	case *models.MultipleChoiceAnswer:
		pb.Answer = &desc.Question_MultipleChoice{
			MultipleChoice: &desc.Question_MultipleChoiceAnswer{
				Options:        answer.Options,
				CorrectOptions: answer.CorrectOptions,
			},
		}
	case *models.FreeTextAnswer:
		pb.Answer = &desc.Question_FreeText{
			FreeText: &desc.Question_FreeTextAnswer{
				CorrectAnswers: answer.CorrectAnswers,
			},
		}
	}

	return pb
}

func QuestionsToProto(questions []*models.Question) []*desc.Question {
	if len(questions) == 0 {
		return nil
	}

	pbQuestions := make([]*desc.Question, 0, len(questions))
	for _, q := range questions {
		pbQuestions = append(pbQuestions, QuestionToProto(q))
	}

	return pbQuestions
}

func ProtoToLanguage(language desc.Language) models.Language {
	switch language {
	case desc.Language_LANGUAGE_PYTHON:
		return models.LanguagePython
	case desc.Language_LANGUAGE_JAVASCRIPT:
		return models.LanguageJavaScript
	case desc.Language_LANGUAGE_GO:
		return models.LanguageGo
	case desc.Language_LANGUAGE_JAVA:
		return models.LanguageJava
	case desc.Language_LANGUAGE_CPP:
		return models.LanguageCPP
	case desc.Language_LANGUAGE_RUST:
		return models.LanguageRust
	case desc.Language_LANGUAGE_TYPESCRIPT:
		return models.LanguageTypeScript
	default:
		return models.LanguageUnspecified
	}
}

func ProtoToDifficulty(difficulty desc.Difficulty) models.Difficulty {
	switch difficulty {
	case desc.Difficulty_DIFFICULTY_BEGINNER:
		return models.DifficultyBeginner
	case desc.Difficulty_DIFFICULTY_INTERMEDIATE:
		return models.DifficultyIntermediate
	case desc.Difficulty_DIFFICULTY_ADVANCED:
		return models.DifficultyAdvanced
	default:
		return models.DifficultyUnspecified
	}
}

func LanguageToProto(language models.Language) desc.Language {
	switch language {
	case models.LanguagePython:
		return desc.Language_LANGUAGE_PYTHON
	case models.LanguageJavaScript:
		return desc.Language_LANGUAGE_JAVASCRIPT
	case models.LanguageGo:
		return desc.Language_LANGUAGE_GO
	case models.LanguageJava:
		return desc.Language_LANGUAGE_JAVA
	case models.LanguageCPP:
		return desc.Language_LANGUAGE_CPP
	case models.LanguageRust:
		return desc.Language_LANGUAGE_RUST
	case models.LanguageTypeScript:
		return desc.Language_LANGUAGE_TYPESCRIPT
	default:
		return desc.Language_LANGUAGE_UNSPECIFIED
	}
}

func DifficultyToProto(difficulty models.Difficulty) desc.Difficulty {
	switch difficulty {
	case models.DifficultyBeginner:
		return desc.Difficulty_DIFFICULTY_BEGINNER
	case models.DifficultyIntermediate:
		return desc.Difficulty_DIFFICULTY_INTERMEDIATE
	case models.DifficultyAdvanced:
		return desc.Difficulty_DIFFICULTY_ADVANCED
	default:
		return desc.Difficulty_DIFFICULTY_UNSPECIFIED
	}
}
