-- Drop indexes
DROP INDEX IF EXISTS idx_user_questions_answered_at;
DROP INDEX IF EXISTS idx_user_questions_question_id;
DROP INDEX IF EXISTS idx_user_questions_user_id;
DROP INDEX IF EXISTS idx_user_questions_user_question;

-- Drop user_questions table
DROP TABLE IF EXISTS user_questions;

