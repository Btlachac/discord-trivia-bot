package db

import (
	"context"
	"database/sql"
)

type TriviaRepository struct {
	db *sql.DB
}

func NewTriviaRepository(db *sql.DB) *TriviaRepository {
	return &TriviaRepository{
		db: db,
	}
}

func (r *TriviaRepository) AddTrivia(ctx context.Context, newTrivia Trivia, audioFileName string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	insertTriviaStatement := `
  INSERT INTO dt.trivia(image_round_theme, image_round_detail, image_round_url, audio_round_theme, answer_url, audio_file_name)
  VALUES($1, $2, $3, $4, $5, $6)`

	result, err := tx.ExecContext(ctx, insertTriviaStatement, newTrivia.ImageRoundTheme, newTrivia.ImageRoundDetail, newTrivia.ImageRoundURL, newTrivia.AudioRoundTheme, newTrivia.AnswersURL, audioFileName)
	if err != nil {
		return err
	}

	triviaID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	for _, round := range newTrivia.Rounds {
		err = addRound(ctx, tx, round, triviaID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func addRound(ctx context.Context, tx *sql.Tx, newRound Round, triviaId int64) error {
	insertRoundStatement := `
  INSERT INTO dt.round(trivia_id, round_number, theme, theme_description, round_type_id)
  VALUES($1, $2, $3, $4, $5)`

	result, err := tx.ExecContext(ctx, insertRoundStatement, triviaId, newRound.RoundNumber, newRound.Theme, newRound.ThemeDescription, newRound.RoundType.Id)
	if err != nil {
		return err
	}

	roundID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	for _, question := range newRound.Questions {
		err = addQuestion(ctx, tx, question, roundID)
		if err != nil {
			return err
		}
	}
	return nil
}

func addQuestion(ctx context.Context, tx *sql.Tx, newQuestion Question, roundId int64) error {
	insertQuestionStatement := `
  INSERT INTO dt.question(round_id, question_number, question)
  VALUES($1, $2, $3)`

	_, err := tx.ExecContext(ctx, insertQuestionStatement, roundId, newQuestion.QuestionNumber, newQuestion.Question)
	if err != nil {
		return err
	}
	return nil
}

func (r *TriviaRepository) MarkTriviaUsed(ctx context.Context, triviaId int64) error {
	updateTriviaStatement := `
	UPDATE dt.trivia
	SET used = true,
		date_used = CURRENT_DATE
	WHERE id = $1`

	_, err := r.db.ExecContext(ctx, updateTriviaStatement, triviaId)
	if err != nil {
		return err
	}
	return nil

}

func (r *TriviaRepository) RoundTypesList(ctx context.Context) ([]RoundType, error) {
	selectStatement := `
	SELECT id, name
	FROM dt.round_type
	`

	var roundTypes []RoundType
	rows, err := r.db.QueryContext(ctx, selectStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var roundType RoundType

		err := rows.Scan(&roundType.Id, &roundType.Name)
		if err != nil {
			return roundTypes, err
		}

		roundTypes = append(roundTypes, roundType)
	}

	return roundTypes, nil
}

func (r *TriviaRepository) GetNewTrivia(ctx context.Context) (Trivia, string, error) {
	selectTriviaStatement := `
  SELECT id, image_round_theme, image_round_detail, image_round_url, audio_round_theme, answer_url, audio_file_name
  FROM dt.trivia
  WHERE used = false
  ORDER BY date_created ASC
  FETCH FIRST ROW ONLY`

	var trivia Trivia
	var audioFileNameHolder sql.NullString
	audioFileName := ""

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Trivia{}, "", err
	}

	err = tx.QueryRowContext(ctx, selectTriviaStatement).Scan(&trivia.Id, &trivia.ImageRoundTheme, &trivia.ImageRoundDetail, &trivia.ImageRoundURL, &trivia.AudioRoundTheme, &trivia.AnswersURL, &audioFileNameHolder)
	if err != nil {
		return Trivia{}, "", err
	}

	if audioFileNameHolder.Valid {
		audioFileName = audioFileNameHolder.String
	}

	trivia.Rounds, err = getRounds(ctx, tx, trivia.Id)
	if err != nil {
		return Trivia{}, "", err
	}

	if err = tx.Commit(); err != nil {
		return Trivia{}, "", err
	}

	return trivia, audioFileName, nil
}

func getRounds(ctx context.Context, tx *sql.Tx, triviaId int64) ([]Round, error) {
	selectRoundsStatement := `
	SELECT r.id, r.round_number, r.theme, r.theme_description, rt.name
	FROM dt.round r JOIN dt.round_type rt ON r.round_type_id = rt.id
	WHERE trivia_id = $1
	ORDER BY r.round_number ASC
	`
	var rounds []Round
	rows, err := tx.QueryContext(ctx, selectRoundsStatement, triviaId)
	if err != nil {
		return rounds, err
	}
	defer rows.Close()

	for rows.Next() {
		var round Round

		err := rows.Scan(&round.Id, &round.RoundNumber, &round.Theme, &round.ThemeDescription, &round.RoundType.Name)
		if err != nil {
			return rounds, err
		}

		round.Questions, err = getQuestions(ctx, tx, round.Id)
		if err != nil {
			return rounds, err
		}

		rounds = append(rounds, round)
	}

	return rounds, nil
}

func getQuestions(ctx context.Context, tx *sql.Tx, roundId int64) ([]Question, error) {
	selectQuestionsStatement := `
  SELECT question_number, question
  FROM dt.question
  WHERE round_id = $1
  ORDER BY question_number ASC
  `
	var questions []Question

	rows, err := tx.QueryContext(ctx, selectQuestionsStatement, roundId)
	if err != nil {
		return questions, err
	}
	defer rows.Close()

	for rows.Next() {
		var question Question

		err := rows.Scan(&question.QuestionNumber, &question.Question)
		if err != nil {
			return questions, err
		}

		questions = append(questions, question)
	}

	return questions, nil
}
