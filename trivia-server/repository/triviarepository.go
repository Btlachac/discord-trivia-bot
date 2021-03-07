package repository

import (
	"database/sql"
	b64 "encoding/base64"
	"fmt"
	"go-trivia-api/model"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
)

type TriviaRepository struct {
	db *sql.DB
}

func NewTriviaRepository(db *sql.DB) *TriviaRepository {
	return &TriviaRepository{
		db: db,
	}
}

func (repository *TriviaRepository) GetNewTrivia() model.Trivia {
	selectTriviaStatement := `
  SELECT id, image_round_theme, image_round_detail, image_round_url, audio_round_theme, answer_url, audio_file_name
  FROM dt.trivia
  WHERE used = false
  ORDER BY date_created ASC
  FETCH FIRST ROW ONLY`

	var trivia model.Trivia
	var audioFileName sql.NullString

	err := repository.db.QueryRow(selectTriviaStatement).Scan(&trivia.Id, &trivia.ImageRoundTheme, &trivia.ImageRoundDetail, &trivia.ImageRoundURL, &trivia.AudioRoundTheme, &trivia.AnswersURL, &audioFileName)
	if err != nil {
		panic(err)
	}

	if audioFileName.Valid && len(audioFileName.String) > 0 {
		trivia.AudioBinary = getAudioBinary(audioFileName.String)
	}

	trivia.Rounds = repository.getRounds(trivia.Id)

	return trivia
}

func (repository *TriviaRepository) getRounds(triviaId int64) []model.Round {
	selectRoundsStatement := `
  SELECT id, round_number, theme, theme_description
  FROM dt.round
  WHERE trivia_id = $1
  `
	rows, err := repository.db.Query(selectRoundsStatement, triviaId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var rounds []model.Round

	for rows.Next() {
		var round model.Round

		err := rows.Scan(&round.Id, &round.RoundNumber, &round.Theme, &round.ThemeDescription)
		if err != nil {
			log.Fatal(err)
		}

		round.Questions = repository.getQuestions(round.Id)

		rounds = append(rounds, round)
	}

	return rounds
}

func (repository *TriviaRepository) getQuestions(roundId int64) []model.Question {
	selectQuestionsStatement := `
  SELECT question_number, question
  FROM dt.question
  WHERE round_id = $1
  `
	rows, err := repository.db.Query(selectQuestionsStatement, roundId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var questions []model.Question

	for rows.Next() {
		var question model.Question

		err := rows.Scan(&question.QuestionNumber, &question.Question)
		if err != nil {
			log.Fatal(err)
		}

		questions = append(questions, question)
	}

	return questions
}

func getAudioBinary(audioFileName string) string {
	audioFileDirectory := os.Getenv("AUDIO_FILE_DIRECTORY")

	fileName := audioFileDirectory + audioFileName

	content, err := ioutil.ReadFile(fileName)

	if err != nil {
		log.Fatal(err)
	}

	encodedFile := b64.StdEncoding.EncodeToString(content)

	return encodedFile
}

func (repository *TriviaRepository) AddTrivia(newTrivia model.Trivia) {

	audioFileName := ""

	if len(newTrivia.AudioBinary) > 0 {
		audioFileName = writeAudioFile(newTrivia.AudioBinary)
	}

	insertTriviaStatement := `
  INSERT INTO dt.trivia(image_round_theme, image_round_detail, image_round_url, audio_round_theme, answer_url, audio_file_name)
  VALUES($1, $2, $3, $4, $5, $6)
  RETURNING id`

	err := repository.db.QueryRow(insertTriviaStatement, newTrivia.ImageRoundTheme, newTrivia.ImageRoundDetail, newTrivia.ImageRoundURL, newTrivia.AudioRoundTheme, newTrivia.AnswersURL, audioFileName).Scan(&newTrivia.Id)
	if err != nil {
		panic(err)
	}

	for _, round := range newTrivia.Rounds {
		repository.addRound(round, newTrivia.Id)
	}

}

func (repository *TriviaRepository) MarkTriviaUsed(triviaId int64) {
	updateTriviaStatement := `
	UPDATE dt.trivia
	SET used = true,
		date_used = CURRENT_DATE
	WHERE id = $1`

	_, err := repository.db.Exec(updateTriviaStatement, triviaId)
	if err != nil {
		panic(err)
	}

}

func (repository *TriviaRepository) addRound(newRound model.Round, triviaId int64) {
	insertRoundStatement := `
  INSERT INTO dt.round(trivia_id, round_number, theme, theme_description)
  VALUES($1, $2, $3, $4)
  RETURNING id`

	err := repository.db.QueryRow(insertRoundStatement, triviaId, newRound.RoundNumber, newRound.Theme, newRound.ThemeDescription).Scan(&newRound.Id)
	if err != nil {
		panic(err)
	}

	for _, question := range newRound.Questions {
		repository.addQuestion(question, newRound.Id)
	}
}

func (repository *TriviaRepository) addQuestion(newQuestion model.Question, roundId int64) {
	insertQuestionStatement := `
  INSERT INTO dt.question(round_id, question_number, question)
  VALUES($1, $2, $3)`

	_, err := repository.db.Exec(insertQuestionStatement, roundId, newQuestion.QuestionNumber, newQuestion.Question)
	if err != nil {
		panic(err)
	}
}

func writeAudioFile(audioBinary string) string {
	// fmt.Println(audioBinary)
	audioFileDirectory := os.Getenv("AUDIO_FILE_DIRECTORY")

	uuidWithHyphen := uuid.New()
	uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)

	fileName := uuid + ".mp3"

	f, err := os.Create(audioFileDirectory + fileName)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	sDec, _ := b64.StdEncoding.DecodeString(audioBinary)
	data := []byte(sDec)

	_, err2 := f.Write(data)

	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println("done")

	return fileName
}
