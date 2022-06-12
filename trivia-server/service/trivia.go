package s

import (
	b64 "encoding/base64"
	db "go-trivia-api/postgres"
	"io/ioutil"
	"os"
	"strings"

	"github.com/google/uuid"
)

type TriviaService struct {
	triviaDB db.TriviaDB
}

func NewTriviaService(triviaDB db.TriviaDB) *TriviaService {
	return &TriviaService{
		triviaDB: triviaDB,
	}
}

func (s *TriviaService) GetNewTrivia() (*db.Trivia, error) {
	trivia, audioFileName, err := s.triviaDB.GetNewTrivia()

	if err != nil {
		return trivia, err
	}

	if len(audioFileName) > 0 {
		trivia.AudioBinary, err = getAudioBinary(audioFileName)
	}

	return trivia, err
}

func (s *TriviaService) AddTrivia(newTrivia *db.Trivia) error {
	audioFileName := ""
	var err error
	if len(newTrivia.AudioBinary) > 0 {
		audioFileName, err = writeAudioFile(newTrivia.AudioBinary)
		if err != nil {
			return err
		}
	}

	return s.triviaDB.AddTrivia(newTrivia, audioFileName)
}

func (s *TriviaService) MarkTriviaUsed(triviaId int64) error {
	return s.triviaDB.MarkTriviaUsed(triviaId)
}

func (s *TriviaService) RoundTypesList() ([]*db.RoundType, error) {
	return s.triviaDB.RoundTypesList()
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
