# Roadmap: Интеграция Backend с LLM через OpenRouter

## Цель

Реализовать backend-часть для генерации вопросов с помощью LLM через OpenRouter API. Roadmap включает конфигурацию, формализацию структуры данных и промпты для генерации.

## Компоненты Roadmap

### 1. Конфигурация OpenRouter

**Файлы:**

- `backend/.env.example` - пример конфигурации
- `backend/internal/config/config.go` - загрузка конфигурации (будущий файл)

**Переменные окружения:**

- `OPENROUTER_API_KEY` - API ключ OpenRouter
- `OPENROUTER_MODEL` - модель LLM (например, `deepseek/deepseek-chat` или `openai/gpt-4o-mini`)
- `OPENROUTER_BASE_URL` - базовый URL API (опционально, по умолчанию `https://openrouter.ai/api/v1`)
- `OPENROUTER_TIMEOUT` - таймаут запросов (опционально, по умолчанию `30s`)
- `OPENROUTER_MAX_TOKENS` - максимальное количество токенов в ответе (опционально, по умолчанию `2000`)

**Библиотека:**Использовать `github.com/revrost/go-openrouter` (стабильная версия v1.1.5) для работы с OpenRouter API.

### 2. Список языков программирования

Поддерживаемые языки (7 самых популярных):

1. **Python** (`python`)
2. **JavaScript** (`javascript`)
3. **Go** (`go`)
4. **Java** (`java`)
5. **C++** (`cpp`)
6. **Rust** (`rust`)
7. **TypeScript** (`typescript`)

### 3. Темы для каждого языка

#### Python

- `variables_types` - Переменные и типы данных
- `lists_arrays` - Списки и массивы
- `dictionaries` - Словари
- `functions` - Функции
- `closures` - Замыкания
- `decorators` - Декораторы
- `generators` - Генераторы
- `classes_oop` - Классы и ООП
- `exceptions` - Обработка исключений
- `context_managers` - Контекстные менеджеры
- `async_await` - Асинхронное программирование

#### JavaScript

- `variables_types` - Переменные и типы
- `arrays` - Массивы
- `objects` - Объекты
- `functions` - Функции
- `closures` - Замыкания
- `this_binding` - Контекст выполнения (this)
- `prototypes` - Прототипы
- `classes` - Классы (ES6+)
- `promises_async` - Промисы и async/await
- `event_loop` - Event Loop
- `destructuring` - Деструктуризация
- `modules` - Модули (ES6+)

#### Go

- `variables_types` - Переменные и типы
- `slices` - Срезы
- `maps` - Мапы
- `functions` - Функции
- `methods` - Методы
- `interfaces` - Интерфейсы
- `goroutines` - Горутины
- `channels` - Каналы
- `select` - Select statement
- `defer_panic_recover` - Defer, panic, recover
- `pointers` - Указатели
- `structs` - Структуры

#### Java

- `variables_types` - Переменные и типы
- `arrays_lists` - Массивы и списки
- `collections` - Коллекции (Set, Map)
- `methods` - Методы
- `classes_objects` - Классы и объекты
- `inheritance` - Наследование
- `interfaces` - Интерфейсы
- `generics` - Дженерики
- `exceptions` - Исключения
- `streams` - Streams API
- `lambda_expressions` - Lambda выражения
- `concurrency` - Многопоточность

#### C++

- `variables_types` - Переменные и типы
- `pointers_references` - Указатели и ссылки
- `arrays_vectors` - Массивы и векторы
- `functions` - Функции
- `classes_objects` - Классы и объекты
- `inheritance` - Наследование
- `templates` - Шаблоны
- `smart_pointers` - Умные указатели
- `stl` - STL контейнеры и алгоритмы
- `move_semantics` - Move семантика
- `lambda` - Lambda выражения
- `multithreading` - Многопоточность

#### Rust

- `variables_types` - Переменные и типы
- `ownership` - Владение (ownership)
- `borrowing` - Заимствование (borrowing)
- `lifetimes` - Время жизни
- `vectors` - Векторы
- `hashmaps` - HashMap
- `functions` - Функции
- `structs` - Структуры
- `enums` - Перечисления
- `pattern_matching` - Сопоставление с образцом
- `error_handling` - Обработка ошибок (Result, Option)
- `concurrency` - Многопоточность

#### TypeScript

