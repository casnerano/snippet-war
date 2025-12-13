package model

// Difficulty представляет уровень сложности вопроса.
type Difficulty string

const (
	// DifficultyBeginner - базовый уровень сложности.
	DifficultyBeginner Difficulty = "beginner"
	// DifficultyIntermediate - средний уровень сложности.
	DifficultyIntermediate Difficulty = "intermediate"
	// DifficultyAdvanced - продвинутый уровень сложности.
	DifficultyAdvanced Difficulty = "advanced"
)

// AllDifficulties возвращает список всех поддерживаемых уровней сложности.
func AllDifficulties() []Difficulty {
	return []Difficulty{
		DifficultyBeginner,
		DifficultyIntermediate,
		DifficultyAdvanced,
	}
}

// String возвращает строковое представление уровня сложности.
func (d Difficulty) String() string {
	return string(d)
}

// IsValid проверяет, является ли значение валидным уровнем сложности.
func (d Difficulty) IsValid() bool {
	switch d {
	case DifficultyBeginner, DifficultyIntermediate, DifficultyAdvanced:
		return true
	default:
		return false
	}
}

// DifficultyDescriptions содержит описания уровней сложности.
var DifficultyDescriptions = map[Difficulty]string{
	DifficultyBeginner:     "Базовые операции и синтаксис. Простые типы данных, базовые структуры данных, простые условия и циклы, простые функции без сложной логики, базовые операции со строками и числами.",
	DifficultyIntermediate: "Более сложные структуры данных, вложенные циклы и условия, функции высшего порядка, работа с коллекциями, базовое ООП, обработка исключений, базовые паттерны проектирования.",
	DifficultyAdvanced:     "Сложные алгоритмы и оптимизация, продвинутые концепции языка, конкурентное/параллельное программирование, продвинутые паттерны проектирования, неочевидное поведение языка, оптимизация производительности, работа с памятью и указателями.",
}

// GetDescription возвращает описание уровня сложности.
func (d Difficulty) GetDescription() string {
	if desc, ok := DifficultyDescriptions[d]; ok {
		return desc
	}
	return ""
}
