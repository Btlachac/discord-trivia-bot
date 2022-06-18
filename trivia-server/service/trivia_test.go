package service

import (
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
		GetNewTrivia().
		Times(1).
		Return(nil, "", errors.New("unexpected db error"))

	triviaService := NewTriviaService(mockDB, "", mockLogger)

	_, err := triviaService.GetNewTrivia()

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
		GetNewTrivia().
		Times(1).
		Return(&postgres.Trivia{}, "nil", nil)

	triviaService := NewTriviaService(mockDB, "", mockLogger)

	if _, err := triviaService.GetNewTrivia(); err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
}
