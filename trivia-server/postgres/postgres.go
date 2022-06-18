package postgres

import "context"

//go:generate mockgen -source=postgres.go -package postgres -destination=postgres_mock.go

type TriviaDB interface {
	GetNewTrivia(ctx context.Context) (*Trivia, string, error)
	AddTrivia(ctx context.Context, newTrivia *Trivia, audioFileName string) error
	MarkTriviaUsed(ctx context.Context, triviaId int64) error
	RoundTypesList(ctx context.Context) ([]*RoundType, error)
}
