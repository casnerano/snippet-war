# План реализации: Интеграция Backend с LLM через OpenRouter

> **Основан на roadmap:** `.cursor/plans/llm_integration_roadmap_2b7365cf.plan.md`

## Цель

Реализовать backend-часть для генерации вопросов с помощью LLM через OpenRouter API. Система должна поддерживать 7 языков программирования, различные темы, уровни сложности и типы ответов.

---

## Этап 1: Подготовка проекта и конфигурация

### Задача 1.1: Создать пример файла конфигурации

**Файл:** `backend/.env.example`

**Описание:** Создать файл с переменными окружения для OpenRouter API.

**Содержимое (см. раздел "1. Конфигурация OpenRouter" roadmap):**

```env
# OpenRouter API Configuration
OPENROUTER_API_KEY=your_api_key_here
OPENROUTER_MODEL=deepseek/deepseek-chat
OPENROUTER_BASE_URL=https://openrouter.ai/api/v1
OPENROUTER_TIMEOUT=30s
OPENROUTER_MAX_TOKENS=2000
```

**Переменные:**
- `OPENROUTER_API_KEY` - API ключ OpenRouter (обязательная)
- `OPENROUTER_MODEL` - модель LLM (обязательная)
- `OPENROUTER_BASE_URL` - базовый URL API (опционально, по умолчанию `https://openrouter.ai/api/v1`)
- `OPENROUTER_TIMEOUT` - таймаут запросов (опционально, по умолчанию `30s`)
- `OPENROUTER_MAX_TOKENS` - максимальное количество токенов в ответе (опционально, по умолчанию `2000`)

---

### Задача 1.2: Добавить зависимость go-openrouter в go.mod

**Файл:** `backend/go.mod`

**Описание:** Добавить библиотеку `github.com/revrost/go-openrouter` версии v1.1.5 для работы с OpenRouter API.

**Команда:**
```bash
cd backend
go get github.com/revrost/go-openrouter@v1.1.5
```

---

### Задача 1.3: Создать структуру конфигурации

**Файл:** `backend/internal/config/config.go`

**Описание:** Создать пакет для загрузки конфигурации из переменных окружения.

**Структура:**
```go
package config

type OpenRouterConfig struct {
    APIKey   string
    Model    string
    BaseURL  string
    Timeout  time.Duration
    MaxTokens int
}

type Config struct {
    OpenRouter OpenRouterConfig
}

func Load() (*Config, error)
```

**Реализация:**
- Использовать `os.Getenv()` для получения переменных окружения
- Применить значения по умолчанию для опциональных параметров
- Парсить `OPENROUTER_TIMEOUT` как `time.Duration` (например, `time.ParseDuration()`)
- Валидировать обязательные поля (API_KEY, MODEL)
- Возвращать понятные ошибки при отсутствии обязательных полей

---

## Этап 2: Модели данных

### Задача 2.1: Создать модель Question

**Файл:** `backend/internal/model/question.go`

**Описание:** Создать структуру данных для хранения вопроса (см. раздел "7. Структура данных для хранения" roadmap).

**Примечание:** Для `Difficulty` и `QuestionType` используются enum типы (см. задачи 3.3 и 3.4).

**Модель:**
```go
package model

import "time"

type Question struct {
    ID                string       `json:"id"`
    Language          string       `json:"language"`
    Topic             string       `json:"topic"`
    Difficulty        Difficulty   `json:"difficulty"` // enum: beginner, intermediate, advanced
    QuestionType      QuestionType `json:"question_type"` // enum: multiple_choice, free_text
    Code              string       `json:"code"`
    QuestionText      string       `json:"question"`
    Options           []string     `json:"options,omitempty"` // для multiple_choice
    CorrectAnswer     string       `json:"correct_answer"` // строка или индекс
    AcceptableVariants []string     `json:"acceptable_variants,omitempty"` // для free_text
    CaseSensitive     bool         `json:"case_sensitive,omitempty"` // для free_text
    Explanation       string       `json:"explanation"`
    CreatedAt         time.Time    `json:"created_at"`
}
```

**Дополнительно:**
- Добавить метод `Validate()` для валидации структуры (включая проверку валидности enum значений)
- Добавить метод `GenerateID()` или использовать UUID при создании

