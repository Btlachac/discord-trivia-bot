ALTER TABLE dt.question
DROP CONSTRAINT question_question_number_limit_chk;

ALTER TABLE dt.question
ADD CONSTRAINT question_question_number_limit_chk CHECK(question_number > 0 AND question_number < 6);