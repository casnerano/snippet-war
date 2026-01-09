"""ProxyAPI client for LLM question generation."""

import json
from typing import TypeVar

from app.config import ProxyAPIConfig
from app.models.llm_response import LLMQuestionResponse, LLMQuestionsResponse
from openai import APIError, APITimeoutError, AsyncOpenAI, RateLimitError
from pydantic import BaseModel

T = TypeVar("T", bound=BaseModel)


class ProxyAPIClient:
    """Client for ProxyAPI (OpenAI-compatible API)."""

    def __init__(self, config: ProxyAPIConfig) -> None:
        """Initialize ProxyAPI client."""
        self.config = config
        self.client = AsyncOpenAI(
            api_key=config.api_key,
            base_url=config.base_url,
            timeout=config.timeout,
        )
        self.model = config.model
        self.max_tokens = config.max_tokens

    async def _generate_response(
        self, prompt: str, response_model: type[T], error_context: str
    ) -> T:
        """
        Generate response from ProxyAPI and validate it.

        Args:
            prompt: Prompt text for LLM
            response_model: Pydantic model class for response validation
            error_context: Context string for error messages

        Returns:
            Validated response model instance

        Raises:
            ValueError: If API request fails or response cannot be parsed
        """
        try:
            response = await self.client.chat.completions.create(
                model=self.model,
                messages=[{"role": "user", "content": prompt}],
                response_format={"type": "json_object"},
                max_tokens=self.max_tokens,
            )

            if not response.choices:
                raise ValueError("empty response from API: no choices returned")

            content = response.choices[0].message.content
            if not content:
                raise ValueError("empty response from API: no text content in response")

            try:
                json_data = json.loads(content)
                return response_model.model_validate(json_data)
            except json.JSONDecodeError as e:
                raise ValueError(f"failed to parse JSON response: {e}") from e
            except Exception as e:
                raise ValueError(f"failed to validate LLM response: {e}") from e

        except RateLimitError as e:
            raise ValueError(f"rate limit exceeded: {e}") from e
        except APITimeoutError as e:
            raise ValueError(
                f"request timeout after {self.config.timeout}s: {e}"
            ) from e
        except APIError as e:
            status_code = getattr(e, "status_code", None)
            if status_code == 401:
                raise ValueError(f"authentication failed: {e}") from e
            elif status_code in (500, 502, 503, 504):
                raise ValueError(f"server error (status {status_code}): {e}") from e
            else:
                raise ValueError(f"API error (status {status_code}): {e}") from e
        except Exception as e:
            raise ValueError(f"failed to {error_context}: {e}") from e

    async def generate_question(self, prompt: str) -> LLMQuestionResponse:
        """
        Generate question response from ProxyAPI.

        Args:
            prompt: Prompt text for LLM

        Returns:
            LLMQuestionResponse with generated question data

        Raises:
            ValueError: If API request fails or response cannot be parsed
        """
        return await self._generate_response(
            prompt, LLMQuestionResponse, "generate question"
        )

    async def generate_questions(self, prompt: str) -> LLMQuestionsResponse:
        """
        Generate multiple questions response from ProxyAPI.

        Args:
            prompt: Prompt text for LLM

        Returns:
            LLMQuestionsResponse with generated questions data

        Raises:
            ValueError: If API request fails or response cannot be parsed
        """
        return await self._generate_response(
            prompt, LLMQuestionsResponse, "generate questions"
        )

    async def close(self) -> None:
        """Close the underlying AsyncOpenAI client."""
        await self.client.close()
