-- Drop questions table and its indexes
DROP INDEX IF EXISTS idx_questions_likes_count;
DROP INDEX IF EXISTS idx_questions_created_at;
DROP INDEX IF EXISTS idx_questions_topic;
DROP INDEX IF EXISTS idx_questions_difficulty;
DROP INDEX IF EXISTS idx_questions_language;

DROP TABLE IF EXISTS questions;

