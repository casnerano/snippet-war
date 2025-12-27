-- Drop index
DROP INDEX IF EXISTS idx_questions_likes_count;

-- Remove likes_count column from questions table
ALTER TABLE questions DROP COLUMN IF EXISTS likes_count;

