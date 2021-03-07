package model

import (
	"database/sql"
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
