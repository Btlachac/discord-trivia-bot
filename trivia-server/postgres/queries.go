package postgres

const GetNewTriviaQuery = `
SELECT id, image_round_theme, image_round_detail, image_round_url, audio_round_theme, answer_url, audio_file_name
FROM dt.trivia
WHERE used = false
ORDER BY date_created ASC
FETCH FIRST ROW ONLY`

const AddNewTriviaQuery = `
INSERT INTO dt.trivia(image_round_theme, image_round_detail, image_round_url, audio_round_theme, answer_url, audio_file_name)
VALUES($1, $2, $3, $4, $5, $6)
RETURNING id`

const MarkTriviaUsedQuery = `
INSERT INTO dt.trivia(image_round_theme, image_round_detail, image_round_url, audio_round_theme, answer_url, audio_file_name)
VALUES($1, $2, $3, $4, $5, $6)
RETURNING id`

const RoundTypesListQuery = `
SELECT id, name
FROM dt.round_type
`

const getRoundsQuery = `
SELECT r.id, r.round_number, r.theme, r.theme_description, rt.name
FROM dt.round r JOIN dt.round_type rt ON r.round_type_id = rt.id
WHERE trivia_id = $1
`

const getQuestionsQuery = `
SELECT question_number, question
FROM dt.question
WHERE round_id = $1
ORDER BY question_number ASC
`

const addRoundQuery = `
INSERT INTO dt.round(trivia_id, round_number, theme, theme_description, round_type_id)
VALUES($1, $2, $3, $4, $5)
RETURNING id`

const addQuestionQuery = `
INSERT INTO dt.question(round_id, question_number, question)
VALUES($1, $2, $3)`
