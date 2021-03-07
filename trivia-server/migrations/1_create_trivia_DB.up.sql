DROP TABLE IF EXISTS dt.question;
DROP Table IF EXISTS dt.round;
DROP Table IF EXISTS dt.trivia;
DROP SCHEMA IF EXISTS dt;

CREATE SCHEMA dt;

CREATE TABLE dt.trivia(
	id BIGSERIAL NOT NULL PRIMARY KEY,
	date_created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, 
	used BOOLEAN NOT NULL DEFAULT false,
	date_used DATE,
	image_round_theme VARCHAR(400) NOT NULL,
	image_round_detail VARCHAR(400) NOT NULL,
	image_round_url VARCHAR(400) NOT NULL,
	audio_round_theme VARCHAR(400),
	answer_url VARCHAR(400) NOT NULL,
	audio_file_name VARCHAR(100)	
);

CREATE TABLE dt.round(
	id BIGSERIAL NOT NULL primary key,
	round_number INTEGER NOT NULL,
	trivia_id BIGINT NOT NULL,
	theme VARCHAR(400) NOT NULL,
	theme_description VARCHAR(400),
	CONSTRAINT round_trivia_fkey foreign key(trivia_id) REFERENCES dt.trivia(id),
	CONSTRAINT round_round_number_limit_chk CHECK(round_number > 0 AND round_number < 7)
);

CREATE TABLE dt.question(
	round_id BIGINT NOT NULL,
	question_number INT NOT NULL,
	question VARCHAR(1000) NOT NULL,
	CONSTRAINT question_pkey PRIMARY KEY(round_id, question_number),
	CONSTRAINT question_round_fkey FOREIGN KEY(round_id) REFERENCES dt.round(id),
	CONSTRAINT question_question_number_limit_chk CHECK(question_number > 0 AND question_number < 6)
);