---

### Задача 2.2: Создать DTO для запроса генерации вопроса

**Файл:** `backend/internal/model/question_request.go`

**Описание:** Создать структуры для запроса на генерацию вопроса.

**Структуры:**
```go
package model

type GenerateQuestionRequest struct {
    Language     Language       `json:"language" binding:"required"`
    Topic        TopicID       `json:"topic" binding:"required"`
    Difficulty   Difficulty   `json:"difficulty" binding:"required"`
    QuestionType QuestionType `json:"question_type" binding:"required"`
}
```

**Валидация:**
- `Language` - должен быть из списка поддерживаемых языков (см. раздел "2. Список языков программирования")
- `Topic` - должен быть из списка тем для выбранного языка (см. раздел "3. Темы для каждого языка")
- `Difficulty` - должен быть валидным значением enum `Difficulty` (использовать метод `IsValid()`)
- `QuestionType` - должен быть валидным значением enum `QuestionType` (использовать метод `IsValid()`)

**Примечание:** Для валидации enum значений можно использовать кастомный валидатор или проверку через методы `IsValid()`.

---

### Задача 2.3: Создать DTO для ответа LLM

**Файл:** `backend/internal/model/llm_response.go`

**Описание:** Создать структуры для парсинга ответа от LLM.

**Структура (соответствует JSON из промпта, раздел "6. Промпт для генерации вопроса"):**
```go
package model

type LLMQuestionResponse struct {
    Code              string       `json:"code"`
    Question          string       `json:"question"`
    QuestionType      QuestionType `json:"question_type"` // enum
    Options           []string     `json:"options,omitempty"`
    CorrectAnswer     interface{}  `json:"correct_answer"` // может быть string или int
    AcceptableVariants []string    `json:"acceptable_variants,omitempty"`
    CaseSensitive     bool         `json:"case_sensitive,omitempty"`
    Explanation       string       `json:"explanation"`
    Difficulty        Difficulty   `json:"difficulty"` // enum
    Topic             TopicID       `json:"topic"`
    Language          Language       `json:"language"`
}
```

**Методы:**
- `ToQuestion() (*Question, error)` - преобразование в модель Question (с валидацией enum значений)
- `Validate() error` - валидация ответа от LLM (включая проверку валидности enum значений)

**Примечание:** При парсинге JSON из LLM ответа, строковые значения для `Difficulty` и `QuestionType` нужно преобразовать в соответствующие enum типы с проверкой валидности.

---

### Задача 2.4: Создать unit tests для моделей

**Файл:** `backend/internal/model/question_test.go`, `backend/internal/model/question_request_test.go`, `backend/internal/model/llm_response_test.go`

**Описание:** Создать тесты для моделей данных.

**Тесты для `Question`:**
- Валидация структуры `Question.Validate()`
- Проверка обязательных полей
- Проверка корректности типов данных
- Проверка валидации enum значений (`Difficulty`, `QuestionType`)

**Тесты для `LLMQuestionResponse`:**
- Метод `ToQuestion()` для `multiple_choice` (преобразование индекса в строку, валидация enum)
- Метод `ToQuestion()` для `free_text` (преобразование строки, валидация enum)
- Метод `Validate()` - проверка всех обязательных полей (включая валидацию enum значений)
- Проверка соответствия типа ответа структуре данных
- Проверка преобразования строковых значений в enum типы при парсинге JSON

**Тесты для `GenerateQuestionRequest`:**
- Проверка структуры полей
- Проверка JSON тегов

---

## Этап 3: Константы и валидация

### Задача 3.1: Определить константы языков программирования

**Файл:** `backend/internal/model/languages.go`

**Описание:** Определить константы для поддерживаемых языков (см. раздел "2. Список языков программирования" roadmap).

