package service

import (
	b64 "encoding/base64"
	"go-trivia-api/model"
	"io/ioutil"
	"os"
	"strings"

	"github.com/google/uuid"
)

type TriviaService struct {
	triviaRepository triviaRepository
}

type triviaRepository interface {
	GetNewTrivia() (model.Trivia, string, error)
	AddTrivia(newTrivia model.Trivia, audioFileName string) error
	MarkTriviaUsed(triviaId int64) error
	RoundTypesList() ([]model.RoundType, error)
}

func NewTriviaService(triviaRepository triviaRepository) *TriviaService {
	return &TriviaService{
		triviaRepository: triviaRepository,
	}
}

func (service *TriviaService) GetNewTrivia() (model.Trivia, error) {
	trivia, audioFileName, err := service.triviaRepository.GetNewTrivia()

	if err != nil {
		return trivia, err
	}

	if len(audioFileName) > 0 {
		trivia.AudioBinary, err = getAudioBinary(audioFileName)
	}

	return trivia, err
}

func (service *TriviaService) AddTrivia(newTrivia model.Trivia) error {
	audioFileName := ""
	var err error
	if len(newTrivia.AudioBinary) > 0 {
		audioFileName, err = writeAudioFile(newTrivia.AudioBinary)
		if err != nil {
			return err
		}
	}

	return service.triviaRepository.AddTrivia(newTrivia, audioFileName)
}

func (service *TriviaService) MarkTriviaUsed(triviaId int64) error {
	return service.triviaRepository.MarkTriviaUsed(triviaId)
}

func (service *TriviaService) RoundTypesList() ([]model.RoundType, error) {
	return service.triviaRepository.RoundTypesList()
}

func writeAudioFile(audioBinary string) (string, error) {
	audioFileDirectory := os.Getenv("AUDIO_FILE_DIRECTORY")

	uuidWithHyphen := uuid.New()
	uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)

	fileName := uuid + ".mp3"

	f, err := os.Create(audioFileDirectory + fileName)
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

func getAudioBinary(audioFileName string) (string, error) {
	audioFileDirectory := os.Getenv("AUDIO_FILE_DIRECTORY")

	fileName := audioFileDirectory + audioFileName

	content, err := ioutil.ReadFile(fileName)

	if err != nil {
		return "", err
	}

	encodedFile := b64.StdEncoding.EncodeToString(content)

	return encodedFile, nil
}
