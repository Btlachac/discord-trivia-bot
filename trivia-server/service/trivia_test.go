package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Btlachac/discord-trivia-bot/postgres"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

func TestNewTriviaService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := postgres.NewMockTriviaDB(ctrl)
	mockLogger := zap.NewNop()

	_ = NewTriviaService(mockDB, "", mockLogger)
}

func TestGetNewTrivia_DB_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := postgres.NewMockTriviaDB(ctrl)
	mockLogger := zap.NewNop()

	mockDB.EXPECT().
		GetNewTrivia(gomock.Any()).
		Times(1).
		Return(nil, "", errors.New("unexpected db error"))

	triviaService := NewTriviaService(mockDB, "", mockLogger)

	ctx := context.Background()
	_, err := triviaService.GetNewTrivia(ctx)

	if err == nil || err.Error() != "unexpected db error" {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func TestGetNewTrivia_Error_Bad_Audio_File(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := postgres.NewMockTriviaDB(ctrl)
	mockLogger := zap.NewNop()

	mockDB.EXPECT().
		GetNewTrivia(gomock.Any()).
		Times(1).
		Return(&postgres.Trivia{}, "nil", nil)

	triviaService := NewTriviaService(mockDB, "", mockLogger)

	ctx := context.Background()
	if _, err := triviaService.GetNewTrivia(ctx); err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
}
