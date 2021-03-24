package service

import (
	b64 "encoding/base64"
	"go-trivia-api/model"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
)

type TriviaService struct {
	triviaRepository triviaRepository
}

type triviaRepository interface {
	GetNewTrivia() (model.Trivia, string, error)
	AddTrivia(newTrivia model.Trivia, audioFileName string)
	MarkTriviaUsed(triviaId int64)
}

func NewTriviaService(triviaRepository triviaRepository) *TriviaService {
	return &TriviaService{
		triviaRepository: triviaRepository,
	}
}

func (service *TriviaService) GetNewTrivia() model.Trivia {
	trivia, audioFileName, _ := service.triviaRepository.GetNewTrivia()

	if len(audioFileName) > 0 {
		trivia.AudioBinary = getAudioBinary(audioFileName)
	}

	return trivia
}

func (service *TriviaService) AddTrivia(newTrivia model.Trivia) {
	audioFileName := ""
	if len(newTrivia.AudioBinary) > 0 {
		audioFileName = writeAudioFile(newTrivia.AudioBinary)
	}

	service.triviaRepository.AddTrivia(newTrivia, audioFileName)
}

func (service *TriviaService) MarkTriviaUsed(triviaId int64) {
	service.triviaRepository.MarkTriviaUsed(triviaId)
}

func writeAudioFile(audioBinary string) string {
	audioFileDirectory := os.Getenv("AUDIO_FILE_DIRECTORY")

	_ = os.Mkdir(audioFileDirectory, os.ModeDir)

	uuidWithHyphen := uuid.New()
	uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)

	fileName := uuid + ".mp3"

	f, err := os.Create(audioFileDirectory + fileName)
	if err != nil {
		log.Print(err)
	}

	defer f.Close()

	sDec, _ := b64.StdEncoding.DecodeString(audioBinary)
	data := []byte(sDec)

	_, err2 := f.Write(data)

	if err2 != nil {
		log.Print(err2)
	}

	return fileName
}

func getAudioBinary(audioFileName string) string {
	audioFileDirectory := os.Getenv("AUDIO_FILE_DIRECTORY")

	fileName := audioFileDirectory + audioFileName

	content, err := ioutil.ReadFile(fileName)

	if err != nil {
		log.Fatal(err)
	}

	encodedFile := b64.StdEncoding.EncodeToString(content)

	return encodedFile
}
