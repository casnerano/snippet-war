"""ProxyAPI client for LLM question generation."""

import json

from app.config import ProxyAPIConfig
from app.models.llm_response import LLMQuestionResponse
from openai import APIError, APITimeoutError, AsyncOpenAI, RateLimitError


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

    async def generate_question(self, prompt: str) -> LLMQuestionResponse:
        """
        Generate question response from ProxyAPI.

        Args:
            prompt: Prompt text for LLM

        Returns:
            LLMQuestionResponse with generated question data

        Raises:
            APIError: If API request fails
            RateLimitError: If rate limit is exceeded
            APITimeoutError: If request times out
            ValueError: If response cannot be parsed
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

            # Parse JSON response into LLMQuestionResponse
            try:
                json_data = json.loads(content)
                return LLMQuestionResponse.model_validate(json_data)
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
            raise ValueError(f"failed to generate question: {e}") from e

    async def close(self) -> None:
        """Close the underlying AsyncOpenAI client."""
        await self.client.close()
