ALTER TABLE dt.round
DROP CONSTRAINT round_round_type_fkey;

ALTER TABLE dt.round
DROP COLUMN round_type_id;

DROP TABLE dt.round_type;

ALTER TABLE dt.round
DROP CONSTRAINT round_round_number_limit_chk;

ALTER TABLE dt.round
ADD CONSTRAINT round_round_number_limit_chk CHECK(round_number > 0 AND round_number < 7);