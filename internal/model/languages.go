package model

// Language представляет язык программирования.
type Language string

const (
	// LanguagePython - Python язык программирования.
	LanguagePython Language = "python"
	// LanguageJavaScript - JavaScript язык программирования.
	LanguageJavaScript Language = "javascript"
	// LanguageGo - Go язык программирования.
	LanguageGo Language = "go"
	// LanguageJava - Java язык программирования.
	LanguageJava Language = "java"
	// LanguageCpp - C++ язык программирования.
	LanguageCpp Language = "cpp"
	// LanguageRust - Rust язык программирования.
	LanguageRust Language = "rust"
	// LanguageTypeScript - TypeScript язык программирования.
	LanguageTypeScript Language = "typescript"
)

// AllLanguages возвращает список всех поддерживаемых языков программирования.
func AllLanguages() []Language {
	return []Language{
		LanguagePython,
		LanguageJavaScript,
		LanguageGo,
		LanguageJava,
		LanguageCpp,
		LanguageRust,
		LanguageTypeScript,
	}
}

// SupportedLanguages содержит список всех поддерживаемых языков программирования (для обратной совместимости).
var SupportedLanguages = []string{
	string(LanguagePython),
	string(LanguageJavaScript),
	string(LanguageGo),
	string(LanguageJava),
	string(LanguageCpp),
	string(LanguageRust),
	string(LanguageTypeScript),
}

// LanguageNames содержит отображаемые названия языков программирования.
var LanguageNames = map[Language]string{
	LanguagePython:     "Python",
	LanguageJavaScript: "JavaScript",
	LanguageGo:         "Go",
	LanguageJava:       "Java",
	LanguageCpp:        "C++",
	LanguageRust:       "Rust",
	LanguageTypeScript: "TypeScript",
}

// String возвращает строковое представление языка программирования.
func (l Language) String() string {
	return string(l)
}

// IsValid проверяет, является ли значение валидным языком программирования.
func (l Language) IsValid() bool {
	switch l {
	case LanguagePython, LanguageJavaScript, LanguageGo, LanguageJava, LanguageCpp, LanguageRust, LanguageTypeScript:
		return true
	default:
		return false
	}
}

// GetName возвращает отображаемое название языка программирования.
func (l Language) GetName() string {
	if name, ok := LanguageNames[l]; ok {
		return name
	}
	return string(l)
}

// IsValidLanguage проверяет, является ли язык поддерживаемым (для обратной совместимости).
func IsValidLanguage(language string) bool {
	return Language(language).IsValid()
}
