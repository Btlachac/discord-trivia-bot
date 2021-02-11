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

var DB *sql.DB

type Question struct {
	Question       string `json:"question"`
	QuestionNumber int    `json:"questionNumber"`
}

type Round struct {
	Id          int64      `json:"id"`
	Questions   []Question `json:"questions"`
	RoundNumber int        `json:"roundNumber"`
	Theme       string     `json:"theme"`
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

func GetNewTrivia() Trivia {
	selectTriviaStatement := `
  SELECT id, image_round_theme, image_round_detail, image_round_url, audio_round_theme, answer_url, audio_file_name
  FROM dt.trivia
  WHERE used = false
  ORDER BY date_created ASC
  FETCH FIRST ROW ONLY`

	var trivia Trivia
	var audioFileName sql.NullString

	err := DB.QueryRow(selectTriviaStatement).Scan(&trivia.Id, &trivia.ImageRoundTheme, &trivia.ImageRoundDetail, &trivia.ImageRoundURL, &trivia.AudioRoundTheme, &trivia.AnswersURL, &audioFileName)
	if err != nil {
		panic(err)
	}

	if audioFileName.Valid && len(audioFileName.String) > 0 {
		trivia.AudioBinary = getAudioBinary(audioFileName.String)
	}

	trivia.Rounds = getRounds(trivia.Id)


	return trivia
}

func getRounds(triviaId int64) []Round {
	selectRoundsStatement := `
  SELECT id, round_number, theme
  FROM dt.round
  WHERE trivia_id = $1
  `
	rows, err := DB.Query(selectRoundsStatement, triviaId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var rounds []Round

	for rows.Next() {
		var round Round

		err := rows.Scan(&round.Id, &round.RoundNumber, &round.Theme)
		if err != nil {
			log.Fatal(err)
		}

		round.Questions = getQuestions(round.Id)

		rounds = append(rounds, round)
	}

	return rounds
}

func getQuestions(roundId int64) []Question {
	selectQuestionsStatement := `
  SELECT question_number, question
  FROM dt.question
  WHERE round_id = $1
  `
	rows, err := DB.Query(selectQuestionsStatement, roundId)
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

func AddTrivia(newTrivia Trivia) {

	audioFileName := ""

	if len(newTrivia.AudioBinary) > 0 {
		audioFileName = writeAudioFile(newTrivia.AudioBinary)
	}

	insertTriviaStatement := `
  INSERT INTO dt.trivia(image_round_theme, image_round_detail, image_round_url, audio_round_theme, answer_url, audio_file_name)
  VALUES($1, $2, $3, $4, $5, $6)
  RETURNING id`

	err := DB.QueryRow(insertTriviaStatement, newTrivia.ImageRoundTheme, newTrivia.ImageRoundDetail, newTrivia.ImageRoundURL, newTrivia.AudioRoundTheme, newTrivia.AnswersURL, audioFileName).Scan(&newTrivia.Id)
	if err != nil {
		panic(err)
	}

	for _, round := range newTrivia.Rounds {
		addRound(round, newTrivia.Id)
	}

}

func addRound(newRound Round, triviaId int64) {
	insertRoundStatement := `
  INSERT INTO dt.round(trivia_id, round_number, theme)
  VALUES($1, $2, $3)
  RETURNING id`

	err := DB.QueryRow(insertRoundStatement, triviaId, newRound.RoundNumber, newRound.Theme).Scan(&newRound.Id)
	if err != nil {
		panic(err)
	}

	for _, question := range newRound.Questions {
		addQuestion(question, newRound.Id)
	}
}

func addQuestion(newQuestion Question, roundId int64) {
	insertQuestionStatement := `
  INSERT INTO dt.question(round_id, question_number, question)
  VALUES($1, $2, $3)`

	_, err := DB.Exec(insertQuestionStatement, roundId, newQuestion.QuestionNumber, newQuestion.Question)
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