**Константы:**
```go
package model

const (
    LanguagePython     = "python"
    LanguageJavaScript = "javascript"
    LanguageGo         = "go"
    LanguageJava       = "java"
    LanguageCpp        = "cpp"
    LanguageRust       = "rust"
    LanguageTypeScript = "typescript"
)

var SupportedLanguages = []string{
    LanguagePython,
    LanguageJavaScript,
    LanguageGo,
    LanguageJava,
    LanguageCpp,
    LanguageRust,
    LanguageTypeScript,
}

var LanguageNames = map[string]string{
    LanguagePython:     "Python",
    LanguageJavaScript: "JavaScript",
    LanguageGo:         "Go",
    LanguageJava:       "Java",
    LanguageCpp:        "C++",
    LanguageRust:       "Rust",
    LanguageTypeScript: "TypeScript",
}
```

---

### Задача 3.2: Определить константы тем для каждого языка

**Файл:** `backend/internal/model/topics.go`

**Описание:** Определить карту тем для каждого языка (см. раздел "3. Темы для каждого языка" roadmap).

**Структура:**
```go
package model

var LanguageTopics = map[string][]Topic{
    LanguagePython: {
        {ID: "variables_types", Name: "Переменные и типы данных"},
        {ID: "lists_arrays", Name: "Списки и массивы"},
        {ID: "dictionaries", Name: "Словари"},
        {ID: "functions", Name: "Функции"},
        {ID: "closures", Name: "Замыкания"},
        {ID: "decorators", Name: "Декораторы"},
        {ID: "generators", Name: "Генераторы"},
        {ID: "classes_oop", Name: "Классы и ООП"},
        {ID: "exceptions", Name: "Обработка исключений"},
        {ID: "context_managers", Name: "Контекстные менеджеры"},
        {ID: "async_await", Name: "Асинхронное программирование"},
    },
    // ... аналогично для остальных языков
}

type Topic struct {
    ID   string
    Name string
}

func GetTopicsForLanguage(language string) []Topic
func IsValidTopic(language, topic string) bool
```

**Примечание:** Реализовать полный список тем для всех 7 языков согласно roadmap.

---

### Задача 3.5: Создать валидатор для запросов

**Файл:** `backend/internal/validator/validator.go`

**Описание:** Создать функции валидации для запросов на генерацию вопросов.

**Функции:**
```go
package validator

import "github.com/casnerano/snippet-war/backend/internal/model"

func ValidateLanguage(language string) error
func ValidateTopic(language, topic string) error
func ValidateDifficulty(difficulty model.Difficulty) error
func ValidateQuestionType(questionType model.QuestionType) error
func ValidateGenerateRequest(req *model.GenerateQuestionRequest) error
```

**Реализация:**
- Использовать enum типы и их методы `IsValid()` из пакета `model`
- Для `ValidateDifficulty` и `ValidateQuestionType` использовать методы `IsValid()` соответствующих enum типов
- Возвращать понятные ошибки валидации
- Использовать `errors.New()` или кастомные типы ошибок

---

## Этап 4: Клиент OpenRouter

### Задача 4.1: Создать интерфейс клиента OpenRouter

**Файл:** `backend/internal/openrouter/client.go`

**Описание:** Создать интерфейс для работы с OpenRouter API.

**Интерфейс:**
```go
package openrouter

import (
    "context"
    "github.com/casnerano/snippet-war/backend/internal/model"
)

type Client interface {
    GenerateQuestion(ctx context.Context, prompt string) (string, error)
}

type client struct {
    apiKey   string
    model    string
    baseURL  string
    timeout  time.Duration
    maxTokens int
}
```

---

### Задача 4.2: Реализовать клиент OpenRouter

**Файл:** `backend/internal/openrouter/client.go` (продолжение)

**Описание:** Реализовать методы клиента с использованием библиотеки `github.com/revrost/go-openrouter`.

**Методы:**
```go
func NewClient(cfg *config.OpenRouterConfig) Client

func (c *client) GenerateQuestion(ctx context.Context, prompt string) (string, error)
```

**Реализация:**
- Использовать библиотеку `github.com/revrost/go-openrouter`
- Создавать контекст с таймаутом из конфигурации
- Отправлять запрос к OpenRouter API с промптом
- Обрабатывать ошибки API (429, 500, и т.д.)
- Возвращать только текст ответа от LLM (без дополнительных метаданных)

**Обработка ошибок:**
- Таймауты
- Ошибки API (ошибки аутентификации, rate limiting, и т.д.)
- Ошибки сети
- Обертывать ошибки с контекстом (`fmt.Errorf("...: %w", err)`)

