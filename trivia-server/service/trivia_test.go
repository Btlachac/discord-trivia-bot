package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Btlachac/discord-trivia-bot/postgres"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

func Test_NewTriviaService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := postgres.NewMockTriviaDB(ctrl)
	mockLogger := zap.NewNop()

	_ = NewTriviaService(mockDB, "", mockLogger)
}

func Test_GetNewTrivia_DB_Error(t *testing.T) {
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

func Test_GetNewTrivia_Error_Bad_Audio_File(t *testing.T) {
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

func Test_AddTrivia_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := postgres.NewMockTriviaDB(ctrl)
	mockLogger := zap.NewNop()

	mockDB.EXPECT().
		AddTrivia(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1).
		Return(nil)

	triviaService := NewTriviaService(mockDB, "", mockLogger)

	ctx := context.Background()
	if err := triviaService.AddTrivia(ctx, &postgres.Trivia{}); err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func Test_AddTrivia_Success_Clean_ImageRoundURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := postgres.NewMockTriviaDB(ctrl)
	mockLogger := zap.NewNop()

	expectedTrivia := &postgres.Trivia{
		ImageRoundURL: "test",
	}
	mockDB.EXPECT().
		AddTrivia(gomock.Any(), expectedTrivia, gomock.Any()).
		Times(1).
		Return(nil)

	triviaService := NewTriviaService(mockDB, "", mockLogger)

	trivia := postgres.Trivia{
		ImageRoundURL: "test/presentxgfdgfdgfd",
	}

	ctx := context.Background()
	if err := triviaService.AddTrivia(ctx, &trivia); err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func Test_AddTrivia_DB_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := postgres.NewMockTriviaDB(ctrl)
	mockLogger := zap.NewNop()

	mockDB.EXPECT().
		AddTrivia(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1).
		Return(errors.New("db error"))

	triviaService := NewTriviaService(mockDB, "", mockLogger)

	ctx := context.Background()
	err := triviaService.AddTrivia(ctx, &postgres.Trivia{})
	if err == nil || err.Error() != "db error" {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func Test_MarkTriviaUsed_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := postgres.NewMockTriviaDB(ctrl)
	mockLogger := zap.NewNop()

	mockDB.EXPECT().
		MarkTriviaUsed(gomock.Any(), gomock.Any()).
		Times(1).
		Return(nil)

	triviaService := NewTriviaService(mockDB, "", mockLogger)

	ctx := context.Background()
	err := triviaService.MarkTriviaUsed(ctx, 1)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func Test_MarkTriviaUsed_DB_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := postgres.NewMockTriviaDB(ctrl)
	mockLogger := zap.NewNop()

	mockDB.EXPECT().
		MarkTriviaUsed(gomock.Any(), gomock.Any()).
		Times(1).
		Return(errors.New("db error"))

	triviaService := NewTriviaService(mockDB, "", mockLogger)

	ctx := context.Background()
	err := triviaService.MarkTriviaUsed(ctx, 1)
	if err == nil || err.Error() != "db error" {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func Test_RoundTypesList_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := postgres.NewMockTriviaDB(ctrl)
	mockLogger := zap.NewNop()

	mockDB.EXPECT().
		RoundTypesList(gomock.Any()).
		Times(1).
		Return(nil, nil)

	triviaService := NewTriviaService(mockDB, "", mockLogger)

	ctx := context.Background()
	_, err := triviaService.RoundTypesList(ctx)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func Test_RoundTypesList_DB_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := postgres.NewMockTriviaDB(ctrl)
	mockLogger := zap.NewNop()

	mockDB.EXPECT().
		RoundTypesList(gomock.Any()).
		Times(1).
		Return(nil, errors.New("db error"))

	triviaService := NewTriviaService(mockDB, "", mockLogger)

	ctx := context.Background()
	_, err := triviaService.RoundTypesList(ctx)
	if err == nil || err.Error() != "db error" {
		t.Errorf("unexpected error: %s", err.Error())
	}
}
