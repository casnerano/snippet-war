# Тестовый скрипт для генерации вопросов

Этот скрипт позволяет протестировать работу генерации вопросов через ProxyAPI (OpenAI-совместимый API).

## Настройка

Перед запуском необходимо установить следующие переменные окружения:

```bash
# Обязательные переменные
export PROXYAPI_API_KEY="your-api-key-here"
export PROXYAPI_MODEL="gpt-4o"  # или другая модель OpenAI (например, "gpt-3.5-turbo", "gpt-4o-mini")

# Опциональные переменные (имеют значения по умолчанию)
export PROXYAPI_BASE_URL="https://api.proxyapi.ru/openai/v1"  # по умолчанию
export PROXYAPI_TIMEOUT="30s"  # по умолчанию
export PROXYAPI_MAX_TOKENS="2000"  # по умолчанию
```

### Получение API ключа ProxyAPI

1. Зарегистрируйтесь на [proxyapi.ru](https://proxyapi.ru)
2. Перейдите в раздел "Ключи API"
3. Создайте новый ключ (обратите внимание: ключ можно увидеть только один раз при создании)

### Доступные модели

ProxyAPI поддерживает все модели OpenAI через OpenAI-совместимый API:
- `gpt-4o` - GPT-4o (рекомендуется)
- `gpt-4o-mini` - GPT-4o Mini (быстрая и экономичная)
- `gpt-4-turbo` - GPT-4 Turbo
- `gpt-4` - GPT-4
- `gpt-3.5-turbo` - GPT-3.5 Turbo

Полный список доступных моделей можно найти в [документации ProxyAPI](https://proxyapi.ru/docs/overview).

## Сборка

Из корня проекта:

```bash
go build -o bin/test-question ./cmd/test-question
```

Или из директории `cmd/test-question`:

```bash
go build -o ../../bin/test-question .
```

## Запуск

```bash
./bin/test-question
```

Или из директории backend:

```bash
../../bin/test-question
```

## Параметры запроса

В текущей версии скрипт использует следующие параметры:

- **Язык**: `python`
- **Тема**: `variables_types` (Переменные и типы данных)
- **Сложность**: `beginner` (Начальный уровень)
- **Тип вопроса**: `multiple_choice` (Множественный выбор)

Эти параметры можно изменить в файле `main.go` в функции `main()`:

```go
req := &model.GenerateQuestionRequest{
    Language:     model.LanguagePython,      // Изменить язык
    Topic:        model.TopicID("variables_types"),  // Изменить тему
    Difficulty:   model.DifficultyBeginner,   // Изменить сложность
    QuestionType: model.QuestionTypeMultipleChoice,  // Изменить тип вопроса
}
```

## Доступные значения

### Языки:
- `LanguagePython` - Python
- `LanguageJavaScript` - JavaScript
- `LanguageGo` - Go
- `LanguageJava` - Java
- `LanguageCpp` - C++
- `LanguageRust` - Rust
- `LanguageTypeScript` - TypeScript

### Уровни сложности:
- `DifficultyBeginner` - Начальный
- `DifficultyIntermediate` - Средний
- `DifficultyAdvanced` - Продвинутый

### Типы вопросов:
- `QuestionTypeMultipleChoice` - Множественный выбор
- `QuestionTypeFreeText` - Свободный текстовый ответ

### Примеры тем для Python:
- `variables_types` - Переменные и типы данных
- `lists_arrays` - Списки и массивы
- `dictionaries` - Словари
- `functions` - Функции
- `classes_oop` - Классы и ООП
- и другие (см. `backend/internal/model/topics.go`)

## Пример вывода

Скрипт выводит:
1. Информацию о процессе генерации
2. Красиво отформатированный вопрос
3. JSON представление вопроса