---

### Задача 4.3: Создать конструктор промпта

**Файл:** `backend/internal/openrouter/prompt.go`

**Описание:** Создать функцию для генерации промпта на основе параметров (см. раздел "6. Промпт для генерации вопроса" roadmap).

**Функция:**
```go
package openrouter

import "github.com/casnerano/snippet-war/backend/internal/model"

func BuildPrompt(req *model.GenerateQuestionRequest) string
```

**Реализация:**
- Использовать шаблон промпта из roadmap (строки 262-302)
- Подставлять параметры:
  - `{language}` - название языка из `LanguageNames`
  - `{topic}` - ID темы
  - `{difficulty}` - уровень сложности
  - `{difficulty_description}` - описание уровня сложности из `DifficultyDescriptions`
  - `{answer_type}` - тип ответа
- Использовать `fmt.Sprintf()` или `strings.ReplaceAll()` для подстановки

**Промпт (из roadmap):**
```
Ты - эксперт по программированию на языке {language}. Твоя задача - создать вопрос для игры "Snippet War", где игроки угадывают вывод кода.

Параметры вопроса:
- Язык программирования: {language}
- Тема: {topic}
- Уровень сложности: {difficulty} ({difficulty_description})
- Тип ответа: {answer_type}

Требования к вопросу:
1. Создай короткий, но понятный фрагмент кода (не более 30 строк)
2. Код должен демонстрировать концепцию из указанной темы
3. Сложность должна соответствовать уровню {difficulty}
4. Вопрос должен быть интересным и обучающим
5. Код должен быть валидным и исполняемым

Требования к ответу:
- Если тип ответа "multiple_choice": создай 4 варианта ответа, где только один правильный. Варианты должны быть правдоподобными, но отличаться от правильного.
- Если тип ответа "free_text": укажи точный правильный ответ и возможные варианты написания (с учетом регистра, пробелов и т.д.)

Формат ответа (JSON):
{
  "code": "фрагмент кода на языке {language}",
  "question": "Что выведет этот код? (или другой вопрос)",
  "question_type": "{answer_type}",
  "options": ["вариант1", "вариант2", "вариант3", "вариант4"], // только для multiple_choice
  "correct_answer": "правильный ответ" или индекс (для multiple_choice),
  "acceptable_variants": ["вариант1", "вариант2"], // только для free_text, опционально
  "case_sensitive": false, // только для free_text
  "explanation": "подробное объяснение, почему правильный ответ именно такой. Объясни, как работает код, какие концепции демонстрируются.",
  "difficulty": "{difficulty}",
  "topic": "{topic}",
  "language": "{language}"
}

Важно:
- Код должен быть правильно отформатирован
- Объяснение должно быть понятным и обучающим
- Для уровня Beginner код должен быть простым и понятным
- Для уровня Advanced можно использовать неочевидные особенности языка
- Всегда возвращай валидный JSON без дополнительного текста
```

---

### Задача 4.4: Создать unit tests для клиента OpenRouter

**Файл:** `backend/internal/openrouter/client_test.go`

**Описание:** Создать тесты для клиента OpenRouter (с моками библиотеки).

**Тесты:**
- Успешный запрос к OpenRouter API
- Обработка ошибок API (429, 500, и т.д.)
- Обработка таймаутов
- Обработка ошибок сети
- Проверка корректности передачи параметров (model, prompt, max_tokens)

**Примечание:** Для unit тестов использовать mock для библиотеки `go-openrouter`.

---

### Задача 4.5: Создать unit tests для конструктора промпта

**Файл:** `backend/internal/openrouter/prompt_test.go`

**Описание:** Создать тесты для функции `BuildPrompt`.

**Тесты:**
- Корректная подстановка всех параметров (language, topic, difficulty, answer_type)
- Корректная подстановка описания уровня сложности
- Проверка включения всех обязательных частей промпта
- Тесты для разных языков, тем, уровней сложности и типов ответов

---

## Этап 5: Сервис генерации вопросов

---

### Задача 5.2: Реализовать сервис генерации вопросов

**Файл:** `backend/internal/service/question_service.go` (продолжение)

**Описание:** Реализовать логику генерации вопросов с использованием клиента OpenRouter.

