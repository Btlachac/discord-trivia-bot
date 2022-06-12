package postgres

type TriviaDB interface {
	GetNewTrivia() (*Trivia, string, error)
	AddTrivia(newTrivia *Trivia, audioFileName string) error
	MarkTriviaUsed(triviaId int64) error
	RoundTypesList() ([]*RoundType, error)
}
