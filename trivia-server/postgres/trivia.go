package postgres

import (
	"context"
	"database/sql"

	"go.uber.org/zap"
)

type Question struct {
	Question       string `json:"question"`
	QuestionNumber int    `json:"questionNumber"`
}

type Round struct {
	Id               int64       `json:"id"`
	Questions        []*Question `json:"questions"`
	RoundNumber      int         `json:"roundNumber"`
	Theme            string      `json:"theme"`
	ThemeDescription string      `json:"themeDescription"`
	RoundType        *RoundType  `json:"roundType"`
}

type Trivia struct {
	Id               int64    `json:"id"`
	Rounds           []*Round `json:"rounds"`
	AnswersURL       string   `json:"answersURL"`
	AudioBinary      string   `json:"audioBinary"`
	AudioRoundTheme  string   `json:"audioRoundTheme"`
	ImageRoundDetail string   `json:"imageRoundDetail"`
	ImageRoundTheme  string   `json:"imageRoundTheme"`
	ImageRoundURL    string   `json:"imageRoundURL"`
}

type RoundType struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type TriviaRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewTriviaRepository(db *sql.DB, logger *zap.Logger) *TriviaRepository {
	return &TriviaRepository{
		db:     db,
		logger: logger,
	}
}

func (r *TriviaRepository) GetNewTrivia(ctx context.Context) (*Trivia, string, error) {
	var trivia Trivia
	var audioFileNameHolder sql.NullString
	audioFileName := ""

	err := r.db.QueryRow(GetNewTriviaQuery).Scan(&trivia.Id, &trivia.ImageRoundTheme, &trivia.ImageRoundDetail, &trivia.ImageRoundURL, &trivia.AudioRoundTheme, &trivia.AnswersURL, &audioFileNameHolder)
	if err != nil {
		return &trivia, audioFileName, err
	}

	if audioFileNameHolder.Valid {
		audioFileName = audioFileNameHolder.String
	}

	trivia.Rounds, err = r.getRounds(trivia.Id)

	return &trivia, audioFileName, err
}

func (r *TriviaRepository) AddTrivia(ctx context.Context, newTrivia *Trivia, audioFileName string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	err = tx.QueryRow(AddNewTriviaQuery, newTrivia.ImageRoundTheme, newTrivia.ImageRoundDetail, newTrivia.ImageRoundURL, newTrivia.AudioRoundTheme, newTrivia.AnswersURL, audioFileName).Scan(&newTrivia.Id)
	if err != nil {
		return err
	}

	for _, round := range newTrivia.Rounds {
		err = r.addRound(ctx, tx, round, newTrivia.Id)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *TriviaRepository) MarkTriviaUsed(ctx context.Context, triviaId int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.Exec(MarkTriviaUsedQuery, triviaId)

	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil

}

func (r *TriviaRepository) RoundTypesList(ctx context.Context) ([]*RoundType, error) {
	var roundTypes []*RoundType
	rows, err := r.db.QueryContext(ctx, RoundTypesListQuery)
	if err != nil {
		return roundTypes, err
	}
	defer rows.Close()

	for rows.Next() {
		var roundType *RoundType

		err := rows.Scan(roundType.Id, roundType.Name)
		if err != nil {
			return roundTypes, err
		}

		roundTypes = append(roundTypes, roundType)
	}

	return roundTypes, nil
}

func (r *TriviaRepository) getRounds(triviaId int64) ([]*Round, error) {
	var rounds []*Round
	rows, err := r.db.Query(getRoundsQuery, triviaId)
	if err != nil {
		return rounds, err
	}
	defer rows.Close()

	for rows.Next() {
		var round *Round

		err := rows.Scan(&round.Id, &round.RoundNumber, &round.Theme, &round.ThemeDescription, &round.RoundType.Name)
		if err != nil {
			return rounds, err
		}

		round.Questions, err = r.getQuestions(round.Id)
		if err != nil {
			return rounds, err
		}

		rounds = append(rounds, round)
	}

	return rounds, nil
}

func (r *TriviaRepository) getQuestions(roundId int64) ([]*Question, error) {
	var questions []*Question

	rows, err := r.db.Query(getQuestionsQuery, roundId)
	if err != nil {
		return questions, err
	}
	defer rows.Close()

	for rows.Next() {
		var question *Question

		err := rows.Scan(question.QuestionNumber, question.Question)
		if err != nil {
			return questions, err
		}

		questions = append(questions, question)
	}

	return questions, nil
}

func (r *TriviaRepository) addRound(ctx context.Context, tx *sql.Tx, newRound *Round, triviaId int64) error {
	err := tx.QueryRow(addRoundQuery, triviaId, newRound.RoundNumber, newRound.Theme, newRound.ThemeDescription, newRound.RoundType.Id).Scan(&newRound.Id)
	if err != nil {
		return err
	}

	for _, question := range newRound.Questions {
		err = r.addQuestion(ctx, tx, question, newRound.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *TriviaRepository) addQuestion(ctx context.Context, tx *sql.Tx, newQuestion *Question, roundId int64) error {
	_, err := tx.Exec(addQuestionQuery, roundId, newQuestion.QuestionNumber, newQuestion.Question)
	if err != nil {
		return err
	}
	return nil
}
