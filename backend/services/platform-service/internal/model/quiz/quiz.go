package quiz

type Question struct {
	ID          string
	Language    Language
	Topic       string
	Difficulty  Difficulty
	Content     Content
	Explanation string
	Answer      Answer
}

type Content struct {
	Text string
	Code *string
}

type Language string

func (l Language) String() string {
	return string(l)
}

const (
	LanguageUnspecified Language = ""
	LanguagePython      Language = "python"
	LanguageJavaScript  Language = "javascript"
	LanguageGo          Language = "go"
	LanguageJava        Language = "java"
	LanguageCPP         Language = "cpp"
	LanguageRust        Language = "rust"
	LanguageTypeScript  Language = "typescript"
)

type Difficulty string

func (d Difficulty) String() string {
	return string(d)
}

const (
	DifficultyUnspecified  Difficulty = ""
	DifficultyBeginner     Difficulty = "beginner"
	DifficultyIntermediate Difficulty = "intermediate"
	DifficultyAdvanced     Difficulty = "advanced"
)

type AnswerType string

func (a AnswerType) String() string {
	return string(a)
}

const (
	AnswerTypeUnspecified    AnswerType = ""
	AnswerTypeMultipleChoice AnswerType = "multiple_choice"
	AnswerTypeFreeText       AnswerType = "free_text"
)

type Answer interface {
	AnswerType() AnswerType
}

type MultipleChoiceAnswer struct {
	Options        []string
	CorrectOptions []string
}

func (a *MultipleChoiceAnswer) AnswerType() AnswerType {
	return AnswerTypeMultipleChoice
}

type FreeTextAnswer struct {
	CorrectAnswers []string
}

func (a *FreeTextAnswer) AnswerType() AnswerType {
	return AnswerTypeFreeText
}
