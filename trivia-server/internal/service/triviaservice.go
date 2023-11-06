package service

import (
	b64 "encoding/base64"
	"io/ioutil"
	"os"
	"strings"

	"github.com/google/uuid"

	"go-trivia-api/internal/db"
)

//TODO: think about imports here
type triviaRepository interface {
	GetNewTrivia() (db.Trivia, string, error)
	AddTrivia(newTrivia db.Trivia, audioFileName string) error
	MarkTriviaUsed(triviaId int64) error
	RoundTypesList() ([]db.RoundType, error)
}

type TriviaService struct {
	triviaRepository triviaRepository
}

func NewTriviaService(triviaRepository triviaRepository) *TriviaService {
	return &TriviaService{
		triviaRepository: triviaRepository,
	}
}

func (service *TriviaService) GetNewTrivia() (db.Trivia, error) {
	trivia, audioFileName, err := service.triviaRepository.GetNewTrivia()

	if err != nil {
		return trivia, err
	}

	if len(audioFileName) > 0 {
		trivia.AudioBinary, err = getAudioBinary(audioFileName)
	}

	return trivia, err
}

func (service *TriviaService) AddTrivia(newTrivia db.Trivia) error {
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

func (service *TriviaService) RoundTypesList() ([]db.RoundType, error) {
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