- `types` - Типы
- `interfaces` - Интерфейсы
- `generics` - Дженерики
- `unions_intersections` - Объединения и пересечения типов
- `type_guards` - Защитники типов
- `decorators` - Декораторы
- `utility_types` - Утилитарные типы
- `modules` - Модули
- `async_promises` - Асинхронность
- `classes` - Классы
- `namespaces` - Пространства имен

### 4. Уровни сложности

#### Beginner (Начальный)

- Базовые операции и синтаксис
- Простые типы данных
- Базовые структуры данных (массивы, списки)
- Простые условия (if/else)
- Простые циклы (for, while)
- Простые функции без сложной логики
- Базовые операции со строками и числами

**Примеры тем:**

- Переменные и присваивание
- Арифметические операции
- Простые условия
- Циклы по элементам коллекции

#### Intermediate (Средний)

- Более сложные структуры данных
- Вложенные циклы и условия
- Функции высшего порядка (map, filter, reduce)
- Работа с коллекциями (словари, множества)
- Базовое ООП (классы, методы)
- Обработка исключений
- Базовые паттерны проектирования

**Примеры тем:**

- Замыкания и лексическое окружение
- Работа с вложенными структурами
- Функции как объекты первого класса
- Базовые алгоритмы (поиск, сортировка)

#### Advanced (Продвинутый)

- Сложные алгоритмы и оптимизация
- Продвинутые концепции языка (декораторы, метаклассы, макросы)
- Конкурентное/параллельное программирование
- Продвинутые паттерны проектирования
- Неочевидное поведение языка (edge cases)
- Оптимизация производительности
- Работа с памятью и указателями (где применимо)

**Примеры тем:**

- Генераторы и итераторы
- Асинхронное программирование
- Многопоточность и конкурентность
- Продвинутые типы и дженерики
- Неочевидные особенности языка

### 5. Типы ответов

#### Multiple Choice (Выбор из вариантов)

- Тип: `multiple_choice`
- Описание: Пользователь выбирает один правильный ответ из предложенных вариантов
- Структура ответа:
- `question_type`: `"multiple_choice"`
- `options`: массив строк с вариантами ответов (минимум 2, максимум 5)
- `correct_answer`: индекс правильного ответа (0-based)
- `explanation`: объяснение правильного ответа

**Пример:**

```json
{
  "question_type": "multiple_choice",
  "options": ["10", "20", "30", "40"],
  "correct_answer": 1,
  "explanation": "Правильный ответ - 20, потому что..."
}
```



#### Free Text (Свободный ввод)

- Тип: `free_text`
- Описание: Пользователь вводит ответ в текстовом поле
- Структура ответа:
- `question_type`: `"free_text"`
- `correct_answer`: строка с правильным ответом
- `acceptable_variants`: массив альтернативных правильных ответов (опционально)
- `explanation`: объяснение правильного ответа
- `case_sensitive`: чувствительность к регистру (по умолчанию `false`)

**Пример:**

```json
{
  "question_type": "free_text",
  "correct_answer": "Hello World",
  "acceptable_variants": ["hello world", "HELLO WORLD"],
  "case_sensitive": false,
  "explanation": "Правильный ответ - 'Hello World'..."
}
```



### 6. Промпт для генерации вопроса

**Структура промпта:**

```javascript
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

**Параметры для подстановки:**

- `{language}` - название языка (Python, JavaScript, Go, и т.д.)
- `{topic}` - тема из списка тем для языка
- `{difficulty}` - уровень сложности (beginner, intermediate, advanced)
- `{difficulty_description}` - описание уровня сложности
- `{answer_type}` - тип ответа (multiple_choice или free_text)

### 7. Структура данных для хранения

**Модель Question:**

```go
type Question struct {
    ID                string   `json:"id"`
    Language          string   `json:"language"`
    Topic             string   `json:"topic"`
    Difficulty        string   `json:"difficulty"` // beginner, intermediate, advanced
    QuestionType      string   `json:"question_type"` // multiple_choice, free_text
    Code              string   `json:"code"`
    QuestionText      string   `json:"question"`
    Options           []string `json:"options,omitempty"` // для multiple_choice
    CorrectAnswer     string   `json:"correct_answer"` // строка или индекс
    AcceptableVariants []string `json:"acceptable_variants,omitempty"` // для free_text
    CaseSensitive     bool     `json:"case_sensitive,omitempty"` // для free_text
    Explanation       string   `json:"explanation"`
    CreatedAt         time.Time `json:"created_at"`
}
```



## Файлы для создания

1. **`todo/llm.md`** - основной файл roadmap со всей информацией