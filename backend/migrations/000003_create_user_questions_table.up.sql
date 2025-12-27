-- Create user_questions table
CREATE TABLE user_questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    question_id UUID NOT NULL,
    seen_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    answered_at TIMESTAMP WITH TIME ZONE,
    is_correct BOOLEAN,
    CONSTRAINT fk_user_questions_user_id 
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_questions_question_id 
        FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE
);

-- Add table comment
COMMENT ON TABLE user_questions IS 'Tracks which users saw and answered which questions';

-- Add column comments
COMMENT ON COLUMN user_questions.id IS 'Unique identifier for the user-question interaction';
COMMENT ON COLUMN user_questions.user_id IS 'Reference to the user';
COMMENT ON COLUMN user_questions.question_id IS 'Reference to the question';
COMMENT ON COLUMN user_questions.seen_at IS 'Timestamp when user first saw the question';
COMMENT ON COLUMN user_questions.answered_at IS 'Timestamp when user answered (NULL if not answered yet)';
COMMENT ON COLUMN user_questions.is_correct IS 'Whether the answer was correct (NULL if not answered yet)';

-- Create unique index to prevent duplicate entries
CREATE UNIQUE INDEX idx_user_questions_user_question 
    ON user_questions(user_id, question_id);

-- Create indexes for efficient queries
CREATE INDEX idx_user_questions_user_id ON user_questions(user_id);
CREATE INDEX idx_user_questions_question_id ON user_questions(question_id);
CREATE INDEX idx_user_questions_answered_at ON user_questions(answered_at);

