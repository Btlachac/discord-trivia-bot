package postgres

import (
	"database/sql"
	"go-trivia-api/model"
)

type TriviaRepository struct {
	db *sql.DB
}

func NewTriviaRepository(db *sql.DB) *TriviaRepository {
	return &TriviaRepository{
		db: db,
	}
}

func (repository *TriviaRepository) GetNewTrivia() (model.Trivia, string, error) {
	selectTriviaStatement := `
  SELECT id, image_round_theme, image_round_detail, image_round_url, audio_round_theme, answer_url, audio_file_name
  FROM dt.trivia
  WHERE used = false
  ORDER BY date_created ASC
  FETCH FIRST ROW ONLY`

	var trivia model.Trivia
	var audioFileNameHolder sql.NullString
	audioFileName := ""

	err := repository.db.QueryRow(selectTriviaStatement).Scan(&trivia.Id, &trivia.ImageRoundTheme, &trivia.ImageRoundDetail, &trivia.ImageRoundURL, &trivia.AudioRoundTheme, &trivia.AnswersURL, &audioFileNameHolder)
	if err != nil {
		return trivia, audioFileName, &model.QueryError{Query: selectTriviaStatement, Err: err}
	}

	if audioFileNameHolder.Valid {
		audioFileName = audioFileNameHolder.String
	}

	trivia.Rounds, err = repository.getRounds(trivia.Id)

	return trivia, audioFileName, err
}

func (repository *TriviaRepository) AddTrivia(newTrivia model.Trivia, audioFileName string) error {

	insertTriviaStatement := `
  INSERT INTO dt.trivia(image_round_theme, image_round_detail, image_round_url, audio_round_theme, answer_url, audio_file_name)
  VALUES($1, $2, $3, $4, $5, $6)
  RETURNING id`

	err := repository.db.QueryRow(insertTriviaStatement, newTrivia.ImageRoundTheme, newTrivia.ImageRoundDetail, newTrivia.ImageRoundURL, newTrivia.AudioRoundTheme, newTrivia.AnswersURL, audioFileName).Scan(&newTrivia.Id)
	if err != nil {
		return &model.QueryError{Query: insertTriviaStatement, Err: err}
	}

	for _, round := range newTrivia.Rounds {
		err = repository.addRound(round, newTrivia.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repository *TriviaRepository) MarkTriviaUsed(triviaId int64) error {
	updateTriviaStatement := `
	UPDATE dt.trivia
	SET used = true,
		date_used = CURRENT_DATE
	WHERE id = $1`

	_, err := repository.db.Exec(updateTriviaStatement, triviaId)
	if err != nil {
		return &model.QueryError{Query: updateTriviaStatement, Err: err}
	}
	return nil

}

func (repository *TriviaRepository) getRounds(triviaId int64) ([]model.Round, error) {
	selectRoundsStatement := `
  SELECT id, round_number, theme, theme_description
  FROM dt.round
  WHERE trivia_id = $1
  `
	var rounds []model.Round
	rows, err := repository.db.Query(selectRoundsStatement, triviaId)
	if err != nil {
		return rounds, &model.QueryError{Query: selectRoundsStatement, Err: err}
	}
	defer rows.Close()

	for rows.Next() {
		var round model.Round

		err := rows.Scan(&round.Id, &round.RoundNumber, &round.Theme, &round.ThemeDescription)
		if err != nil {
			return rounds, &model.QueryError{Query: selectRoundsStatement, Err: err}
		}

		round.Questions, err = repository.getQuestions(round.Id)
		if err != nil {
			return rounds, err
		}

		rounds = append(rounds, round)
	}

	return rounds, nil
}

func (repository *TriviaRepository) getQuestions(roundId int64) ([]model.Question, error) {
	selectQuestionsStatement := `
  SELECT question_number, question
  FROM dt.question
  WHERE round_id = $1
  `
	var questions []model.Question

	rows, err := repository.db.Query(selectQuestionsStatement, roundId)
	if err != nil {
		return questions, &model.QueryError{Query: selectQuestionsStatement, Err: err}
	}
	defer rows.Close()

	for rows.Next() {
		var question model.Question

		err := rows.Scan(&question.QuestionNumber, &question.Question)
		if err != nil {
			return questions, &model.QueryError{Query: selectQuestionsStatement, Err: err}
		}

		questions = append(questions, question)
	}

	return questions, nil
}

func (repository *TriviaRepository) addRound(newRound model.Round, triviaId int64) error {
	insertRoundStatement := `
  INSERT INTO dt.round(trivia_id, round_number, theme, theme_description)
  VALUES($1, $2, $3, $4)
  RETURNING id`

	err := repository.db.QueryRow(insertRoundStatement, triviaId, newRound.RoundNumber, newRound.Theme, newRound.ThemeDescription).Scan(&newRound.Id)
	if err != nil {
		return &model.QueryError{Query: insertRoundStatement, Err: err}
	}

	for _, question := range newRound.Questions {
		err = repository.addQuestion(question, newRound.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repository *TriviaRepository) addQuestion(newQuestion model.Question, roundId int64) error {
	insertQuestionStatement := `
  INSERT INTO dt.question(round_id, question_number, question)
  VALUES($1, $2, $3)`

	_, err := repository.db.Exec(insertQuestionStatement, roundId, newQuestion.QuestionNumber, newQuestion.Question)
	if err != nil {
		return &model.QueryError{Query: insertQuestionStatement, Err: err}
	}
	return nil
}
