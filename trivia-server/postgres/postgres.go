package postgres

import "context"

//go:generate mockgen -source=postgres.go -package postgres -destination=postgres_mock.go

type TriviaDB interface {
	GetNewTrivia() (*Trivia, string, error)
	AddTrivia(ctx *context.Context, newTrivia *Trivia, audioFileName string) error
	MarkTriviaUsed(triviaId int64) error
	RoundTypesList() ([]*RoundType, error)
}
