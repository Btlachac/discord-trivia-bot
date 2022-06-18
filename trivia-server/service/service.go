package service

import (
	"context"

	db "github.com/Btlachac/discord-trivia-bot/postgres"
)

//go:generate mockgen -source=service.go -package service -destination=service_mock.go

type Service interface {
	GetNewTrivia(ctx context.Context) (*db.Trivia, error)
	AddTrivia(ctx context.Context, newTrivia *db.Trivia) error
	MarkTriviaUsed(ctx context.Context, triviaId int64) error
	RoundTypesList(ctx context.Context) ([]*db.RoundType, error)
}
