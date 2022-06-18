package service

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	db "github.com/Btlachac/discord-trivia-bot/postgres"
	"go.uber.org/zap"

	"github.com/google/uuid"
)

type TriviaService struct {
	triviaDB           db.TriviaDB
	audioFileDirectory string
	logger             *zap.Logger
}

func NewTriviaService(triviaDB db.TriviaDB, audioFileDirectory string, logger *zap.Logger) *TriviaService {
	return &TriviaService{
		triviaDB:           triviaDB,
		audioFileDirectory: audioFileDirectory,
		logger:             logger,
	}
}

func (s *TriviaService) GetNewTrivia() (*db.Trivia, error) {
	trivia, audioFileName, err := s.triviaDB.GetNewTrivia()

	if err != nil {
		return nil, err
	}

	if len(audioFileName) > 0 {
		trivia.AudioBinary, err = s.getAudioBinary(audioFileName)
	}

	if err != nil {
		s.logger.Warn("error occurred retrieving audio file ", zap.String("filename", audioFileName), zap.Error(err))
		trivia.AudioBinary = ""
	}

	return trivia, nil
}

func (s *TriviaService) AddTrivia(ctx *context.Context, newTrivia *db.Trivia) error {
	audioFileName := ""
	var err error
	if len(newTrivia.AudioBinary) > 0 {
		audioFileName, err = s.writeAudioFile(newTrivia.AudioBinary)
		if err != nil {
			return err
		}
	}

	return s.triviaDB.AddTrivia(ctx, newTrivia, audioFileName)
}

func (s *TriviaService) MarkTriviaUsed(triviaId int64) error {
	return s.triviaDB.MarkTriviaUsed(triviaId)
}

func (s *TriviaService) RoundTypesList() ([]*db.RoundType, error) {
	return s.triviaDB.RoundTypesList()
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

	if _, err := os.Stat(fileName); err != nil {
		return "", fmt.Errorf("file did not exist.  %s", err.Error())
	}

	content, err := ioutil.ReadFile(fileName)

	if err != nil {
		return "", err
	}

	encodedFile := b64.StdEncoding.EncodeToString(content)

	return encodedFile, nil
}
