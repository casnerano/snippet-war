# Platform Service

Platform service for generating questions for Snippet War game using LLM.

## Description

This service generates programming questions using Large Language Models (LLM). It supports multiple LLM providers (ProxyAPI, OpenRouter) and uses structured output (JSON mode) for reliable parsing.

## Technology Stack

- **FastAPI** - Modern async web framework
- **Pydantic** - Data validation using Python type annotations
- **OpenAI SDK** - For LLM API interactions with structured output
- **Loguru** - Structured logging
- **Uvicorn** - ASGI server

## Setup

### Prerequisites

- Python 3.13+
- uv (Python package manager)

### Installation

1. Install dependencies:
```bash
cd backend/services/platform-service
uv sync
```

2. Set environment variables:
```bash
# Required
export PROXYAPI_API_KEY=your_api_key
export PROXYAPI_MODEL=gpt-4.1-mini

# Optional
export PROXYAPI_BASE_URL=https://api.proxyapi.ru/openai/v1
export PROXYAPI_TIMEOUT=30
export PROXYAPI_MAX_TOKENS=2000
```

Or create a `.env` file:
```
PROXYAPI_API_KEY=your_api_key
PROXYAPI_MODEL=gpt-4.1-mini
PROXYAPI_BASE_URL=https://api.proxyapi.ru/openai/v1
PROXYAPI_TIMEOUT=30
PROXYAPI_MAX_TOKENS=2000
```

### Running

#### Development

```bash
uv run uvicorn app.main:app --reload --host 0.0.0.0 --port 8081
```

#### Production

```bash
uv run uvicorn app.main:app --host 0.0.0.0 --port 8081
```

#### Docker

```bash
docker build -t platform-service .
docker run -p 8081:8081 --env-file .env platform-service
```

## API Endpoints

### POST /api/questions/generate

Generate a question based on request parameters.

**Request Body:**
```json
{
  "language": "python",
  "topic": "functions",
  "difficulty": "beginner",
  "question_type": "multiple_choice"
}
```

**Response:**
```json
{
  "id": "uuid",
  "language": "python",
  "topic": "functions",
  "difficulty": "beginner",
  "question_type": "multiple_choice",
  "code": "def add(a, b):\n    return a + b\n\nprint(add(2, 3))",
  "question_text": "Что выведет этот код?",
  "options": ["5", "6", "7", "8"],
  "correct_answers": ["5"],
  "explanation": "...",
  "created_at": "2024-01-01T00:00:00Z"
}
```

### GET /health

Health check endpoint.

**Response:**
```json
{
  "status": "ok"
}
```

## Supported Languages

- Python
- JavaScript
- Go
- Java
- C++
- Rust
- TypeScript

## Supported Question Types

- `multiple_choice` - Multiple choice questions with 2-5 options
- `free_text` - Free text answer questions

## Supported Difficulty Levels

- `beginner` - Basic operations and syntax
- `intermediate` - More complex structures and patterns
- `advanced` - Advanced concepts and optimizations

## Project Structure

```
app/
├── main.py              # FastAPI application entry point
├── config.py            # Configuration management
├── exceptions.py        # Custom exceptions
├── models/              # Pydantic models
│   ├── enums.py        # Enumerations (Language, Difficulty, etc.)
│   ├── question.py     # Question models
│   ├── llm_response.py # LLM response models
│   └── validation.py   # Validation functions
├── services/            # Business logic
│   └── question_service.py
├── clients/             # LLM clients
│   ├── llm_client.py   # LLM client interface
│   ├── proxyapi_client.py
│   └── prompt_builder.py
└── routers/            # API routes
    └── questions.py
```

## Development

### Code Formatting

```bash
uv run ruff check --fix .
uv run ruff format .
```

### Testing

```bash
uv run pytest
```

## Migration from Go

This service was migrated from Go to Python. Key improvements:

1. **Structured Output**: Uses OpenAI SDK JSON mode instead of manual parsing
2. **Async/Await**: Native async support for I/O operations
3. **Pydantic Validation**: Automatic validation at model level
4. **Type Safety**: Full type hints for better IDE support

