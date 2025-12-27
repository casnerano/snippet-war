"""Test script to generate a random question from the content service."""

import json
import random
import sys
from pathlib import Path

import httpx

# Add parent directory to path to import app models
sys.path.insert(0, str(Path(__file__).parent.parent))

from app.models.enums import Difficulty, Language, QuestionType
from app.models.topics import get_topics_for_language


def get_random_request() -> dict[str, str]:
    """Generate a random question generation request."""
    # Random language
    language = random.choice(list(Language))

    # Random topic for the selected language
    topics = get_topics_for_language(language)
    topic = random.choice(topics)

    # Random difficulty
    difficulty = random.choice(list(Difficulty))

    # Random question type
    question_type = random.choice(list(QuestionType))

    return {
        "language": language.value,
        "topic": topic.topic_id,
        "difficulty": difficulty.value,
        "question_type": question_type.value,
    }


def print_question(question: dict) -> None:
    """Print question in a formatted way."""
    print("\n" + "=" * 80)
    print("GENERATED QUESTION")
    print("=" * 80)
    print(f"\nID: {question.get('id')}")
    print(f"Language: {question.get('language')}")
    print(f"Topic: {question.get('topic')}")
    print(f"Difficulty: {question.get('difficulty')}")
    print(f"Question Type: {question.get('question_type')}")
    print("\nCode:")
    print("-" * 80)
    print(question.get("code", ""))
    print("-" * 80)
    print(f"\nQuestion: {question.get('question_text', '')}")

    if question.get("options"):
        print("\nOptions:")
        correct_answers = question.get("correct_answers", [])
        if not isinstance(correct_answers, list):
            correct_answers = [correct_answers]
        for i, option in enumerate(question["options"], 1):
            marker = "✓" if option in correct_answers else " "
            print(f"  {marker} {i}. {option}")

    correct_answers = question.get("correct_answers", [])
    if isinstance(correct_answers, list):
        print(f"\nCorrect Answer(s): {', '.join(correct_answers)}")
    else:
        print(f"\nCorrect Answer: {correct_answers}")

    print("\nExplanation:")
    print("-" * 80)
    print(question.get("explanation", ""))
    print("-" * 80)
    print(f"\nCreated At: {question.get('created_at', '')}")
    print("=" * 80)
    print()


def main() -> None:
    """Main function to test question generation."""
    # Default URL, can be overridden with command line argument
    base_url = sys.argv[1] if len(sys.argv) > 1 else "http://localhost:8081"
    endpoint = f"{base_url}/api/questions/generate"

    print(f"Connecting to: {endpoint}")

    # Generate random request
    request_data = get_random_request()
    print("\nRequest parameters:")
    print(json.dumps(request_data, indent=2))

    # Send request
    try:
        with httpx.Client(timeout=60.0) as client:
            response = client.post(endpoint, json=request_data)

            if response.status_code == 200:
                question = response.json()
                print_question(question)
            else:
                print(f"\nError: {response.status_code}")
                print(f"Response: {response.text}")
                sys.exit(1)

    except httpx.ConnectError:
        print(f"\nError: Could not connect to {base_url}")
        print("Make sure the server is running.")
        sys.exit(1)
    except httpx.TimeoutException:
        print("\nError: Request timed out")
        print("The question generation might be taking too long.")
        sys.exit(1)
    except Exception as e:
        print(f"\nError: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()
