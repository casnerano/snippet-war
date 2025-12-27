-- Add likes_count column to questions table
ALTER TABLE questions 
    ADD COLUMN likes_count INTEGER NOT NULL DEFAULT 0;

-- Add column comment
COMMENT ON COLUMN questions.likes_count IS 'Number of likes the question has received';

-- Create index for sorting by popularity
CREATE INDEX idx_questions_likes_count ON questions(likes_count);

