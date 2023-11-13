package db

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
	RoundType        RoundType  `json:"roundType"`
}

type Trivia struct {
	Id          int64   `json:"id"`
	Rounds      []Round `json:"rounds"`
	AnswersURL  string  `json:"answersURL"`
	AudioBinary string  `json:"audioBinary"`
	//TODO: ideally we should have separate models for the DB and for the API - but this works for now
	AudioFileName    string
	AudioRoundTheme  string `json:"audioRoundTheme"`
	ImageRoundDetail string `json:"imageRoundDetail"`
	ImageRoundTheme  string `json:"imageRoundTheme"`
	ImageRoundURL    string `json:"imageRoundURL"`
}

type RoundType struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