**Структура:**
```go
type QuestionService struct {
    openRouterClient openrouter.Client
    logger          Logger // если есть система логирования
}

func NewQuestionService(client openrouter.Client, logger Logger) QuestionService
```

**Метод `GenerateQuestion`:**
1. Валидировать запрос (`validator.ValidateGenerateRequest`)
2. Построить промпт (`openrouter.BuildPrompt`)
3. Отправить запрос к OpenRouter API (`openRouterClient.GenerateQuestion`)
4. Распарсить JSON ответ от LLM в `LLMQuestionResponse`
5. Преобразовать `LLMQuestionResponse` в `Question` (метод `ToQuestion`)
6. Валидировать полученный вопрос
7. Установить `ID` и `CreatedAt` для вопроса
8. Вернуть вопрос

**Обработка ошибок:**
- Валидация входных данных
- Ошибки парсинга JSON (может быть невалидный JSON от LLM)
- Ошибки валидации ответа от LLM
- Ошибки API OpenRouter

---

### Задача 5.3: Создать парсер ответа LLM

**Файл:** `backend/internal/service/llm_parser.go`

**Описание:** Создать функции для парсинга и валидации ответа от LLM.

**Функции:**
```go
package service

import (
    "encoding/json"
    "github.com/casnerano/snippet-war/backend/internal/model"
)

func ParseLLMResponse(jsonStr string) (*model.LLMQuestionResponse, error)
func ValidateLLMResponse(resp *model.LLMQuestionResponse, req *model.GenerateQuestionRequest) error
```

**Реализация `ParseLLMResponse`:**
- Использовать `json.Unmarshal()` для парсинга JSON
- Обрабатывать ошибки парсинга (может быть невалидный JSON)
- Удалять markdown code blocks (если LLM обернул ответ в ```json ... ```)

**Реализация `ValidateLLMResponse`:**
- Проверить наличие всех обязательных полей
- Проверить соответствие языка, темы, сложности и типа ответа запросу
- Проверить структуру в зависимости от типа ответа:
  - Для `multiple_choice`: наличие `options` (минимум 2, максимум 5), валидный `correct_answer` (число в диапазоне индексов)
  - Для `free_text`: наличие `correct_answer` (непустая строка)
- Проверить, что `code` не пустой
- Проверить, что `question` не пустой
- Проверить, что `explanation` не пустой

**Метод `ToQuestion` в `LLMQuestionResponse`:**
```go
func (r *LLMQuestionResponse) ToQuestion() (*model.Question, error) {
    // Преобразовать correct_answer (может быть string или int) в string
    // Для multiple_choice: correct_answer должен быть индексом
    // Для free_text: correct_answer - строка
}
```

---

### Задача 5.4: Создать unit tests для парсера LLM ответа

**Файл:** `backend/internal/service/llm_parser_test.go`

**Описание:** Создать тесты для парсинга ответов от LLM.

**Тесты для `ParseLLMResponse`:**
- Парсинг валидного JSON ответа (multiple_choice)
- Парсинг валидного JSON ответа (free_text)
- Парсинг JSON с markdown code blocks (удаление ```json ... ```)
- Обработка невалидного JSON
- Обработка пустого ответа
- Обработка ответа с отсутствующими полями

**Тесты для `ValidateLLMResponse`:**
- Проверка всех обязательных полей
- Проверка соответствия языка, темы, сложности и типа ответа запросу
- Валидация структуры для `multiple_choice`:
  - Наличие `options` (минимум 2, максимум 5)
  - Валидный `correct_answer` (число в диапазоне индексов)
- Валидация структуры для `free_text`:
  - Наличие `correct_answer` (непустая строка)
- Проверка непустых полей: `code`, `question`, `explanation`

---

### Задача 5.5: Создать integration tests для сервиса

**Файл:** `backend/internal/service/question_service_test.go`

**Описание:** Создать интеграционные тесты для сервиса генерации вопросов.

**Тесты:**
- Успешная генерация вопроса типа `multiple_choice`
- Успешная генерация вопроса типа `free_text`
- Валидация входных данных (невалидный язык, тема, и т.д.)
- Обработка ошибок от OpenRouter API (mock клиента)
- Обработка невалидного JSON от LLM
- Проверка корректности установки `ID` и `CreatedAt`

