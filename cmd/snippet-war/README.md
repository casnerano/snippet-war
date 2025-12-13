
## 📡 API Документация

### Генерация вопроса

Генерирует новый вопрос на основе указанных параметров.

**Endpoint:** `POST /api/questions/generate`

**Content-Type:** `application/json`

**Request Body:**

```json
{
  "language": "python",
  "topic": "variables_types",
  "difficulty": "beginner",
  "question_type": "multiple_choice"
}
```

**Параметры запроса:**

- `language` (string, обязательно) - Язык программирования. Доступные значения:
  - `python`
  - `javascript`
  - `go`
  - `java`
  - `cpp`
  - `rust`
  - `typescript`

- `topic` (string, обязательно) - Тема вопроса. Зависит от выбранного языка. Примеры для Python:
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

- `difficulty` (string, обязательно) - Уровень сложности. Доступные значения:
  - `beginner` - Начальный
  - `intermediate` - Средний
  - `advanced` - Продвинутый

- `question_type` (string, обязательно) - Тип вопроса. Доступные значения:
  - `multiple_choice` - Множественный выбор
  - `free_text` - Свободный текст

**Response 200 OK:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "language": "python",
  "topic": "variables_types",
  "difficulty": "beginner",
  "question_type": "multiple_choice",
  "code": "x = 5\ny = 3\nprint(x + y)",
  "question": "Что выведет этот код?",
  "options": [
    "8",
    "53",
    "5 + 3",
    "Ошибка"
  ],
  "correct_answer": "8",
  "explanation": "Код выполняет арифметическую операцию сложения двух чисел...",
  "created_at": "2024-01-15T10:30:00Z"
}
```

**Response 400 Bad Request:**

```json
{
  "error": "invalid request: language is required"
}
```

**Response 500 Internal Server Error:**

```json
{
  "error": "failed to generate question from API: authentication failed"
}
```

### Примеры использования

#### Пример 1: Python, начальный уровень, множественный выбор

```bash
curl -X POST http://localhost:8081/api/questions/generate \
  -H "Content-Type: application/json" \
  -d '{
    "language": "python",
    "topic": "variables_types",
    "difficulty": "beginner",
    "question_type": "multiple_choice"
  }'
```

#### Пример 2: JavaScript, средний уровень, свободный текст

```bash
curl -X POST http://localhost:8081/api/questions/generate \
  -H "Content-Type: application/json" \
  -d '{
    "language": "javascript",
    "topic": "closures",
    "difficulty": "intermediate",
    "question_type": "free_text"
  }'
```

#### Пример 3: Go, продвинутый уровень, множественный выбор

```bash
curl -X POST http://localhost:8081/api/questions/generate \
  -H "Content-Type: application/json" \
  -d '{
    "language": "go",
    "topic": "goroutines",
    "difficulty": "advanced",
    "question_type": "multiple_choice"
  }'
```

#### Пример 4: Использование с JavaScript (fetch)

```javascript
const response = await fetch('http://localhost:8081/api/questions/generate', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    language: 'python',
    topic: 'functions',
    difficulty: 'intermediate',
    question_type: 'multiple_choice'
  })
});

const question = await response.json();
console.log(question);
```
