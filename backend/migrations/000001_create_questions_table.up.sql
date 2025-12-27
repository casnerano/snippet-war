-- Create questions table
CREATE TABLE questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    language VARCHAR(50) NOT NULL,
    topic VARCHAR(255) NOT NULL,
    difficulty VARCHAR(20) NOT NULL,
    question_type VARCHAR(20) NOT NULL,
    code TEXT NOT NULL,
    question_text TEXT NOT NULL,
    options JSONB,
    correct_answers JSONB NOT NULL,
    explanation TEXT NOT NULL,
    likes_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

-- Add table comment
COMMENT ON TABLE questions IS 'Stores programming questions for the Snippet War game';

-- Add column comments
COMMENT ON COLUMN questions.id IS 'Unique identifier for the question';
COMMENT ON COLUMN questions.language IS 'Programming language (python, javascript, go, java, cpp, rust, typescript)';
COMMENT ON COLUMN questions.topic IS 'Topic of the question (e.g., functions, variables)';
COMMENT ON COLUMN questions.difficulty IS 'Difficulty level (beginner, intermediate, advanced)';
COMMENT ON COLUMN questions.question_type IS 'Type of question (multiple_choice, free_text)';
COMMENT ON COLUMN questions.code IS 'Code snippet for the question';
COMMENT ON COLUMN questions.question_text IS 'The question text';
COMMENT ON COLUMN questions.options IS 'Array of options for multiple choice questions (JSONB)';
COMMENT ON COLUMN questions.correct_answers IS 'Array of correct answers (JSONB). For multiple_choice: array of exact texts from options. For free_text: array of all acceptable variants.';
COMMENT ON COLUMN questions.explanation IS 'Explanation of the correct answer';
COMMENT ON COLUMN questions.likes_count IS 'Number of likes the question has received';
COMMENT ON COLUMN questions.created_at IS 'Timestamp when the question was created';

-- Create indexes for filtering and sorting
CREATE INDEX idx_questions_language ON questions(language);
CREATE INDEX idx_questions_difficulty ON questions(difficulty);
CREATE INDEX idx_questions_topic ON questions(topic);
CREATE INDEX idx_questions_created_at ON questions(created_at);
CREATE INDEX idx_questions_likes_count ON questions(likes_count);

-- Add CHECK constraints for enum-like values
ALTER TABLE questions ADD CONSTRAINT chk_questions_language 
    CHECK (language IN ('python', 'javascript', 'go', 'java', 'cpp', 'rust', 'typescript'));

ALTER TABLE questions ADD CONSTRAINT chk_questions_difficulty 
    CHECK (difficulty IN ('beginner', 'intermediate', 'advanced'));

ALTER TABLE questions ADD CONSTRAINT chk_questions_question_type 
    CHECK (question_type IN ('multiple_choice', 'free_text'));

