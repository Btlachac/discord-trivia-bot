package service

import (
	"context"
	b64 "encoding/base64"
	"os"
	"strings"

	"github.com/google/uuid"

	"go-trivia-api/internal/db"
)

// TODO: think about imports here
type triviaRepository interface {
	GetNewTrivia(ctx context.Context) (db.Trivia, string, error)
	AddTrivia(ctx context.Context, newTrivia db.Trivia, audioFileName string) error
	MarkTriviaUsed(ctx context.Context, triviaId int64) error
	RoundTypesList(ctx context.Context) ([]db.RoundType, error)
}

type TriviaService struct {
	triviaRepository   triviaRepository
	audioFileDirectory string
}

func NewTriviaService(triviaRepository triviaRepository, audioFileDirectory string) *TriviaService {
	return &TriviaService{
		triviaRepository:   triviaRepository,
		audioFileDirectory: audioFileDirectory,
	}
}

func (s *TriviaService) GetNewTrivia(ctx context.Context) (db.Trivia, error) {
	trivia, audioFileName, err := s.triviaRepository.GetNewTrivia(ctx)

	if err != nil {
		return trivia, err
	}

	if len(audioFileName) > 0 {
		trivia.AudioBinary, err = s.getAudioBinary(audioFileName)
	}

	return trivia, err
}

func (s *TriviaService) AddTrivia(ctx context.Context, newTrivia db.Trivia) error {
	audioFileName := ""
	var err error
	if len(newTrivia.AudioBinary) > 0 {
		audioFileName, err = s.writeAudioFile(newTrivia.AudioBinary)
		if err != nil {
			return err
		}
	}

	return s.triviaRepository.AddTrivia(ctx, newTrivia, audioFileName)
}

func (s *TriviaService) MarkTriviaUsed(ctx context.Context, triviaId int64) error {
	return s.triviaRepository.MarkTriviaUsed(ctx, triviaId)
}

func (s *TriviaService) RoundTypesList(ctx context.Context) ([]db.RoundType, error) {
	return s.triviaRepository.RoundTypesList(ctx)
}

func (s *TriviaService) writeAudioFile(audioBinary string) (string, error) {
	uuidWithHyphen := uuid.New()
	uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)

	fileName := uuid + ".mp3"

	f, err := os.Create(s.audioFileDirectory + fileName)
	if err != nil {
		return "", err
	}

	defer f.Close()

	sDec, err := b64.StdEncoding.DecodeString(audioBinary)
	if err != nil {
		return "", err
	}

	data := []byte(sDec)

	_, err = f.Write(data)

	if err != nil {
		return fileName, err
	}

	return fileName, nil
}

func (s *TriviaService) getAudioBinary(audioFileName string) (string, error) {
	fileName := s.audioFileDirectory + audioFileName

	content, err := os.ReadFile(fileName)

	if err != nil {
		return "", err
	}

	encodedFile := b64.StdEncoding.EncodeToString(content)

	return encodedFile, nil
}
