package model

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

// func (q *Question) Validate() error {
// 	var err error
// 	if len(q.Question) == 0 {

// 	}

// 	if q.QuestionNumber < 1 || q.QuestionNumber > 5 {

// 	}
// }
