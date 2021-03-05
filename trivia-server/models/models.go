package models

import (
	"database/sql"
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
)

type Model struct {
	*sql.DB
}

type Question struct {
	Question       string `json:"question"`
	QuestionNumber int    `json:"questionNumber"`
}

type Round struct {
	Id               int64      `json:"id"`
	Questions        []Question `json:"questions"`
	RoundNumber      int        `json:"roundNumber"`
	Theme            string     `json:"theme"`
	ThemeDescription string     `json:"themeDescription"`
}

type Trivia struct {
	Id               int64   `json:"id"`
	Rounds           []Round `json:"rounds"`
	AnswersURL       string  `json:"answersURL"`
	AudioBinary      string  `json:"audioBinary"`
	AudioRoundTheme  string  `json:"audioRoundTheme"`
	ImageRoundDetail string  `json:"imageRoundDetail"`
	ImageRoundTheme  string  `json:"imageRoundTheme"`
	ImageRoundURL    string  `json:"imageRoundURL"`
}

func NewModel(db *sql.DB) *Model {
	m := Model{}
	m.DB = db
	return &m
}

func (m *Model) GetNewTrivia() Trivia {
	selectTriviaStatement := `
  SELECT id, image_round_theme, image_round_detail, image_round_url, audio_round_theme, answer_url, audio_file_name
  FROM dt.trivia
  WHERE used = false
  ORDER BY date_created ASC
  FETCH FIRST ROW ONLY`

	var trivia Trivia
	var audioFileName sql.NullString

	err := m.DB.QueryRow(selectTriviaStatement).Scan(&trivia.Id, &trivia.ImageRoundTheme, &trivia.ImageRoundDetail, &trivia.ImageRoundURL, &trivia.AudioRoundTheme, &trivia.AnswersURL, &audioFileName)
	if err != nil {
		panic(err)
	}

	if audioFileName.Valid && len(audioFileName.String) > 0 {
		trivia.AudioBinary = getAudioBinary(audioFileName.String)
	}

	trivia.Rounds = m.getRounds(trivia.Id)

	return trivia
}

func (m *Model) getRounds(triviaId int64) []Round {
	selectRoundsStatement := `
  SELECT id, round_number, theme, theme_description
  FROM dt.round
  WHERE trivia_id = $1
  `
	rows, err := m.DB.Query(selectRoundsStatement, triviaId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var rounds []Round

	for rows.Next() {
		var round Round

		err := rows.Scan(&round.Id, &round.RoundNumber, &round.Theme, &round.ThemeDescription)
		if err != nil {
			log.Fatal(err)
		}

		round.Questions = m.getQuestions(round.Id)

		rounds = append(rounds, round)
	}

	return rounds
}

func (m *Model) getQuestions(roundId int64) []Question {
	selectQuestionsStatement := `
  SELECT question_number, question
  FROM dt.question
  WHERE round_id = $1
  `
	rows, err := m.DB.Query(selectQuestionsStatement, roundId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var questions []Question

	for rows.Next() {
		var question Question

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

func (m *Model) AddTrivia(newTrivia Trivia) {

	audioFileName := ""

	if len(newTrivia.AudioBinary) > 0 {
		audioFileName = writeAudioFile(newTrivia.AudioBinary)
	}

	insertTriviaStatement := `
  INSERT INTO dt.trivia(image_round_theme, image_round_detail, image_round_url, audio_round_theme, answer_url, audio_file_name)
  VALUES($1, $2, $3, $4, $5, $6)
  RETURNING id`

	err := m.DB.QueryRow(insertTriviaStatement, newTrivia.ImageRoundTheme, newTrivia.ImageRoundDetail, newTrivia.ImageRoundURL, newTrivia.AudioRoundTheme, newTrivia.AnswersURL, audioFileName).Scan(&newTrivia.Id)
	if err != nil {
		panic(err)
	}

	for _, round := range newTrivia.Rounds {
		m.addRound(round, newTrivia.Id)
	}

}

func (m *Model) MarkTriviaUsed(triviaId int64) {
	updateTriviaStatement := `
	UPDATE dt.trivia
	SET used = true,
		date_used = CURRENT_DATE
	WHERE id = $1`

	_, err := m.DB.Exec(updateTriviaStatement, triviaId)
	if err != nil {
		panic(err)
	}

}

func (m *Model) addRound(newRound Round, triviaId int64) {
	insertRoundStatement := `
  INSERT INTO dt.round(trivia_id, round_number, theme, theme_description)
  VALUES($1, $2, $3, $4)
  RETURNING id`

	err := m.DB.QueryRow(insertRoundStatement, triviaId, newRound.RoundNumber, newRound.Theme, newRound.ThemeDescription).Scan(&newRound.Id)
	if err != nil {
		panic(err)
	}

	for _, question := range newRound.Questions {
		m.addQuestion(question, newRound.Id)
	}
}

func (m *Model) addQuestion(newQuestion Question, roundId int64) {
	insertQuestionStatement := `
  INSERT INTO dt.question(round_id, question_number, question)
  VALUES($1, $2, $3)`

	_, err := m.DB.Exec(insertQuestionStatement, roundId, newQuestion.QuestionNumber, newQuestion.Question)
	if err != nil {
		panic(err)
	}
}

//TODO: retrieving this isn't working

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