**Примечание:** Использовать mock клиента OpenRouter для контроля возвращаемых данных.

---

### Задача 3.6: Создать unit tests для валидаторов

**Файл:** `backend/internal/validator/validator_test.go`

**Описание:** Создать тесты для всех валидаторов.

**Тесты:**
- `ValidateLanguage` - валидные и невалидные языки
- `ValidateTopic` - валидные и невалидные темы для каждого языка
- `ValidateDifficulty` - валидные и невалидные уровни сложности (enum значения)
- `ValidateQuestionType` - валидные и невалидные типы ответов (enum значения)
- `ValidateGenerateRequest` - полная валидация запроса со всеми комбинациями ошибок (включая enum валидацию)

**Структура тестов:**
- Table-driven tests для каждого валидатора
- Проверка валидных значений
- Проверка невалидных значений
- Проверка граничных случаев

---

## Этап 4: Клиент OpenRouter

## Этап 6: Документация (опционально)

### Задача 6.1: Обновить README

**Файл:** `backend/README.md` или `README.md`

**Описание:** Добавить документацию по использованию сервиса генерации вопросов.

**Содержимое:**
- Описание сервиса `QuestionService`
- Примеры использования
- Описание параметров запроса
- Список поддерживаемых языков и тем
- Уровни сложности
- Типы ответов
- Примеры работы с сервисом

---

## Чеклист реализации

- [ ] Этап 1: Подготовка проекта и конфигурация
  - [ ] Задача 1.1: Создать `backend/.env.example`
  - [ ] Задача 1.2: Добавить зависимость `go-openrouter`
  - [ ] Задача 1.3: Создать структуру конфигурации
- [ ] Этап 2: Модели данных
  - [ ] Задача 2.1: Создать модель `Question`
  - [ ] Задача 2.2: Создать DTO для запроса
  - [ ] Задача 2.3: Создать DTO для ответа LLM
  - [ ] Задача 2.4: Создать unit tests для моделей
- [ ] Этап 3: Константы и валидация
  - [ ] Задача 3.1: Определить константы языков
  - [ ] Задача 3.2: Определить константы тем
  - [ ] Задача 3.3: Определить константы уровней сложности
  - [ ] Задача 3.4: Определить константы типов ответов
  - [ ] Задача 3.5: Создать валидатор
  - [ ] Задача 3.6: Создать unit tests для валидаторов
- [ ] Этап 4: Клиент OpenRouter
  - [ ] Задача 4.1: Создать интерфейс клиента
  - [ ] Задача 4.2: Реализовать клиент
  - [ ] Задача 4.3: Создать конструктор промпта
  - [ ] Задача 4.4: Создать unit tests для клиента
  - [ ] Задача 4.5: Создать unit tests для конструктора промпта
- [ ] Этап 5: Сервис генерации вопросов
  - [ ] Задача 5.1: Создать интерфейс сервиса
  - [ ] Задача 5.2: Реализовать сервис
  - [ ] Задача 5.3: Создать парсер ответа LLM
  - [ ] Задача 5.4: Создать unit tests для парсера
  - [ ] Задача 5.5: Создать integration tests для сервиса
- [ ] Этап 6: Документация
  - [ ] Задача 6.1: Обновить README

---

## Примечания

1. **Порядок выполнения:** Задачи должны выполняться последовательно по этапам, так как каждый этап зависит от предыдущего.

2. **Обработка ошибок:** Все ошибки должны обрабатываться с добавлением контекста через `fmt.Errorf("...: %w", err)`.

3. **Валидация:** Валидация должна выполняться на нескольких уровнях:
   - На уровне сервиса (валидация бизнес-логики через валидатор)
   - На уровне модели (валидация данных через методы Validate)

4. **Тестирование:** Для интеграционных тестов может потребоваться mock клиента OpenRouter, чтобы не тратить API ключи при каждом тесте.

5. **Конфигурация:** Все настройки должны быть конфигурируемыми через переменные окружения с разумными значениями по умолчанию.

6. **Логирование:** Использовать структурированное логирование для всех важных операций (запросы к API, ошибки, генерация вопросов).

