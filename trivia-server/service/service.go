package service

import (
	"context"

	db "github.com/Btlachac/discord-trivia-bot/postgres"
)

//go:generate mockgen -source=service.go -package service -destination=service_mock.go

type Service interface {
	GetNewTrivia() (*db.Trivia, error)
	AddTrivia(ctx *context.Context, newTrivia *db.Trivia) error
	MarkTriviaUsed(triviaId int64) error
	RoundTypesList() ([]*db.RoundType, error)
}
