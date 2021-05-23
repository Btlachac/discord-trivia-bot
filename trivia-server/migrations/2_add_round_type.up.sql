CREATE TABLE dt.round_type(
	id BIGSERIAL NOT NULL primary key,
	name VARCHAR(50)
);

INSERT INTO dt.round_type(name)
VALUES ('regular'),
	('lightning'),
	('list');

ALTER TABLE dt.round
ADD COLUMN round_type_id BIGINT;

ALTER TABLE dt.round
ADD CONSTRAINT round_round_type_fkey foreign key(round_type_id) REFERENCES dt.round_type(id);

UPDATE dt.round
SET round_type_id = (SELECT id
					 FROM dt.round_type
					 WHERE name = 'regular');

ALTER TABLE dt.round
DROP CONSTRAINT round_round_number_limit_chk;

ALTER TABLE dt.round
ADD CONSTRAINT round_round_number_limit_chk CHECK(round_number > 0 AND round_number < 8);