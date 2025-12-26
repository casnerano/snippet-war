package model

// TopicID представляет идентификатор темы для языка программирования.
type TopicID string

// Topic представляет тему для языка программирования.
type Topic struct {
	ID   TopicID
	Name string
}

// LanguageTopics содержит карту тем для каждого языка программирования.
var LanguageTopics = map[Language][]Topic{
	LanguagePython: {
		{ID: TopicID("variables_types"), Name: "Переменные и типы данных"},
		{ID: TopicID("lists_arrays"), Name: "Списки и массивы"},
		{ID: TopicID("dictionaries"), Name: "Словари"},
		{ID: TopicID("functions"), Name: "Функции"},
		{ID: TopicID("closures"), Name: "Замыкания"},
		{ID: TopicID("decorators"), Name: "Декораторы"},
		{ID: TopicID("generators"), Name: "Генераторы"},
		{ID: TopicID("classes_oop"), Name: "Классы и ООП"},
		{ID: TopicID("exceptions"), Name: "Обработка исключений"},
		{ID: TopicID("context_managers"), Name: "Контекстные менеджеры"},
		{ID: TopicID("async_await"), Name: "Асинхронное программирование"},
	},
	LanguageJavaScript: {
		{ID: TopicID("variables_types"), Name: "Переменные и типы"},
		{ID: TopicID("arrays"), Name: "Массивы"},
		{ID: TopicID("objects"), Name: "Объекты"},
		{ID: TopicID("functions"), Name: "Функции"},
		{ID: TopicID("closures"), Name: "Замыкания"},
		{ID: TopicID("this_binding"), Name: "Контекст выполнения (this)"},
		{ID: TopicID("prototypes"), Name: "Прототипы"},
		{ID: TopicID("classes"), Name: "Классы (ES6+)"},
		{ID: TopicID("promises_async"), Name: "Промисы и async/await"},
		{ID: TopicID("event_loop"), Name: "Event Loop"},
		{ID: TopicID("destructuring"), Name: "Деструктуризация"},
		{ID: TopicID("modules"), Name: "Модули (ES6+)"},
	},
	LanguageGo: {
		{ID: TopicID("variables_types"), Name: "Переменные и типы"},
		{ID: TopicID("slices"), Name: "Срезы"},
		{ID: TopicID("maps"), Name: "Мапы"},
		{ID: TopicID("functions"), Name: "Функции"},
		{ID: TopicID("methods"), Name: "Методы"},
		{ID: TopicID("interfaces"), Name: "Интерфейсы"},
		{ID: TopicID("goroutines"), Name: "Горутины"},
		{ID: TopicID("channels"), Name: "Каналы"},
		{ID: TopicID("select"), Name: "Select statement"},
		{ID: TopicID("defer_panic_recover"), Name: "Defer, panic, recover"},
		{ID: TopicID("pointers"), Name: "Указатели"},
		{ID: TopicID("structs"), Name: "Структуры"},
	},
	LanguageJava: {
		{ID: TopicID("variables_types"), Name: "Переменные и типы"},
		{ID: TopicID("arrays_lists"), Name: "Массивы и списки"},
		{ID: TopicID("collections"), Name: "Коллекции (Set, Map)"},
		{ID: TopicID("methods"), Name: "Методы"},
		{ID: TopicID("classes_objects"), Name: "Классы и объекты"},
		{ID: TopicID("inheritance"), Name: "Наследование"},
		{ID: TopicID("interfaces"), Name: "Интерфейсы"},
		{ID: TopicID("generics"), Name: "Дженерики"},
		{ID: TopicID("exceptions"), Name: "Исключения"},
		{ID: TopicID("streams"), Name: "Streams API"},
		{ID: TopicID("lambda_expressions"), Name: "Lambda выражения"},
		{ID: TopicID("concurrency"), Name: "Многопоточность"},
	},
	LanguageCpp: {
		{ID: TopicID("variables_types"), Name: "Переменные и типы"},
		{ID: TopicID("pointers_references"), Name: "Указатели и ссылки"},
		{ID: TopicID("arrays_vectors"), Name: "Массивы и векторы"},
		{ID: TopicID("functions"), Name: "Функции"},
		{ID: TopicID("classes_objects"), Name: "Классы и объекты"},
		{ID: TopicID("inheritance"), Name: "Наследование"},
		{ID: TopicID("templates"), Name: "Шаблоны"},
		{ID: TopicID("smart_pointers"), Name: "Умные указатели"},
		{ID: TopicID("stl"), Name: "STL контейнеры и алгоритмы"},
		{ID: TopicID("move_semantics"), Name: "Move семантика"},
		{ID: TopicID("lambda"), Name: "Lambda выражения"},
		{ID: TopicID("multithreading"), Name: "Многопоточность"},
	},
	LanguageRust: {
		{ID: TopicID("variables_types"), Name: "Переменные и типы"},
		{ID: TopicID("ownership"), Name: "Владение (ownership)"},
		{ID: TopicID("borrowing"), Name: "Заимствование (borrowing)"},
		{ID: TopicID("lifetimes"), Name: "Время жизни"},
		{ID: TopicID("vectors"), Name: "Векторы"},
		{ID: TopicID("hashmaps"), Name: "HashMap"},
		{ID: TopicID("functions"), Name: "Функции"},
		{ID: TopicID("structs"), Name: "Структуры"},
		{ID: TopicID("enums"), Name: "Перечисления"},
		{ID: TopicID("pattern_matching"), Name: "Сопоставление с образцом"},
		{ID: TopicID("error_handling"), Name: "Обработка ошибок (Result, Option)"},
		{ID: TopicID("concurrency"), Name: "Многопоточность"},
	},
	LanguageTypeScript: {
		{ID: TopicID("types"), Name: "Типы"},
		{ID: TopicID("interfaces"), Name: "Интерфейсы"},
		{ID: TopicID("generics"), Name: "Дженерики"},
		{ID: TopicID("unions_intersections"), Name: "Объединения и пересечения типов"},
		{ID: TopicID("type_guards"), Name: "Защитники типов"},
		{ID: TopicID("decorators"), Name: "Декораторы"},
		{ID: TopicID("utility_types"), Name: "Утилитарные типы"},
		{ID: TopicID("modules"), Name: "Модули"},
		{ID: TopicID("async_promises"), Name: "Асинхронность"},
		{ID: TopicID("classes"), Name: "Классы"},
		{ID: TopicID("namespaces"), Name: "Пространства имен"},
	},
}

// String возвращает строковое представление идентификатора темы.
func (tid TopicID) String() string {
	return string(tid)
}

// GetTopicsForLanguage возвращает список тем для указанного языка.
func GetTopicsForLanguage(language Language) []Topic {
	topics, ok := LanguageTopics[language]
	if !ok {
		return nil
	}
	return topics
}

// IsValidTopic проверяет, является ли тема валидной для указанного языка.
func IsValidTopic(language Language, topic TopicID) bool {
	topics := GetTopicsForLanguage(language)
	if topics == nil {
		return false
	}

	for _, t := range topics {
		if t.ID == topic {
			return true
		}
	}

	return false
}

// IsValidTopicString проверяет, является ли тема валидной для указанного языка (для обратной совместимости).
func IsValidTopicString(language, topic string) bool {
	return IsValidTopic(Language(language), TopicID(topic))
}
